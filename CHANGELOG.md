# Changelog

## [v0.4.0](https://github.com/runapi-ai/cli/releases/tag/v0.4.0) - 2026-07-18

### Changed
- Publish precompiled Windows amd64 and arm64 archives containing `runapi.exe`.
- Update CLI Go module dependencies to the latest published SDK releases.

## [v0.3.1](https://github.com/runapi-ai/cli/releases/tag/v0.3.1) - 2026-07-17

### Changed
- Expose the Fast model and its request constraints in Grok Imagine CLI help and validation.

## [v0.3.0](https://github.com/runapi-ai/cli/releases/tag/v0.3.0) - 2026-07-17

### Added
- Add Midjourney commands for image generation, editing, image-to-video, image-to-prompt, and seed lookup.
- Add interactive and JSON API key discovery, project `.runapi.toml` selection, positional forwarding URLs, and per-key Listen Signing Secret output.

### Security
- Scope callback delivery, listener limits, and signing secrets to the explicitly selected API key.
- Stop listeners when credentials or selected keys become unusable instead of silently selecting another key.

### Upgrade
- Existing listener users must update the CLI, run `runapi login` again, update their webhook secret, and restart local listeners.

## [v0.2.17](https://github.com/runapi-ai/cli/releases/tag/v0.2.17) - 2026-07-16

### Changed
- Add CLI contract metadata for Kling V3 Turbo text-to-video and image-to-video.
- Include generated enum validation for duration, aspect ratio, and output resolution.
- Refresh the CLI generated contract so Grok Imagine text-to-video and image-to-video commands recognize `grok-imagine-video-1.5-preview`.

## [v0.2.16](https://github.com/runapi-ai/cli/releases/tag/v0.2.16) - 2026-07-08

### Changed
- Refresh bundled CLI contract metadata for current SDK validation.

## [v0.2.15](https://github.com/runapi-ai/cli/releases/tag/v0.2.15) - 2026-07-08

### Changed
- Refresh the CLI Seedance commands to depend on the published Seedance Go SDK release.

## [v0.2.14](https://github.com/runapi-ai/cli/releases/tag/v0.2.14) - 2026-07-07

### Added
- Add Nano Banana 2 Lite command metadata.

### Changed
- Publish v0.2.14.

## [v0.2.12](https://github.com/runapi-ai/cli/releases/tag/v0.2.12), [v0.2.13](https://github.com/runapi-ai/cli/releases/tag/v0.2.13) - 2026-07-02

### Fixed
- Command help and input validation now reflect current field names and allowed values from the RunAPI request contract.

## [v0.2.11](https://github.com/runapi-ai/cli/releases/tag/v0.2.11) - 2026-06-25

### Added
- Upload local media file paths automatically in CLI create commands before task creation.

### Changed
- Preserve explicit JSON null entries in media arrays while processing local file uploads.

## [v0.2.10](https://github.com/runapi-ai/cli/releases/tag/v0.2.10) - 2026-06-18

### Added
- File upload command
- Universal account and files resources

## [v0.2.9](https://github.com/runapi-ai/cli/releases/tag/v0.2.9) - 2026-06-01

### Changed
- Align CLI with upstream Input Contract and public API vocabulary changes
- Update all SDK dependencies to latest versions

## [v0.2.5](https://github.com/runapi-ai/cli/releases/tag/v0.2.5) - 2026-05-22

### Changed
- Build the CLI against RunAPI SDK artifacts v0.2.4.
- Refresh CLI README and public release metadata.

## [v0.2.3](https://github.com/runapi-ai/cli/releases/tag/v0.2.3), [v0.2.4](https://github.com/runapi-ai/cli/releases/tag/v0.2.4) - 2026-05-22

Initial release.

## [v0.2.2](https://github.com/runapi-ai/cli/releases/tag/v0.2.2) - 2026-05-20

Initial release.

## [v0.2.1](https://github.com/runapi-ai/cli/releases/tag/v0.2.1) - 2026-05-19

Initial release.
