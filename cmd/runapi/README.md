# RunAPI CLI

[![Release](https://img.shields.io/github/v/release/runapi-ai/cli?sort=semver)](https://github.com/runapi-ai/cli/releases)
[![License](https://img.shields.io/github/license/runapi-ai/cli)](https://github.com/runapi-ai/cli/blob/main/LICENSE)

`runapi` is the command-line client for [RunAPI](https://runapi.ai). One binary calls 100+ AI models — Suno, Veo, Nano Banana, GPT Image, Imagen, Kling, Flux, Wan, Claude, and the rest — with a JSON-first interface that maps cleanly to scripts, CI pipelines, and agent runtimes.

## Install

**Linux / macOS (servers, CI, Dockerfiles):**

```bash
curl -fsSL https://runapi.ai/cli/install.sh | sh
```

**Homebrew (macOS / Linux):**

```bash
brew install runapi-ai/tap/runapi
```

**Windows (precompiled archive):**

Download `windows-amd64` or `windows-arm64` from the
[latest release](https://github.com/runapi-ai/cli/releases/latest), extract
`runapi.exe`, and place it on `PATH`.

**Go source build (Go 1.22+):**

```bash
go install github.com/runapi-ai/cli/cmd/runapi@latest
```

**Or paste this prompt to your AI agent** (Claude Code, Codex, Gemini CLI, OpenClaw, Hermes):

```
Install the RunAPI CLI on this machine and wire it into this agent runtime:

1. On Linux or macOS, run `curl -fsSL https://runapi.ai/cli/install.sh | sh`.
   On Windows, download the matching `windows-amd64` or `windows-arm64` archive
   from https://github.com/runapi-ai/cli/releases/latest, extract `runapi.exe`,
   and place it on `PATH`.
2. Save my API key with `runapi auth import-token`; it reads `RUNAPI_API_KEY`
   from the current environment. (Or run `runapi login` on a laptop with a
   browser.)
3. Install the CLI skill into this runtime with `runapi agent install-skill --target <claude|codex|gemini|openclaw|hermes>`.
4. Confirm it works by running `runapi version` and `runapi auth status`.

Docs: https://runapi.ai/docs#runapi-cli
```

Precompiled archives are published for macOS, Linux, and Windows on amd64 and
arm64. The curl install script supports Linux and macOS: it detects OS and CPU
architecture, verifies a SHA-256 checksum from
`https://runapi.ai/cli/latest.json`, and writes the binary to `/usr/local/bin`
(or `~/.local/bin` if it cannot get write access). On a musl libc distro
(Alpine) it exits early with a clear message rather than installing a binary
that will crash at runtime.

### Pin a release or pick a directory

```bash
# Pin to a specific tag
curl -fsSL https://runapi.ai/cli/install.sh | sh -s -- --version v0.1.0

# Install somewhere else
curl -fsSL https://runapi.ai/cli/install.sh | sh -s -- --dir ~/.local/bin
```

Environment variables also work: `RUNAPI_VERSION`, `RUNAPI_INSTALL_DIR`, `RUNAPI_INSTALL_BASE`, `RUNAPI_DOWNLOAD_BASE` (set to a private mirror to fetch the release archive from your own CDN), and `RUNAPI_SKIP_LIBC_CHECK=1` (bypass the musl guard once you have a glibc shim).

## Quickstart

```bash
# 1. Authenticate (browser flow on a laptop, headless on a server — see below).
runapi login
runapi auth status

# 2. Generate an image with Nano Banana from JSON.
runapi nano-banana generate --input '{"prompt":"a hummingbird drinking espresso","aspect_ratio":"1:1"}'

# 3. Same call, but read the request from a file (recommended for any non-trivial payload).
cat > request.json <<'JSON'
{
  "prompt": "a hummingbird drinking espresso",
  "aspect_ratio": "1:1",
  "negative_prompt": "blurry, low resolution"
}
JSON
runapi nano-banana generate --input-file request.json

# 4. Fire-and-forget for long jobs, then poll separately.
TASK_ID=$(runapi suno generate --async --input-file song.json | jq -r .id)
runapi wait "$TASK_ID" --service suno --action generate
```

For top-level model media URL fields, readable local file paths are accepted directly. The CLI uploads them before submitting the request and keeps `http://` or `https://` values unchanged:

```bash
runapi gpt-image edit-image --input \
  '{"model":"gpt-image-1.5","prompt":"remove the background","source_image_urls":["./image.png"],"aspect_ratio":"1:1","quality":"medium"}'
```

Use `runapi files create` when you need a temporary URL to reuse, or when the source is a remote URL or Base64 data.

`runapi` writes JSON to stdout and progress lines to stderr, so it composes naturally with `jq` and any pipeline:

```bash
runapi nano-banana generate --input-file request.json \
  | jq -r '.images[].url' \
  | xargs -I{} curl -OL {}
```

## Headless / CI usage

For server images, Dockerfiles, and automation that cannot run an interactive browser:

```bash
# Save the API key without exposing it in `ps -ef` or shell history.
printf '%s' "$RUNAPI_API_KEY" | runapi auth import-token --token -

# Or keep the key in the environment for the duration of the job.
export RUNAPI_API_KEY=<your-key>
runapi auth status
```

`auth import-token` verifies the key against `/api/v1/me` before writing it to `~/.config/runapi/config.json` with mode `0600`. Pass `--skip-verify` when bootstrapping an offline image. `RUNAPI_BASE_URL` overrides the API endpoint.

Build pipelines that need to verify downloads themselves can consume the stable JSON contract directly:

- `https://runapi.ai/cli/latest.json` — current stable release with `version`, asset URLs, SHA-256 sums
- `https://runapi.ai/cli/v<x.y.z>/manifest.json` — pinned, immutable

```bash
curl -fsSL https://runapi.ai/cli/latest.json | jq '.assets["linux-amd64"]'
```

## Local callback listener

Local listeners require the credential issued by browser login. On first use, the CLI lists your enabled API keys and lets you choose one by name, stable ID, and masked value. The validated stable ID is written to `.runapi.toml` at the git root, or in the current directory outside a git repository; no credential or signing secret is stored there.

```bash
runapi login
runapi listen localhost:3000/webhooks/runapi
```

Only tasks created with the selected API key are copied to that listener. Tasks created with another key stay isolated, including keys owned by the same Account member. A committed `.runapi.toml` is reusable by the same member on another machine; another member is prompted to select a key they own.

Selection precedence is `--callback-api-key-id` (one invocation only), then project `.runapi.toml`, then the TTY selector. The selector writes the config only after the server validates the session. The config has one allowed field:

```toml
callback_api_key_id = "token_abc123"
```

Renaming the API key does not invalidate this stable ID. Do not add credentials, names, masks, signing secrets, forwarding URLs, or `base_url`; unknown fields are rejected.

Agents and non-interactive shells can discover and select the key explicitly:

```bash
runapi api-keys list --json
runapi listen localhost:3000/webhooks/runapi --callback-api-key-id token_abc123
```

Without `--callback-api-key-id` or `.runapi.toml`, a non-interactive invocation returns `callback_api_key_required` with the available key metadata and never chooses automatically. `cli_listen_required` means the active credential did not come from browser login: the imported API key keeps its existing API access, but cannot list or select listener keys. Remove any `--api-key` or `RUNAPI_API_KEY` override, run `runapi login`, then retry `runapi listen`. If the selected key becomes unusable, the listener exits without falling back to another key.

Print the selected key's stable Listen Signing Secret without starting a listener with:

```bash
RUNAPI_WEBHOOK_SECRET="$(runapi listen --print-secret --callback-api-key-id token_abc123)"
```

After upgrading from the previous Account-wide listener behavior, update the CLI, run `runapi login` again, and restart listeners. Previous listener sessions and Account-wide listener secrets are invalid after the migration.

## Agent runtimes

`runapi` ships a portable skill for the major AI agent runtimes. Install it once and the runtime can run RunAPI commands with full inline documentation:

```bash
runapi agent install-skill --target claude    # → ~/.claude/skills/runapi-cli
runapi agent install-skill --target codex     # → ~/.agents/skills/runapi-cli
runapi agent install-skill --target gemini    # → ~/.gemini/skills/runapi-cli
runapi agent install-skill --target openclaw  # → ~/.openclaw/skills/runapi-cli
runapi agent install-skill --target hermes    # → ~/.hermes/skills/runapi-cli
runapi agent list-targets                     # JSON list of supported targets
```

The skill content is published at [`runapi-ai/cli-skill`](https://github.com/runapi-ai/cli-skill). By default, `install-skill` installs the latest stable skill tag. Use `--version <tag>` to pin an exact skill release, or `--target-dir <path>` to install into a non-standard location.

## Configuration reference

| Source | Precedence |
|---|---|
| `--api-key <key>` flag | 1 (highest) |
| `RUNAPI_API_KEY` env var | 2 |
| `~/.config/runapi/config.json` (or `$XDG_CONFIG_HOME/runapi/config.json`) | 3 |
| Default base URL `https://runapi.ai` | last resort |

The CLI exits with non-zero status codes you can switch on in scripts:

| Code | Meaning |
|---|---|
| 0 | Success |
| 2 | Authentication failure / unsupported platform |
| 3 | Insufficient credits / missing dependency (`curl`, `sha256sum`) |
| 4 | Validation, not found, or manifest parse error |
| 5 | Timeout, download, or checksum mismatch |
| 6 | Rate limit |
| 7 | Task failed at the model |

## Links

- Documentation: <https://runapi.ai/docs#runapi-cli>
- Model catalog: <https://runapi.ai/models>
- Install script source: <https://runapi.ai/cli/install.sh>
- Release manifest: <https://runapi.ai/cli/latest.json>
- CLI repository: <https://github.com/runapi-ai/cli>
- Homebrew tap: <https://github.com/runapi-ai/homebrew-tap>
- Agent skill: <https://github.com/runapi-ai/cli-skill>

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
