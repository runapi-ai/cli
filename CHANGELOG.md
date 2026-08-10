# Changelog

## [v0.10.0](https://github.com/runapi-ai/cli/releases/tag/v0.10.0) - 2026-08-10

### Breaking
- Document the 300-character prompt limit for OmniHuman audio-to-video requests.
  Migration: Upgrade the CLI and keep OmniHuman audio-to-video prompts at 300 characters or fewer.

### Added
- Add create, poll, and task lookup commands for Suno music inspiration.

### Changed
- Update Seedream 5 Lite output quality constraints in generated CLI help.


## [v0.9.0](https://github.com/runapi-ai/cli/releases/tag/v0.9.0) - 2026-08-07

### Added
- Add audio transcription commands with multipart inputs and JSON, text, and subtitle output preservation.
- Add text-to-video and image-to-video commands for MiniMax H3.
- Add Qwen 3 text-to-image and image editing commands with generated contract help.
- Add s2.1-pro and MP3 or WAV output fields to Fish Audio text-to-speech help.

### Changed
- Update Grok Imagine Preview resolution and reference image constraints.


## [v0.8.7](https://github.com/runapi-ai/cli/releases/tag/v0.8.7) - 2026-08-06

### Added
- Add CLI commands for stitching audio, remastering audio, and adding samples from a selected time range.

### Fixed
- Show the optional PixVerse enable_audio field for text-to-video and image-to-video requests.


## [v0.8.6](https://github.com/runapi-ai/cli/releases/tag/v0.8.6) - 2026-08-04

### Added
- Add PixVerse commands for text-to-video, image-to-video, reference-guided video, transitions, and video extension.

### Changed
- Support more concurrent local listeners with lower idle polling frequency and immediate event draining.

### Fixed
- Honor the server's Retry-After delay when retrying local listener requests.


## [v0.8.5](https://github.com/runapi-ai/cli/releases/tag/v0.8.5) - 2026-07-31

### Removed
- Remove seedance-v1-lite from Seedance CLI contract metadata.
  Migration: Use seedance-v1-pro or another supported Seedance model.

### Fixed
- Acknowledge each valid Local Callback Listener event before forwarding so local HTTP failures do not replay or block later events.


## [v0.8.4](https://github.com/runapi-ai/cli/releases/tag/v0.8.4) - 2026-07-29

### Removed
- Remove seedance-v1-lite from Seedance CLI contract metadata.
  Migration: Use seedance-v1-pro or another supported Seedance model.


## [v0.8.3](https://github.com/runapi-ai/cli/releases/tag/v0.8.3) - 2026-07-28

### Changed
- Describe and validate the documented Gemini Omni, Grok Imagine, and Topaz request fields.

### Fixed
- Expose supported Wan Flash image-to-video duration values in generated command metadata; request defaults remain API-applied.


## [v0.8.2](https://github.com/runapi-ai/cli/releases/tag/v0.8.2) - 2026-07-28

### Added
- Add public commands for current Price Schedules and task reservation quotes.
- Add Flux 2 Max text-to-image and remix-image model options with their supported request constraints.
- Expose Kling O1 models and reference image, reference video, reference role, and audio-preservation fields in generated Kling commands.

### Changed
- Send the CLI version during browser sign-in so the authorization page can identify an available update and link to its release.


## [v0.8.1](https://github.com/runapi-ai/cli/releases/tag/v0.8.1) - 2026-07-23

### Added
- Add Kling 2.6 motion-control model and conditional field requirements to generated Kling command help.


## [v0.8.0](https://github.com/runapi-ai/cli/releases/tag/v0.8.0) - 2026-07-23

### Added
- Add the Kling V3 Omni model, output resolution, sound, duration, aspect ratio, and frame constraints to generated Kling command help.
- Add seven Producer FUZZ versions to generated text-to-music model help and validation.
- Add CLI support for continuing completed Kling v2.5 Turbo videos.


## [v0.7.0](https://github.com/runapi-ai/cli/releases/tag/v0.7.0) - 2026-07-22

### Added
- Add Kling 2.6 model, mode, sound, duration, and frame constraints to generated Kling command help.
- Add the Midjourney extend-video command with task creation, polling, and public input help.
- Add Flux text-to-image and remix-image commands with generated request validation and task polling.
- Expose Veo 3.1 Lite model and input constraints in generated CLI help.
- Add Qwen Image text-to-image, remix-image, and edit-image commands with generated contract help.


## [v0.6.0](https://github.com/runapi-ai/cli/releases/tag/v0.6.0) - 2026-07-21

### Added
- Expose request-scoped reference audio entries for Fish Audio text-to-speech input files.
- Add lyric blending command flags and generated contract help.


## [v0.5.0](https://github.com/runapi-ai/cli/releases/tag/v0.5.0) - 2026-07-20

### Breaking
- Replace Grok Imagine image-to-video `source_image_urls` with scalar `source_image_url`.
  Migration: Pass `--source-image-url URL` or provide `source_image_url` in JSON input.

### Added
- Add the synchronous `runapi midjourney shorten-prompt` command with generated input validation and help.
- Add synchronous text-to-speech commands for OpenAI TTS and Fish Audio.
- Expose model and text input help from the generated contract.
- Add Gemini Omni Flash Preview model selection and model-specific text-to-video help.
- Add the gemini-tts text-to-speech command with generated input help and validation.
- Add Seedream 5 Pro text-to-image and edit-image models to CLI help and validation.
- Expose advanced stem separation mode and supported stem names in CLI help and validation.
- Add the Producer text-to-music command with generated contract help and typed SDK wiring.

### Changed
- Expose Seedream 5-Lite `output_format` values and default in CLI help and validation.
- Include minimum and maximum item counts for constrained array fields in generated CLI metadata and help.

### Fixed
- Preserve API-provided error codes in JSON output and omit the code when the response does not provide one.
- Classify command exit and retry behavior by HTTP status or SDK error type instead of generating substitute codes.


## [v0.4.1](https://github.com/runapi-ai/cli/releases/tag/v0.4.1) - 2026-07-20

### Added
- Add `runapi listen --rotate-secret` for explicit per-key Listen Signing Secret recovery.

### Changed
- Verify the selected key echoed by the API, print only the rotated secret, and document listener restart behavior.


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
