package main

var generatedContract = map[string]generatedContractAction{
	"elevenlabs/isolate-audio": {
		Models: []string{"audio-isolation"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"audio-isolation": {},
		},
	},
	"elevenlabs/speech-to-text": {
		Models: []string{"speech-to-text"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"speech-to-text": {},
		},
	},
	"elevenlabs/text-to-dialogue": {
		Models: []string{"text-to-dialogue-v3"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"text-to-dialogue-v3": {
				"stability": {Enum: []any{0.0, 0.5, 1.0}},
			},
		},
	},
	"elevenlabs/text-to-sound": {
		Models: []string{"sound-effect-v2"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"sound-effect-v2": {
				"output_format": {Enum: []any{"mp3_22050_32", "mp3_44100_32", "mp3_44100_64", "mp3_44100_96", "mp3_44100_128", "mp3_44100_192", "pcm_8000", "pcm_16000", "pcm_22050", "pcm_24000", "pcm_44100", "pcm_48000", "ulaw_8000", "alaw_8000", "opus_48000_32", "opus_48000_64", "opus_48000_96", "opus_48000_128", "opus_48000_192"}},
			},
		},
	},
	"elevenlabs/text-to-speech": {
		Models: []string{"text-to-speech-multilingual-v2", "text-to-speech-turbo-v2.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"text-to-speech-multilingual-v2": {},
			"text-to-speech-turbo-v2.5":      {},
		},
	},
	"flux-2/remix-image": {
		Models: []string{"flux-2-flex-remix-image", "flux-2-pro-remix-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-2-flex-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
			"flux-2-pro-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
		},
	},
	"flux-2/text-to-image": {
		Models: []string{"flux-2-flex-text-to-image", "flux-2-pro-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-2-flex-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
			"flux-2-pro-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
		},
	},
	"flux-kontext/text-to-image": {
		Models: []string{"flux-kontext-max", "flux-kontext-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-kontext-max": {
				"aspect_ratio":  {Enum: []any{"21:9", "16:9", "4:3", "1:1", "3:4", "9:16"}},
				"output_format": {Enum: []any{"jpeg", "png"}},
			},
			"flux-kontext-pro": {
				"aspect_ratio":  {Enum: []any{"21:9", "16:9", "4:3", "1:1", "3:4", "9:16"}},
				"output_format": {Enum: []any{"jpeg", "png"}},
			},
		},
	},
	"gemini-omni/create-audio": {
		Models: []string{"gemini-omni-audio"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gemini-omni-audio": {},
		},
	},
	"gemini-omni/create-character": {
		Models: []string{"gemini-omni-character"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gemini-omni-character": {},
		},
	},
	"gemini-omni/text-to-video": {
		Models: []string{"gemini-omni-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gemini-omni-text-to-video": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16"}},
				"duration_seconds":  {Enum: []any{4, 6, 8, 10}},
				"output_resolution": {Enum: []any{"720p", "1080p", "4k"}},
			},
		},
	},
	"gpt-4o-image/text-to-image": {
		Models: []string{"gpt-4o-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-4o-image": {
				"aspect_ratio": {Enum: []any{"1:1", "3:2", "2:3"}},
				"output_count": {Enum: []any{1, 2, 4}},
			},
		},
	},
	"gpt-image-2/edit-image": {
		Models: []string{"gpt-image-2"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-image-2": {
				"aspect_ratio":      {Enum: []any{"auto", "1:1", "3:2", "2:3", "4:3", "3:4", "5:4", "4:5", "16:9", "9:16", "2:1", "1:2", "3:1", "1:3", "21:9", "9:21"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"gpt-image-2/text-to-image": {
		Models: []string{"gpt-image-2"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-image-2": {
				"aspect_ratio":      {Enum: []any{"auto", "1:1", "3:2", "2:3", "4:3", "3:4", "5:4", "4:5", "16:9", "9:16", "2:1", "1:2", "3:1", "1:3", "21:9", "9:21"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"gpt-image/edit-image": {
		Models: []string{"gpt-image-1.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-image-1.5": {
				"aspect_ratio": {Enum: []any{"1:1", "2:3", "3:2"}},
				"quality":      {Enum: []any{"medium", "high"}},
			},
		},
	},
	"gpt-image/text-to-image": {
		Models: []string{"gpt-image-1.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-image-1.5": {
				"aspect_ratio": {Enum: []any{"1:1", "2:3", "3:2"}},
				"quality":      {Enum: []any{"medium", "high"}},
			},
		},
	},
	"grok-imagine/edit-image": {
		Models: []string{"grok-imagine-edit-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-edit-image": {},
		},
	},
	"grok-imagine/extend": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"extension_duration_seconds": {Enum: []any{6, 10}},
			},
		},
	},
	"grok-imagine/image-to-video": {
		Models: []string{"grok-imagine-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-image-to-video": {
				"aspect_ratio":      {Enum: []any{"2:3", "3:2", "1:1", "16:9", "9:16"}},
				"motion_style":      {Enum: []any{"fun", "normal", "spicy"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
		},
	},
	"grok-imagine/text-to-image": {
		Models: []string{"grok-imagine-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-text-to-image": {
				"aspect_ratio": {Enum: []any{"2:3", "3:2", "1:1", "16:9", "9:16"}},
			},
		},
	},
	"grok-imagine/text-to-video": {
		Models: []string{"grok-imagine-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-text-to-video": {
				"aspect_ratio":      {Enum: []any{"2:3", "3:2", "1:1", "16:9", "9:16"}},
				"motion_style":      {Enum: []any{"fun", "normal", "spicy"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
		},
	},
	"grok-imagine/upscale-image": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"hailuo/image-to-video": {
		Models: []string{"hailuo-02-image-to-video-pro", "hailuo-02-image-to-video-standard", "hailuo-2.3-image-to-video-pro", "hailuo-2.3-image-to-video-standard"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"hailuo-02-image-to-video-pro": {},
			"hailuo-02-image-to-video-standard": {
				"duration_seconds":  {Enum: []any{6, 10}},
				"output_resolution": {Enum: []any{"512p", "768p"}},
			},
			"hailuo-2.3-image-to-video-pro": {
				"duration_seconds":  {Enum: []any{6, 10}},
				"output_resolution": {Enum: []any{"768p", "1080p"}},
			},
			"hailuo-2.3-image-to-video-standard": {
				"duration_seconds":  {Enum: []any{6, 10}},
				"output_resolution": {Enum: []any{"768p", "1080p"}},
			},
		},
	},
	"hailuo/text-to-video": {
		Models: []string{"hailuo-02-text-to-video-pro", "hailuo-02-text-to-video-standard"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"hailuo-02-text-to-video-pro": {},
			"hailuo-02-text-to-video-standard": {
				"duration_seconds": {Enum: []any{6, 10}},
			},
		},
	},
	"happyhorse/edit-video": {
		Models: []string{"happyhorse-edit-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"happyhorse-edit-video": {
				"audio_setting":     {Enum: []any{"auto", "original"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"happyhorse/image-to-video": {
		Models: []string{"happyhorse-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"happyhorse-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"happyhorse/text-to-video": {
		Models: []string{"happyhorse-character", "happyhorse-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"happyhorse-character": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"happyhorse-text-to-video": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"ideogram-v3/edit-image": {
		Models: []string{"ideogram-v3-character-edit", "ideogram-v3-edit"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"ideogram-v3-character-edit": {
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "realistic", "fiction"}},
			},
			"ideogram-v3-edit": {
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
			},
		},
	},
	"ideogram-v3/reframe-image": {
		Models: []string{"ideogram-v3-reframe"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"ideogram-v3-reframe": {
				"aspect_ratio":    {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "general", "realistic", "design"}},
			},
		},
	},
	"ideogram-v3/remix-image": {
		Models: []string{"ideogram-v3-character-remix", "ideogram-v3-remix"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"ideogram-v3-character-remix": {
				"aspect_ratio":    {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "realistic", "fiction"}},
			},
			"ideogram-v3-remix": {
				"aspect_ratio":    {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "general", "realistic", "design"}},
			},
		},
	},
	"ideogram-v3/text-to-image": {
		Models: []string{"ideogram-v3-character", "ideogram-v3-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"ideogram-v3-character": {
				"aspect_ratio":    {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "realistic", "fiction"}},
			},
			"ideogram-v3-text-to-image": {
				"aspect_ratio":    {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_count":    {Enum: []any{1, 2, 3, 4}},
				"rendering_speed": {Enum: []any{"turbo", "balanced", "quality"}},
				"style":           {Enum: []any{"auto", "general", "realistic", "design"}},
			},
		},
	},
	"imagen-4/remix-image": {
		Models: []string{"imagen-4-pro-remix-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"imagen-4-pro-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9", "auto"}},
				"output_format":     {Enum: []any{"png", "jpg"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"imagen-4/text-to-image": {
		Models: []string{"imagen-4", "imagen-4-fast", "imagen-4-ultra"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"imagen-4": {
				"aspect_ratio": {Enum: []any{"1:1", "16:9", "9:16", "3:4", "4:3"}},
			},
			"imagen-4-fast": {
				"aspect_ratio": {Enum: []any{"1:1", "16:9", "9:16", "3:4", "4:3"}},
				"output_count": {Enum: []any{1, 2, 3, 4}},
			},
			"imagen-4-ultra": {
				"aspect_ratio": {Enum: []any{"1:1", "16:9", "9:16", "3:4", "4:3"}},
			},
		},
	},
	"infinitetalk/audio-to-video": {
		Models: []string{"infinitetalk-from-audio"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"infinitetalk-from-audio": {
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
		},
	},
	"kling/avatar": {
		Models: []string{"kling-ai-avatar-pro", "kling-ai-avatar-standard", "kling-ai-avatar-v1-pro", "kling-v1-avatar-standard"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-ai-avatar-pro":      {},
			"kling-ai-avatar-standard": {},
			"kling-ai-avatar-v1-pro":   {},
			"kling-v1-avatar-standard": {},
		},
	},
	"kling/image-to-video": {
		Models: []string{"kling-v2.1-master-image-to-video", "kling-v2.1-pro", "kling-v2.1-standard", "kling-v2.5-turbo-image-to-video-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-v2.1-master-image-to-video": {
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.1-pro": {
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.1-standard": {
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.5-turbo-image-to-video-pro": {
				"duration_seconds": {Enum: []any{5, 10}},
			},
		},
	},
	"kling/motion-control": {
		Models: []string{"kling-3.0"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-3.0": {
				"background_source":     {Enum: []any{"video", "image"}},
				"character_orientation": {Enum: []any{"video", "image"}},
				"output_resolution":     {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"kling/text-to-video": {
		Models: []string{"kling-3.0", "kling-v2.1-master-text-to-video", "kling-v2.5-turbo-text-to-video-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-3.0": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p", "4k"}},
			},
			"kling-v2.1-master-text-to-video": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.5-turbo-text-to-video-pro": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
			},
		},
	},
	"luma/modify-video": {
		Models: []string{"luma-modify-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"luma-modify-video": {},
		},
	},
	"nano-banana/edit-image": {
		Models: []string{"nano-banana-edit"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"nano-banana-edit": {
				"aspect_ratio":  {Enum: []any{"1:1", "9:16", "16:9", "3:4", "4:3", "3:2", "2:3", "5:4", "4:5", "21:9", "auto"}},
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"nano-banana/text-to-image": {
		Models: []string{"nano-banana", "nano-banana-2", "nano-banana-2-lite", "nano-banana-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"nano-banana": {
				"aspect_ratio":  {Enum: []any{"1:1", "9:16", "16:9", "3:4", "4:3", "3:2", "2:3", "5:4", "4:5", "21:9", "auto"}},
				"output_format": {Enum: []any{"png", "jpeg", "jpg"}},
			},
			"nano-banana-2": {
				"aspect_ratio":      {Enum: []any{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9", "auto"}},
				"output_format":     {Enum: []any{"png", "jpeg", "jpg"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
			"nano-banana-2-lite": {
				"aspect_ratio": {Enum: []any{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9", "auto"}},
			},
			"nano-banana-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9", "auto"}},
				"output_format":     {Enum: []any{"png", "jpeg", "jpg"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"qwen-2/edit-image": {
		Models: []string{"qwen-2-edit-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-2-edit-image": {
				"aspect_ratio":  {Enum: []any{"1:1", "2:3", "3:2", "3:4", "4:3", "9:16", "16:9", "21:9"}},
				"output_format": {Enum: []any{"jpeg", "png"}},
			},
		},
	},
	"qwen-2/remix-image": {
		Models: []string{"qwen-2-remix-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-2-remix-image": {
				"acceleration":  {Enum: []any{"none", "regular", "high"}},
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"qwen-2/text-to-image": {
		Models: []string{"qwen-2-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-2-text-to-image": {
				"aspect_ratio":  {Enum: []any{"1:1", "3:4", "4:3", "9:16", "16:9"}},
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"recraft/remove-background": {
		Models: []string{"recraft-remove-background"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"recraft-remove-background": {},
		},
	},
	"recraft/upscale-image": {
		Models: []string{"recraft-crisp-upscale"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"recraft-crisp-upscale": {},
		},
	},
	"runway-aleph/edit-video": {
		Models: []string{"runway-aleph"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"runway-aleph": {
				"aspect_ratio": {Enum: []any{"16:9", "9:16", "4:3", "3:4", "1:1", "21:9"}},
			},
		},
	},
	"runway/extend-video": {
		Models: []string{"runway"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"runway": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"runway/text-to-video": {
		Models: []string{"runway"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"runway": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"duration_seconds":  {Enum: []any{5, 10}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"seedance/text-to-video": {
		Models: []string{"seedance-1.5-pro", "seedance-2.0", "seedance-2.0-fast", "seedance-v1-lite", "seedance-v1-pro", "seedance-v1-pro-fast"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedance-1.5-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_resolution": {Enum: []any{"480p", "720p", "1080p"}},
			},
			"seedance-2.0": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "auto"}},
				"output_resolution": {Enum: []any{"480p", "720p", "1080p"}},
			},
			"seedance-2.0-fast": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "auto"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"seedance-v1-lite": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "9:21"}},
				"duration_seconds":  {Enum: []any{5, 10}},
				"output_resolution": {Enum: []any{"480p", "720p", "1080p"}},
			},
			"seedance-v1-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"duration_seconds":  {Enum: []any{5, 10}},
				"output_resolution": {Enum: []any{"480p", "720p", "1080p"}},
			},
			"seedance-v1-pro-fast": {
				"duration_seconds":  {Enum: []any{5, 10}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"seedream/edit-image": {
		Models: []string{"seedream-4.5-edit", "seedream-5-lite-edit", "seedream-v4-edit"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedream-4.5-edit": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-5-lite-edit": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-v4-edit": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "3:2", "2:3", "16:9", "9:16", "21:9"}},
				"output_count":      {Enum: []any{1, 2, 3, 4, 5, 6}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"seedream/text-to-image": {
		Models: []string{"seedream-4.5-text-to-image", "seedream-5-lite-text-to-image", "seedream-v4-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedream-4.5-text-to-image": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-5-lite-text-to-image": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-v4-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "3:2", "2:3", "16:9", "9:16", "21:9"}},
				"output_count":      {Enum: []any{1, 2, 3, 4, 5, 6}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"suno/add-instrumental": {
		Models: []string{"suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4.5-plus": {
				"vocal_gender": {Enum: []any{"male", "female"}},
			},
			"suno-v5": {
				"vocal_gender": {Enum: []any{"male", "female"}},
			},
			"suno-v5.5": {
				"vocal_gender": {Enum: []any{"male", "female"}},
			},
		},
	},
	"suno/add-vocals": {
		Models: []string{"suno-v4.5-plus", "suno-v5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4.5-plus": {
				"vocal_gender": {Enum: []any{"male", "female"}},
			},
			"suno-v5": {
				"vocal_gender": {Enum: []any{"male", "female"}},
			},
		},
	},
	"suno/boost-style": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/check-voice": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/convert-audio": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/cover-audio": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-all", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-all": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-plus": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
		},
	},
	"suno/create-mashup": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-all", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-all": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-plus": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
		},
	},
	"suno/extend-music": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-all", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
			"suno-v4.5": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
			"suno-v4.5-all": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
			"suno-v4.5-plus": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
			"suno-v5": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
			"suno-v5.5": {
				"parameter_mode": {Enum: []any{"source", "custom"}},
				"persona_type":   {Enum: []any{"style", "voice"}},
				"vocal_gender":   {Enum: []any{"male", "female"}},
			},
		},
	},
	"suno/generate-artwork": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/generate-lyrics": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/generate-midi": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/generate-persona": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/generate-voice": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"singer_skill_level": {Enum: []any{"beginner", "intermediate", "advanced", "professional"}},
			},
		},
	},
	"suno/get-timestamped-lyrics": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/regenerate-validation-phrase": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/replace-section": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/separate-audio-stems": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/text-to-music": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-all", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-all": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-plus": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5.5": {
				"persona_type": {Enum: []any{"style", "voice"}},
				"vocal_gender": {Enum: []any{"male", "female"}},
				"vocal_mode":   {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
		},
	},
	"suno/text-to-sound": {
		Models: []string{"suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v5": {
				"sound_key": {Enum: []any{"Cm", "C#m", "Dm", "D#m", "Em", "Fm", "F#m", "Gm", "G#m", "Am", "A#m", "Bm", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}},
			},
			"suno-v5.5": {
				"sound_key": {Enum: []any{"Cm", "C#m", "Dm", "D#m", "Em", "Fm", "F#m", "Gm", "G#m", "Am", "A#m", "Bm", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}},
			},
		},
	},
	"suno/visualize-music": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"suno/voice-to-validation-phrase": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"language": {Enum: []any{"en", "zh", "es", "fr", "pt", "de", "ja", "ko", "hi", "ru"}},
			},
		},
	},
	"topaz/upscale-image": {
		Models: []string{"topaz-upscale-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"topaz-upscale-image": {
				"upscale_factor": {Enum: []any{1, 2, 4, 8}},
			},
		},
	},
	"topaz/upscale-video": {
		Models: []string{"topaz-upscale-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"topaz-upscale-video": {
				"upscale_factor": {Enum: []any{1, 2, 4}},
			},
		},
	},
	"veo-3-1/extend-video": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"veo-3-1/text-to-video": {
		Models: []string{"veo-3.1", "veo-3.1-fast"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"veo-3.1": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "auto"}},
				"duration_seconds": {Enum: []any{4, 6, 8}},
				"input_mode":       {Enum: []any{"text", "first_and_last_frames", "reference"}},
			},
			"veo-3.1-fast": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "auto"}},
				"duration_seconds": {Enum: []any{4, 6, 8}},
				"input_mode":       {Enum: []any{"text", "first_and_last_frames", "reference"}},
			},
		},
	},
	"veo-3-1/upscale-video": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"output_resolution": {Enum: []any{"1080p", "4k"}},
			},
		},
	},
	"wan/animate": {
		Models: []string{"wan-2.2-animate-move", "wan-2.2-animate-replace"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.2-animate-move": {
				"output_resolution": {Enum: []any{"480p", "580p", "720p"}},
			},
			"wan-2.2-animate-replace": {
				"output_resolution": {Enum: []any{"480p", "580p", "720p"}},
			},
		},
	},
	"wan/edit-video": {
		Models: []string{"wan-2.6-edit-video", "wan-2.6-flash-edit-video", "wan-2.7-edit-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.6-edit-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.6-flash-edit-video": {},
			"wan-2.7-edit-video": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"wan/image-to-video": {
		Models: []string{"wan-2.2-a14b-image-to-video-turbo", "wan-2.5-image-to-video", "wan-2.6-flash-image-to-video", "wan-2.6-image-to-video", "wan-2.7-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.2-a14b-image-to-video-turbo": {
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"wan-2.5-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.6-flash-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.6-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.7-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"wan/speech-to-video": {
		Models: []string{"wan-2.2-a14b-speech-to-video-turbo"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.2-a14b-speech-to-video-turbo": {
				"output_resolution": {Enum: []any{"480p", "580p", "720p"}},
			},
		},
	},
	"wan/text-to-image": {
		Models: []string{"wan-2.7-image", "wan-2.7-image-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.7-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "16:9", "4:3", "21:9", "3:4", "9:16", "8:1", "1:8"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
			"wan-2.7-image-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "16:9", "4:3", "21:9", "3:4", "9:16", "8:1", "1:8"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"wan/text-to-video": {
		Models: []string{"wan-2.2-a14b-text-to-video-turbo", "wan-2.5-text-to-video", "wan-2.6-text-to-video", "wan-2.7-r2v", "wan-2.7-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"wan-2.2-a14b-text-to-video-turbo": {
				"output_resolution": {Enum: []any{"480p", "580p", "720p"}},
			},
			"wan-2.5-text-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.6-text-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.7-r2v": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"wan-2.7-text-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"z-image/text-to-image": {
		Models: []string{"z-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"z-image": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16"}},
			},
		},
	},
}
