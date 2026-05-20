# RunAPI CLI

[![Release](https://img.shields.io/github/v/release/runapi-ai/cli?sort=semver)](https://github.com/runapi-ai/cli/releases)
[![License](https://img.shields.io/github/license/runapi-ai/cli)](https://github.com/runapi-ai/cli/blob/main/LICENSE)
[![Go module](https://pkg.go.dev/badge/github.com/runapi-ai/cli.svg)](https://pkg.go.dev/github.com/runapi-ai/cli)

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

**Go source build (Go 1.22+):**

```bash
go install github.com/runapi-ai/cli@latest
```

**Or paste this prompt to your AI agent** (Claude Code, Codex, Gemini CLI, OpenClaw, Hermes):

```
Install the RunAPI CLI on this machine and wire it into this agent runtime:

1. Run `curl -fsSL https://runapi.ai/cli/install.sh | sh`.
2. Save my API key with `printf '%s' "$RUNAPI_API_KEY" | runapi auth import-token --token -`
   (or run `runapi login` if I'm on a laptop with a browser).
3. Install the CLI skill into this runtime with `runapi agent install-skill --target <claude|codex|gemini|openclaw|hermes>`.
4. Confirm it works by running `runapi version` and `runapi auth status`.

Docs: https://runapi.ai/docs#runapi-cli
```

The curl install script detects OS and CPU architecture (Linux and macOS, amd64 and arm64), verifies a SHA-256 checksum from `https://runapi.ai/cli/latest.json`, and writes the binary to `/usr/local/bin` (or `~/.local/bin` if it cannot get write access). On a musl libc distro (Alpine) it exits early with a clear message rather than installing a binary that will crash at runtime.

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

`runapi` writes JSON to stdout and progress lines to stderr, so it composes naturally with `jq` and any pipeline:

```bash
runapi nano-banana generate --input-file request.json \
  | jq -r '.result_urls[]' \
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

## Agent runtimes

`runapi` ships a portable skill for the major AI agent runtimes. Install it once and the runtime can run RunAPI commands with full inline documentation:

```bash
runapi agent install-skill --target claude    # → ~/.claude/skills/cli
runapi agent install-skill --target codex     # → ~/.agents/skills/cli
runapi agent install-skill --target gemini    # → ~/.gemini/skills/cli
runapi agent install-skill --target openclaw  # → ~/.openclaw/skills/cli
runapi agent install-skill --target hermes    # → ~/.hermes/skills/cli
runapi agent list-targets                     # JSON list of supported targets
```

The skill content is published at [`runapi-ai/cli-skill`](https://github.com/runapi-ai/cli-skill) and tagged in lockstep with CLI releases. Use `--target-dir <path>` to install into a non-standard location.

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
