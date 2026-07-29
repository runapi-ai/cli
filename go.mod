module github.com/runapi-ai/cli

go 1.26

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/manifoldco/promptui v0.9.0
	github.com/runapi-ai/core-sdk/go v0.2.16
	github.com/runapi-ai/elevenlabs-sdk/go v0.2.10
	github.com/runapi-ai/fish-audio-sdk/go v0.1.3
	github.com/runapi-ai/flux-2-sdk/go v0.3.0
	github.com/runapi-ai/flux-kontext-sdk/go v0.2.8
	github.com/runapi-ai/flux-sdk/go v0.1.1
	github.com/runapi-ai/gemini-omni-sdk/go v0.3.3
	github.com/runapi-ai/gemini-tts-sdk/go v0.1.2
	github.com/runapi-ai/gpt-4o-image-sdk/go v0.2.9
	github.com/runapi-ai/gpt-image-2-sdk/go v0.2.8
	github.com/runapi-ai/gpt-image-sdk/go v0.2.8
	github.com/runapi-ai/grok-imagine-sdk/go v0.2.12
	github.com/runapi-ai/hailuo-sdk/go v0.2.8
	github.com/runapi-ai/happyhorse-sdk/go v0.2.9
	github.com/runapi-ai/ideogram-v3-sdk/go v0.2.9
	github.com/runapi-ai/imagen-4-sdk/go v0.2.10
	github.com/runapi-ai/infinitetalk-sdk/go v0.2.8
	github.com/runapi-ai/kling-sdk/go v0.2.13
	github.com/runapi-ai/luma-sdk/go v0.2.8
	github.com/runapi-ai/midjourney-sdk/go v0.3.1
	github.com/runapi-ai/nano-banana-sdk/go v0.2.12
	github.com/runapi-ai/omnihuman-sdk/go v0.2.10
	github.com/runapi-ai/openai-tts-sdk/go v0.1.2
	github.com/runapi-ai/producer-sdk/go v0.2.1
	github.com/runapi-ai/qwen-2-sdk/go v0.2.9
	github.com/runapi-ai/qwen-image-sdk/go v0.1.1
	github.com/runapi-ai/recraft-sdk/go v0.2.8
	github.com/runapi-ai/runway-aleph-sdk/go v0.2.8
	github.com/runapi-ai/runway-sdk/go v0.2.9
	github.com/runapi-ai/seedance-sdk/go v0.2.13
	github.com/runapi-ai/seedream-sdk/go v0.2.11
	github.com/runapi-ai/suno-sdk/go v0.3.2
	github.com/runapi-ai/topaz-sdk/go v0.2.9
	github.com/runapi-ai/veo-3.1-sdk/go v0.2.11
	github.com/runapi-ai/volcengine-lip-sync-sdk/go v0.2.10
	github.com/runapi-ai/wan-sdk/go v0.2.11
	github.com/runapi-ai/z-image-sdk/go v0.2.8
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20181122145206-62eef0e2fa9b // indirect
)

retract [v0.0.0, v0.8.3] // Upgrade to v0.8.4 or later to continue using RunAPI.
