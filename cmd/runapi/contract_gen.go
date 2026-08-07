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
	"fish-audio/text-to-speech": {
		Models: []string{"s1", "s2-pro", "s2.1-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"s1": {
				"bitrate_kbps":   {Enum: []any{64, 128, 192}},
				"output_format":  {Enum: []any{"mp3", "wav"}},
				"sample_rate_hz": {Enum: []any{8000, 16000, 24000, 32000, 44100}},
			},
			"s2-pro": {
				"bitrate_kbps":   {Enum: []any{64, 128, 192}},
				"output_format":  {Enum: []any{"mp3", "wav"}},
				"sample_rate_hz": {Enum: []any{8000, 16000, 24000, 32000, 44100}},
			},
			"s2.1-pro": {
				"bitrate_kbps":   {Enum: []any{64, 128, 192}},
				"output_format":  {Enum: []any{"mp3", "wav"}},
				"sample_rate_hz": {Enum: []any{8000, 16000, 24000, 32000, 44100}},
			},
		},
	},
	"flux-2/remix-image": {
		Models: []string{"flux-2-flex-remix-image", "flux-2-max-remix-image", "flux-2-pro-remix-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-2-flex-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 8},
			},
			"flux-2-max-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count":      {Enum: []any{1}},
				"output_resolution": {Enum: []any{"1k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 1},
			},
			"flux-2-pro-remix-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 8},
			},
		},
	},
	"flux-2/text-to-image": {
		Models: []string{"flux-2-flex-text-to-image", "flux-2-max-text-to-image", "flux-2-pro-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-2-flex-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
			"flux-2-max-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count":      {Enum: []any{1}},
				"output_resolution": {Enum: []any{"1k"}},
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
	"flux/remix-image": {
		Models: []string{"flux-dev", "flux-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-dev": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count": {Enum: []any{1}},
			},
			"flux-pro": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count": {Enum: []any{1}},
			},
		},
	},
	"flux/text-to-image": {
		Models: []string{"flux-2-klein", "flux-dev", "flux-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"flux-2-klein": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count": {Enum: []any{1}},
			},
			"flux-dev": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count": {Enum: []any{1}},
			},
			"flux-pro": {
				"aspect_ratio": {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3"}},
				"output_count": {Enum: []any{1}},
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
		Models: []string{"gemini-omni-flash-preview", "gemini-omni-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gemini-omni-flash-preview": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16"}},
				"output_resolution": {Enum: []any{"720p"}},
			},
			"gemini-omni-text-to-video": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16"}},
				"audio_ids":            {MaxItems: 3},
				"character_ids":        {MaxItems: 3},
				"duration_seconds":     {Enum: []any{4, 6, 8, 10}},
				"output_resolution":    {Enum: []any{"720p", "1080p", "4k"}},
				"reference_image_urls": {MaxItems: 7},
				"video_list":           {MaxItems: 1},
			},
		},
	},
	"gemini-tts/text-to-speech": {
		Models: []string{"gemini-2.5-pro-tts", "gemini-3.1-flash-tts"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gemini-2.5-pro-tts": {
				"dialogue_turns": {MinItems: 1},
				"speakers":       {MinItems: 1},
			},
			"gemini-3.1-flash-tts": {
				"dialogue_turns": {MinItems: 1},
				"speakers":       {MinItems: 1},
			},
		},
	},
	"gpt-4o-image/text-to-image": {
		Models: []string{"gpt-4o-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-4o-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "3:2", "2:3"}},
				"output_count":      {Enum: []any{1, 2, 4}},
				"source_image_urls": {MaxItems: 5},
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
		Models: []string{"grok-imagine-image-to-video", "grok-imagine-video-1.5-fast", "grok-imagine-video-1.5-preview"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-image-to-video": {
				"aspect_ratio":      {Enum: []any{"2:3", "3:2", "1:1", "16:9", "9:16"}},
				"motion_style":      {Enum: []any{"fun", "normal", "spicy"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"grok-imagine-video-1.5-fast": {
				"aspect_ratio":      {Enum: []any{"1:1", "16:9", "9:16", "3:2", "2:3"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"grok-imagine-video-1.5-preview": {
				"aspect_ratio":         {Enum: []any{"1:1", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution":    {Enum: []any{"480p", "720p", "1080p"}},
				"reference_image_urls": {MaxItems: 6},
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
		Models: []string{"grok-imagine-text-to-video", "grok-imagine-video-1.5-fast", "grok-imagine-video-1.5-preview"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"grok-imagine-text-to-video": {
				"aspect_ratio":      {Enum: []any{"2:3", "3:2", "1:1", "16:9", "9:16"}},
				"motion_style":      {Enum: []any{"fun", "normal", "spicy"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"grok-imagine-video-1.5-fast": {
				"aspect_ratio":      {Enum: []any{"1:1", "16:9", "9:16", "3:2", "2:3"}},
				"output_resolution": {Enum: []any{"480p", "720p"}},
			},
			"grok-imagine-video-1.5-preview": {
				"aspect_ratio":         {Enum: []any{"1:1", "16:9", "9:16", "3:2", "2:3", "auto"}},
				"output_resolution":    {Enum: []any{"480p", "720p", "1080p"}},
				"reference_image_urls": {MaxItems: 7},
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
				"audio_setting":        {Enum: []any{"auto", "original"}},
				"output_resolution":    {Enum: []any{"720p", "1080p"}},
				"reference_image_urls": {MaxItems: 5},
			},
		},
	},
	"happyhorse/image-to-video": {
		Models: []string{"happyhorse-1.0-i2v", "happyhorse-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"happyhorse-1.0-i2v": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"happyhorse-image-to-video": {
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"happyhorse/text-to-video": {
		Models: []string{"happyhorse-1.0-r2v", "happyhorse-1.0-t2v", "happyhorse-character", "happyhorse-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"happyhorse-1.0-r2v": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution":    {Enum: []any{"720p", "1080p"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 9},
			},
			"happyhorse-1.0-t2v": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
			"happyhorse-character": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "1:1", "4:3", "3:4"}},
				"output_resolution":    {Enum: []any{"720p", "1080p"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 9},
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
				"source_image_urls": {MinItems: 1, MaxItems: 8},
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
				"aspect_ratio": {Enum: []any{"1:1", "16:9", "9:16", "3:4", "4:3", "auto"}},
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
	"kling/extend-video": {
		Models: []string{"kling-v2.5-turbo-image-to-video-pro", "kling-v2.5-turbo-text-to-video-pro"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-v2.5-turbo-image-to-video-pro": {
				"mode": {Enum: []any{"std", "pro"}},
			},
			"kling-v2.5-turbo-text-to-video-pro": {
				"mode": {Enum: []any{"std", "pro"}},
			},
		},
	},
	"kling/image-to-video": {
		Models: []string{"kling-o1", "kling-v2.1-master-image-to-video", "kling-v2.1-pro", "kling-v2.1-standard", "kling-v2.5-turbo-image-to-video-pro", "kling-v2.6", "kling-v3-omni", "kling-v3-turbo-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-o1": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":     {Enum: []any{5}},
				"mode":                 {Enum: []any{"std", "pro"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 7},
				"reference_video_type": {Enum: []any{"base", "feature"}},
			},
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
			"kling-v2.6": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
				"mode":             {Enum: []any{"std", "pro"}},
			},
			"kling-v3-omni": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p", "4k"}},
			},
			"kling-v3-turbo-image-to-video": {
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"kling/motion-control": {
		Models: []string{"kling-3.0", "kling-v2.6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-3.0": {
				"background_source":     {Enum: []any{"video", "image"}},
				"character_orientation": {Enum: []any{"video", "image"}},
				"output_resolution":     {Enum: []any{"720p", "1080p"}},
			},
			"kling-v2.6": {
				"character_orientation": {Enum: []any{"video", "image"}},
				"output_resolution":     {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"kling/text-to-video": {
		Models: []string{"kling-3.0", "kling-o1", "kling-v2.1-master-text-to-video", "kling-v2.5-turbo-text-to-video-pro", "kling-v2.6", "kling-v3-omni", "kling-v3-turbo-text-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"kling-3.0": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p", "4k"}},
			},
			"kling-o1": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":     {Enum: []any{5}},
				"mode":                 {Enum: []any{"std", "pro"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 7},
				"reference_video_type": {Enum: []any{"base", "feature"}},
			},
			"kling-v2.1-master-text-to-video": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.5-turbo-text-to-video-pro": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
			},
			"kling-v2.6": {
				"aspect_ratio":     {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds": {Enum: []any{5, 10}},
				"mode":             {Enum: []any{"std", "pro"}},
			},
			"kling-v3-omni": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p", "4k"}},
			},
			"kling-v3-turbo-text-to-video": {
				"aspect_ratio":      {Enum: []any{"16:9", "9:16", "1:1"}},
				"duration_seconds":  {Enum: []any{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"luma/modify-video": {
		Models: []string{"luma-modify-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"luma-modify-video": {},
		},
	},
	"midjourney/edit-image": {
		Models: []string{"midjourney-edit-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"midjourney-edit-image": {},
		},
	},
	"midjourney/extend-video": {
		Models: []string{"midjourney-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"midjourney-image-to-video": {},
		},
	},
	"midjourney/get-seed": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"midjourney/image-to-prompt": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"midjourney/image-to-video": {
		Models: []string{"midjourney-image-to-video"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"midjourney-image-to-video": {
				"output_resolution": {Enum: []any{"480p"}},
			},
		},
	},
	"midjourney/shorten-prompt": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
		},
	},
	"midjourney/text-to-image": {
		Models: []string{"midjourney-v8.1"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"midjourney-v8.1": {},
		},
	},
	"minimax-h3/image-to-video": {
		Models: []string{"minimax-h3"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"minimax-h3": {
				"duration_seconds":  {Enum: []any{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution": {Enum: []any{"768p", "2k"}},
			},
		},
	},
	"minimax-h3/text-to-video": {
		Models: []string{"minimax-h3"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"minimax-h3": {
				"aspect_ratio":         {Enum: []any{"adaptive", "21:9", "16:9", "4:3", "1:1", "3:4", "9:16"}},
				"duration_seconds":     {Enum: []any{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				"output_resolution":    {Enum: []any{"768p", "2k"}},
				"reference_audio_urls": {MaxItems: 3},
				"reference_image_urls": {MaxItems: 9},
				"reference_video_urls": {MaxItems: 3},
			},
		},
	},
	"nano-banana/edit-image": {
		Models: []string{"nano-banana-2-lite", "nano-banana-edit"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"nano-banana-2-lite": {
				"aspect_ratio":      {Enum: []any{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9", "auto"}},
				"source_image_urls": {MinItems: 1, MaxItems: 10},
			},
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
				"aspect_ratio":         {Enum: []any{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9", "auto"}},
				"reference_image_urls": {MaxItems: 10},
			},
			"nano-banana-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9", "auto"}},
				"output_format":     {Enum: []any{"png", "jpeg", "jpg"}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
			},
		},
	},
	"omnihuman/audio-to-video": {
		Models: []string{"omnihuman-1.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"omnihuman-1.5": {
				"mask_urls":         {MaxItems: 5},
				"output_resolution": {Enum: []any{"720p", "1080p"}},
			},
		},
	},
	"omnihuman/human-identification": {
		Models: []string{"omnihuman-1.5-human-identification"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"omnihuman-1.5-human-identification": {},
		},
	},
	"omnihuman/subject-detection": {
		Models: []string{"omnihuman-1.5-subject-detection"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"omnihuman-1.5-subject-detection": {},
		},
	},
	"openai-transcription/speech-to-text": {
		Models: []string{"gpt-transcribe", "whisper-1"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"gpt-transcribe": {
				"model":           {Enum: []any{"gpt-transcribe"}},
				"response_format": {Enum: []any{"json", "text"}},
			},
			"whisper-1": {
				"model":                   {Enum: []any{"whisper-1"}},
				"response_format":         {Enum: []any{"json", "text", "srt", "verbose_json", "vtt"}},
				"timestamp_granularities": {MaxItems: 2},
			},
		},
	},
	"openai-tts/text-to-speech": {
		Models: []string{"tts-1", "tts-1-hd"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"tts-1":    {},
			"tts-1-hd": {},
		},
	},
	"pixverse/edit-video": {
		Models: []string{"pixverse-v6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"pixverse-v6": {
				"aspect_ratio":         {Enum: []any{"16:9", "4:3", "1:1", "3:4", "9:16", "2:3", "3:2", "21:9"}},
				"enable_audio":         {Enum: []any{true, false}},
				"output_resolution":    {Enum: []any{"360p", "540p", "720p", "1080p"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 7},
			},
		},
	},
	"pixverse/extend-video": {
		Models: []string{"pixverse-v6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"pixverse-v6": {
				"enable_audio":      {Enum: []any{true, false}},
				"output_resolution": {Enum: []any{"360p", "540p", "720p", "1080p"}},
			},
		},
	},
	"pixverse/image-to-video": {
		Models: []string{"pixverse-v6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"pixverse-v6": {
				"enable_audio":      {Enum: []any{true, false}},
				"output_resolution": {Enum: []any{"360p", "540p", "720p", "1080p"}},
			},
		},
	},
	"pixverse/text-to-video": {
		Models: []string{"pixverse-v6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"pixverse-v6": {
				"aspect_ratio":      {Enum: []any{"16:9", "4:3", "1:1", "3:4", "9:16", "2:3", "3:2", "21:9"}},
				"enable_audio":      {Enum: []any{true, false}},
				"output_resolution": {Enum: []any{"360p", "540p", "720p", "1080p"}},
			},
		},
	},
	"pixverse/transition-video": {
		Models: []string{"pixverse-v6"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"pixverse-v6": {
				"enable_audio":      {Enum: []any{true, false}},
				"output_resolution": {Enum: []any{"360p", "540p", "720p", "1080p"}},
			},
		},
	},
	"producer/text-to-music": {
		Models: []string{"fuzz-0.8", "fuzz-1.0", "fuzz-1.0-pro", "fuzz-1.1", "fuzz-1.1-pro", "fuzz-2.0", "fuzz-2.0-pro", "fuzz-2.0-raw"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"fuzz-0.8": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-1.0": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-1.0-pro": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-1.1": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-1.1-pro": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-2.0": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-2.0-pro": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
			},
			"fuzz-2.0-raw": {
				"vocal_mode": {Enum: []any{"exact_lyrics", "instrumental"}},
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
	"qwen-2/text-to-image": {
		Models: []string{"qwen-2-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-2-text-to-image": {
				"aspect_ratio":  {Enum: []any{"1:1", "3:4", "4:3", "9:16", "16:9"}},
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"qwen-3/edit-image": {
		Models: []string{"qwen-3-edit-image", "qwen-3-pro-edit-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-3-edit-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "3:2", "2:3", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 3},
			},
			"qwen-3-pro-edit-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "3:2", "2:3", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 3},
			},
		},
	},
	"qwen-3/text-to-image": {
		Models: []string{"qwen-3-pro-text-to-image", "qwen-3-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-3-pro-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "3:2", "2:3", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
			"qwen-3-text-to-image": {
				"aspect_ratio":      {Enum: []any{"1:1", "3:2", "2:3", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_resolution": {Enum: []any{"1k", "2k"}},
			},
		},
	},
	"qwen-image/edit-image": {
		Models: []string{"qwen-image-edit-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-image-edit-image": {
				"aspect_ratio":  {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"qwen-image/remix-image": {
		Models: []string{"qwen-image-remix-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-image-remix-image": {
				"output_format": {Enum: []any{"png", "jpeg"}},
			},
		},
	},
	"qwen-image/text-to-image": {
		Models: []string{"qwen-image-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"qwen-image-text-to-image": {
				"aspect_ratio":  {Enum: []any{"1:1", "3:4", "9:16", "4:3", "16:9"}},
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
		Models: []string{"seedance-1.5-pro", "seedance-2-mini", "seedance-2.0", "seedance-2.0-fast", "seedance-v1-pro", "seedance-v1-pro-fast"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedance-1.5-pro": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9"}},
				"output_resolution": {Enum: []any{"480p", "720p", "1080p"}},
				"source_image_urls": {MaxItems: 2},
			},
			"seedance-2-mini": {
				"aspect_ratio":         {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "auto"}},
				"output_resolution":    {Enum: []any{"480p", "720p"}},
				"reference_audio_urls": {MaxItems: 3},
				"reference_image_urls": {MaxItems: 9},
				"reference_video_urls": {MaxItems: 3},
			},
			"seedance-2.0": {
				"aspect_ratio":         {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "auto"}},
				"output_resolution":    {Enum: []any{"480p", "720p", "1080p", "4k"}},
				"reference_audio_urls": {MaxItems: 3},
				"reference_image_urls": {MaxItems: 9},
				"reference_video_urls": {MaxItems: 3},
			},
			"seedance-2.0-fast": {
				"aspect_ratio":         {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "auto"}},
				"output_resolution":    {Enum: []any{"480p", "720p"}},
				"reference_audio_urls": {MaxItems: 3},
				"reference_image_urls": {MaxItems: 9},
				"reference_video_urls": {MaxItems: 3},
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
		Models: []string{"seedream-4.5-edit", "seedream-5-lite-edit", "seedream-5-pro-edit", "seedream-v4-edit"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedream-4.5-edit": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality":    {Enum: []any{"basic", "high"}},
				"source_image_urls": {MinItems: 1, MaxItems: 14},
			},
			"seedream-5-lite-edit": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_quality":    {Enum: []any{"basic", "high"}},
				"source_image_urls": {MinItems: 1, MaxItems: 14},
			},
			"seedream-5-pro-edit": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_format":     {Enum: []any{"png", "jpeg"}},
				"output_quality":    {Enum: []any{"basic", "high"}},
				"source_image_urls": {MinItems: 1, MaxItems: 10},
			},
			"seedream-v4-edit": {
				"aspect_ratio":      {Enum: []any{"1:1", "4:3", "3:4", "3:2", "2:3", "16:9", "9:16", "21:9"}},
				"output_count":      {Enum: []any{1, 2, 3, 4, 5, 6}},
				"output_resolution": {Enum: []any{"1k", "2k", "4k"}},
				"source_image_urls": {MinItems: 1, MaxItems: 10},
			},
		},
	},
	"seedream/text-to-image": {
		Models: []string{"seedream-4.5-text-to-image", "seedream-5-lite-text-to-image", "seedream-5-pro-text-to-image", "seedream-v4-text-to-image"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"seedream-4.5-text-to-image": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-5-lite-text-to-image": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_format":  {Enum: []any{"png", "jpeg"}},
				"output_quality": {Enum: []any{"basic", "high"}},
			},
			"seedream-5-pro-text-to-image": {
				"aspect_ratio":   {Enum: []any{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}},
				"output_format":  {Enum: []any{"png", "jpeg"}},
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
	"suno/add-samples": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4":        {},
			"suno-v4.5":      {},
			"suno-v4.5-plus": {},
			"suno-v5":        {},
			"suno-v5.5":      {},
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
	"suno/blend-lyrics": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {},
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
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5": {
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-all": {
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v4.5-plus": {
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5": {
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
			},
			"suno-v5.5": {
				"persona_type":    {Enum: []any{"style", "voice"}},
				"upload_url_list": {MinItems: 2, MaxItems: 2},
				"vocal_gender":    {Enum: []any{"male", "female"}},
				"vocal_mode":      {Enum: []any{"auto_lyrics", "exact_lyrics", "instrumental"}},
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
	"suno/remaster-audio": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4":        {},
			"suno-v4.5":      {},
			"suno-v4.5-plus": {},
			"suno-v5":        {},
			"suno-v5.5":      {},
		},
	},
	"suno/replace-section": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"model": {Enum: []any{"suno-v4", "suno-v4.5", "suno-v4.5-all", "suno-v4.5-plus", "suno-v5", "suno-v5.5"}},
			},
		},
	},
	"suno/separate-audio-stems": {
		Models: []string{},
		FieldsByModel: map[string]map[string]generatedContractField{
			"_": {
				"stem_name": {Enum: []any{"Lead Vocal", "Drum Kit", "Kick", "Snare", "Risers", "Bass", "Backing Vocals", "Piano", "Electric Guitar", "Percussion", "String Section", "Synth", "Acoustic Guitar", "Sound Effects", "Synth Pad", "Synth Bass", "Guitar", "Brass Section", "Organ", "Electronic Drum Kit", "Lead Electric Guitar", "Synth Keys", "Rhythm Electric Guitar", "Electric Piano", "Upright Bass", "Keyboards", "Distorted Electric Guitar", "Synth Strings", "Synth Lead", "Woodwinds", "Rhythm Acoustic Guitar", "Flute", "Harp", "Tambourine", "Trumpet", "Arpeggiator", "Accordion", "Fiddle", "Pedal Steel Guitar", "Synth Voice", "Violin", "Digital Piano", "Synth Brass", "Mandolin", "Choir", "Banjo", "Bells", "Clarinet", "Tenor Saxophone", "Trombone", "Shaker", "French Horn", "Glockenspiel", "Electric Bass", "Cello", "Timpani", "Harmonica", "Marimba", "Vibraphone", "Lap Steel Guitar", "Saxophone", "Orchestra", "Horns", "Cymbals", "Hand Clap", "Oboe", "Celesta", "Congas", "Drone", "Alto Saxophone", "Double Bass", "Ukulele", "Harpsichord", "Baritone Saxophone", "Xylophone", "Tuba", "Bass Guitar", "Whistle", "Lead Guitar", "Rhodes", "808", "Bongos", "Bassoon", "Cowbell", "Viola", "Sitar", "Steel Drums", "Piccolo", "Theremin", "Bagpipes", "Hi-Hat", "Music Box", "Melodica", "Tabla", "Koto", "Djembe", "Taiko", "Didgeridoo"}},
				"type":      {Enum: []any{"separate_vocal", "split_stem", "split_stem_advanced"}},
			},
		},
	},
	"suno/stitch-audio": {
		Models: []string{"suno-v4", "suno-v4.5", "suno-v4.5-plus", "suno-v5", "suno-v5.5"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"suno-v4":        {},
			"suno-v4.5":      {},
			"suno-v4.5-plus": {},
			"suno-v5":        {},
			"suno-v5.5":      {},
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
		Models: []string{"veo-3.1", "veo-3.1-fast", "veo-3.1-lite"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"veo-3.1": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "auto"}},
				"duration_seconds":     {Enum: []any{4, 6, 8}},
				"input_mode":           {Enum: []any{"text", "first_and_last_frames", "reference"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 3},
			},
			"veo-3.1-fast": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16", "auto"}},
				"duration_seconds":     {Enum: []any{4, 6, 8}},
				"input_mode":           {Enum: []any{"text", "first_and_last_frames", "reference"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 3},
			},
			"veo-3.1-lite": {
				"aspect_ratio":         {Enum: []any{"16:9", "9:16"}},
				"duration_seconds":     {Enum: []any{4, 6, 8}},
				"input_mode":           {Enum: []any{"text", "first_and_last_frames", "reference"}},
				"reference_image_urls": {MinItems: 1, MaxItems: 3},
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
	"volcengine-lip-sync/lip-sync-video": {
		Models: []string{"volcengine-lip-sync"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"volcengine-lip-sync": {
				"mode": {Enum: []any{"lite", "basic"}},
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
				"duration_seconds":  {Enum: []any{5, 10, 15}},
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
