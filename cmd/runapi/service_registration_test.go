package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestSynchronousRawTextIsWrittenWithoutJSONEncoding(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}

	if err := c.writeSynchronousResponse(json.RawMessage("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n")); err != nil {
		t.Fatal(err)
	}

	if got := c.stdout.(*bytes.Buffer).String(); got != "WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n" {
		t.Fatalf("expected exact raw response, got %q", got)
	}
}

func helpHasField(output, field string) bool {
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Fields(line)
		if len(parts) > 0 && parts[0] == field {
			return true
		}
	}

	return false
}

func helpFieldCount(output, field string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Fields(line)
		if len(parts) > 0 && parts[0] == field {
			count++
		}
	}

	return count
}

func TestGeminiOmniServiceCommandIsRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"gemini-omni", "create-audio", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "audio_id") {
		t.Fatalf("expected Gemini Omni create-audio help to include audio_id, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	cmd = c.command()
	cmd.SetArgs([]string{"gemini-omni", "create-character", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "descriptions") || !helpHasField(output, "reference_image_url") || helpHasField(output, "image_urls") {
		t.Fatalf("expected Gemini Omni create-character help to include character fields, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	cmd = c.command()
	cmd.SetArgs([]string{"gemini-omni", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "model") || !strings.Contains(output, "gemini-omni-flash-preview") || !strings.Contains(output, "duration_seconds") || !strings.Contains(output, "character_ids") || !helpHasField(output, "reference_image_urls") || helpHasField(output, "image_urls") {
		t.Fatalf("expected Gemini Omni text-to-video help to include video fields, got:\n%s", output)
	}
}

func TestMidjourneyShortenPromptCommandIsRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"midjourney", "shorten-prompt", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "prompt") {
		t.Fatalf("expected Midjourney shorten-prompt help to include prompt, got:\n%s", output)
	}
	if helpHasField(output, "model") || helpHasField(output, "callback_url") {
		t.Fatalf("expected Midjourney shorten-prompt help to expose only its synchronous input, got:\n%s", output)
	}
}

func TestMidjourneyExtendVideoCommandIsRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"midjourney", "extend-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	for _, field := range []string{"source_task_id", "prompt", "callback_url"} {
		if !helpHasField(output, field) {
			t.Fatalf("expected Midjourney extend-video help to include %s, got:\n%s", field, output)
		}
	}
	for _, field := range []string{"model", "video_id", "video_index", "output_resolution"} {
		if helpHasField(output, field) {
			t.Fatalf("expected Midjourney extend-video help to omit %s, got:\n%s", field, output)
		}
	}
}

func TestTTSServiceCommandsAreRegistered(t *testing.T) {
	for _, service := range []string{"openai-tts", "fish-audio"} {
		c := newCLI()
		c.stdout = &bytes.Buffer{}
		c.stderr = &bytes.Buffer{}
		cmd := c.command()
		cmd.SetArgs([]string{service, "text-to-speech", "--help"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		output := c.stdout.(*bytes.Buffer).String()
		if !helpHasField(output, "model") || !helpHasField(output, "text") {
			t.Fatalf("expected %s text-to-speech help to include model and text, got:\n%s", service, output)
		}
		if service == "fish-audio" && (!helpHasField(output, "references") || !helpHasField(output, "output_format") || !helpHasField(output, "sample_rate_hz") || !helpHasField(output, "bitrate_kbps")) {
			t.Fatalf("expected Fish Audio help to include references and output controls, got:\n%s", output)
		}
	}
}

func TestOpenAITranscriptionCommandIsRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	cmd := c.command()
	cmd.SetArgs([]string{"openai-transcription", "speech-to-text", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	for _, field := range []string{"file", "model", "languages", "keywords", "response_format", "timestamp_granularities"} {
		if !helpHasField(output, field) {
			t.Fatalf("expected OpenAI transcription help to include %s, got:\n%s", field, output)
		}
	}
}

func TestHappyHorseServiceCommandIsRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"happyhorse", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "resolution") || !strings.Contains(output, "aspect_ratio") || !strings.Contains(output, "duration_seconds") || !strings.Contains(output, "reference_image_urls") {
		t.Fatalf("expected HappyHorse text-to-video help to include video fields, got:\n%s", output)
	}

	removedAction := strings.Join([]string{"reference", "to", "video"}, "-")
	if strings.Contains(output, removedAction) {
		t.Fatalf("expected HappyHorse text-to-video help not to include removed command, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"happyhorse", "image-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "first_frame_image_url") || !strings.Contains(output, "resolution") || !strings.Contains(output, "duration_seconds") {
		t.Fatalf("expected HappyHorse image-to-video help to include image fields, got:\n%s", output)
	}
	if helpHasField(output, "image_urls") {
		t.Fatalf("expected HappyHorse image-to-video help not to expose image_urls, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"happyhorse", "edit-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_video_url") || !strings.Contains(output, "audio_setting") || !helpHasField(output, "reference_image_urls") {
		t.Fatalf("expected HappyHorse edit-video help to include edit fields, got:\n%s", output)
	}
	if helpHasField(output, "video_url") || helpHasField(output, "reference_image") {
		t.Fatalf("expected HappyHorse edit-video help not to expose provider media fields, got:\n%s", output)
	}
}

func TestWanServiceCommandUsesCanonicalMediaFields(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		expected  []string
		forbidden []string
	}{
		{
			name:      "text-to-video",
			args:      []string{"wan", "text-to-video", "--help"},
			expected:  []string{"multi_shots"},
			forbidden: []string{},
		},
		{
			name:      "image-to-video",
			args:      []string{"wan", "image-to-video", "--help"},
			expected:  []string{"first_frame_image_url", "last_frame_image_url", "source_video_url", "background_audio_url", "multi_shots"},
			forbidden: []string{"image_url", "image_urls", "first_frame_url", "last_frame_url", "first_clip_url", "audio_url"},
		},
		{
			name:      "speech-to-video",
			args:      []string{"wan", "speech-to-video", "--help"},
			expected:  []string{"source_image_url", "source_audio_url"},
			forbidden: []string{"image_url", "audio_url"},
		},
		{
			name:      "animate",
			args:      []string{"wan", "animate", "--help"},
			expected:  []string{"source_image_url", "reference_video_url"},
			forbidden: []string{"image_url", "video_url"},
		},
		{
			name:      "text-to-image",
			args:      []string{"wan", "text-to-image", "--help"},
			expected:  []string{"source_image_urls"},
			forbidden: []string{"input_urls"},
		},
		{
			name:      "edit-video",
			args:      []string{"wan", "edit-video", "--help"},
			expected:  []string{"source_video_url", "source_video_urls", "reference_image_url", "multi_shots"},
			forbidden: []string{"reference_image"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}

			cmd := c.command()
			cmd.SetArgs(tc.args)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			output := c.stdout.(*bytes.Buffer).String()
			for _, field := range tc.expected {
				if !helpHasField(output, field) {
					t.Fatalf("expected Wan %s help to include %s, got:\n%s", tc.name, field, output)
				}
			}
			for _, field := range tc.forbidden {
				if helpHasField(output, field) {
					t.Fatalf("expected Wan %s help to omit stale field %s, got:\n%s", tc.name, field, output)
				}
			}
		})
	}
}

func TestHailuoImageToVideoHelpUsesCanonicalImageFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"hailuo", "image-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "first_frame_image_url") || !helpHasField(output, "last_frame_image_url") {
		t.Fatalf("expected Hailuo image-to-video help to include canonical image fields, got:\n%s", output)
	}
	if helpHasField(output, "image_url") || helpHasField(output, "end_image_url") {
		t.Fatalf("expected Hailuo image-to-video help not to expose roleless image fields, got:\n%s", output)
	}
}

func TestInfinitetalkAudioToVideoHelpUsesCanonicalSourceMediaFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"infinitetalk", "audio-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_image_url") || !helpHasField(output, "source_audio_url") {
		t.Fatalf("expected InfiniteTalk audio-to-video help to include canonical source media fields, got:\n%s", output)
	}
	if helpHasField(output, "image_url") || helpHasField(output, "audio_url") {
		t.Fatalf("expected InfiniteTalk audio-to-video help not to expose roleless media fields, got:\n%s", output)
	}
}

func TestKlingTextToVideoHelpUsesCanonicalFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"kling", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "output_resolution") || !strings.Contains(output, "Accepted values: kling-3.0, kling-o1, kling-v2.1-master-text-to-video, kling-v2.5-turbo-text-to-video-pro, kling-v2.6, kling-v3-omni, kling-v3-turbo-text-to-video.") {
		t.Fatalf("expected Kling text-to-video help to include generated model values, got:\n%s", output)
	}
	if !helpHasField(output, "first_frame_image_url") || !helpHasField(output, "last_frame_image_url") {
		t.Fatalf("expected Kling text-to-video help to include frame media fields, got:\n%s", output)
	}
	if !strings.Contains(output, "duration_seconds") || !strings.Contains(output, "Accepted values by model: kling-3.0, kling-v3-omni, kling-v3-turbo-text-to-video: 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15; kling-o1: 5; kling-v2.1-master-text-to-video, kling-v2.5-turbo-text-to-video-pro, kling-v2.6: 5, 10.") {
		t.Fatalf("expected Kling text-to-video help to include model-specific duration_seconds values, got:\n%s", output)
	}
	if strings.Contains(output, "duration_seconds          string     optional; duration in seconds. Accepted values: 5, 10.") {
		t.Fatalf("expected Kling text-to-video help not to advertise action-level duration_seconds values, got:\n%s", output)
	}
	if !helpHasField(output, "mode") || !strings.Contains(output, "Accepted values by model: kling-o1, kling-v2.6: std, pro.") {
		t.Fatalf("expected Kling text-to-video help to include Kling O1 and 2.6 mode values, got:\n%s", output)
	}
	if !strings.Contains(output, "kling-v3-omni: 720p, 1080p, 4k") {
		t.Fatalf("expected Kling text-to-video help to include Kling V3 Omni output resolution values, got:\n%s", output)
	}
	if helpHasField(output, "image_urls") {
		t.Fatalf("expected Kling text-to-video help not to expose roleless media fields, got:\n%s", output)
	}
}

func TestGeneratedContractHelpIncludesArrayItemCounts(t *testing.T) {
	common := generatedContractAction{
		Models: []string{"m-a", "m-b"},
		FieldsByModel: map[string]map[string]generatedContractField{
			"m-a": {"reference_image_urls": {"min_items": 1, "max_items": 3}},
			"m-b": {"reference_image_urls": {"min_items": 1, "max_items": 3}},
		},
	}
	if got := generatedContractHelpSentenceFor(common, "reference_image_urls"); got != "Item count: 1-3." {
		t.Fatalf("got %q", got)
	}

	divergent := common
	divergent.FieldsByModel["m-a"] = map[string]generatedContractField{
		"reference_image_urls": {"min_items": 1, "max_items": 2},
	}
	if got := generatedContractHelpSentenceFor(divergent, "reference_image_urls"); got != "Item count by model: m-a: 1-2; m-b: 1-3." {
		t.Fatalf("got %q", got)
	}
}

func TestGeneratedContractHelpIncludesConditionalRules(t *testing.T) {
	cases := []struct {
		args []string
		want []string
	}{
		{
			args: []string{"runway", "text-to-video", "--help"},
			want: []string{
				"When first_frame_image_url is absent: require aspect_ratio.",
				"When first_frame_image_url is present: forbid aspect_ratio.",
			},
		},
		{
			args: []string{"suno", "text-to-music", "--help"},
			want: []string{
				"When vocal_mode=auto_lyrics: require prompt; forbid lyrics, style, title, negative_tags, vocal_gender, duration_seconds.",
				"When model=suno-v5: forbid duration_seconds.",
			},
		},
		{
			args: []string{"kling", "text-to-video", "--help"},
			want: []string{
				"When model=kling-v3-turbo-text-to-video: forbid enable_sound",
			},
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.args[:2], "/"), func(t *testing.T) {
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}

			cmd := c.command()
			cmd.SetArgs(tc.args)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			output := c.stdout.(*bytes.Buffer).String()
			for _, want := range tc.want {
				if !strings.Contains(output, want) {
					t.Fatalf("expected help to contain %q, got:\n%s", want, output)
				}
			}
		})
	}
}

func TestGeneratedContractHelpIncludesUnconditionalRules(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"minimax-h3", "image-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "Always: require one of first_frame_image_url, last_frame_image_url.") {
		t.Fatalf("expected unconditional rule help, got:\n%s", output)
	}
	if strings.Contains(output, "When :") {
		t.Fatalf("expected no empty conditional clause, got:\n%s", output)
	}
}

func TestGeminiTTSHelpIncludesNestedItemFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"gemini-tts", "text-to-speech", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	for _, want := range []string{
		"speakers[].speaker_id",
		"speakers[].voice_name",
		"Accepted values: Achernar, Achird",
		"dialogue_turns[].speaker_id",
		"dialogue_turns[].text",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected Gemini TTS help to contain %q, got:\n%s", want, output)
		}
	}
}

func TestKlingImageToVideoHelpUsesCanonicalImageFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"kling", "image-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "first_frame_image_url") || !helpHasField(output, "last_frame_image_url") {
		t.Fatalf("expected Kling image-to-video help to include canonical image fields, got:\n%s", output)
	}
	if helpHasField(output, "image_url") || helpHasField(output, "tail_image_url") {
		t.Fatalf("expected Kling image-to-video help not to expose roleless image fields, got:\n%s", output)
	}
}

func TestKlingAvatarHelpUsesCanonicalSourceFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"kling", "avatar", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_image_url") || !helpHasField(output, "source_audio_url") {
		t.Fatalf("expected Kling avatar help to include canonical source fields, got:\n%s", output)
	}
	if helpHasField(output, "image_url") || helpHasField(output, "audio_url") {
		t.Fatalf("expected Kling avatar help not to expose roleless media fields, got:\n%s", output)
	}
}

func TestKlingMotionControlHelpUsesCanonicalMediaFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"kling", "motion-control", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_image_url") || !helpHasField(output, "reference_video_url") {
		t.Fatalf("expected Kling motion-control help to include canonical media fields, got:\n%s", output)
	}
	if helpHasField(output, "input_urls") || helpHasField(output, "video_urls") {
		t.Fatalf("expected Kling motion-control help not to expose provider media fields, got:\n%s", output)
	}
	if !strings.Contains(output, "required for kling-v2.6") {
		t.Fatalf("expected Kling motion-control help to explain v2.6 required fields, got:\n%s", output)
	}
}

func TestVeo31TextToVideoHelpIncludesDurationSeconds(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"veo-3-1", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "duration_seconds") || !strings.Contains(output, "Accepted values: 4, 6, 8.") {
		t.Fatalf("expected Veo 3.1 text-to-video help to include generated duration_seconds values, got:\n%s", output)
	}
}

func TestVeo31UpscaleVideoHelpIncludesTargetResolution(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"veo-3-1", "upscale-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "output_resolution") || !strings.Contains(output, "Accepted values: 1080p, 4k.") {
		t.Fatalf("expected Veo 3.1 upscale-video help to include generated output_resolution values, got:\n%s", output)
	}
}

func TestRunwayAlephHelpUsesCanonicalSourceVideoURL(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"runway-aleph", "edit-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_video_url") {
		t.Fatalf("expected Runway Aleph edit-video help to include source_video_url, got:\n%s", output)
	}
	if helpHasField(output, "video_url") {
		t.Fatalf("expected Runway Aleph edit-video help not to include stale video_url, got:\n%s", output)
	}
}

func TestRunwayTextToVideoHelpUsesCanonicalFirstFrameImageURL(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"runway", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "first_frame_image_url") {
		t.Fatalf("expected Runway text-to-video help to include first_frame_image_url, got:\n%s", output)
	}
	if helpHasField(output, "image_url") {
		t.Fatalf("expected Runway text-to-video help not to include stale image_url, got:\n%s", output)
	}
}

func TestLumaHelpUsesCanonicalSourceVideoURL(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"luma", "modify-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_video_url") {
		t.Fatalf("expected Luma modify-video help to include source_video_url, got:\n%s", output)
	}
	if helpHasField(output, "video_url") {
		t.Fatalf("expected Luma modify-video help not to include stale video_url, got:\n%s", output)
	}
}

func TestTopazHelpUsesCanonicalSourceURLs(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"topaz", "upscale-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_image_url") {
		t.Fatalf("expected Topaz upscale-image help to include source_image_url, got:\n%s", output)
	}
	if helpHasField(output, "image_url") {
		t.Fatalf("expected Topaz upscale-image help to omit image_url, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"topaz", "upscale-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_video_url") {
		t.Fatalf("expected Topaz upscale-video help to include source_video_url, got:\n%s", output)
	}
	if helpHasField(output, "video_url") {
		t.Fatalf("expected Topaz upscale-video help to omit video_url, got:\n%s", output)
	}
}

func TestFluxKontextHelpUsesCanonicalSourceImageURL(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"flux-kontext", "text-to-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "source_image_url") {
		t.Fatalf("expected Flux Kontext text-to-image help to include source_image_url, got:\n%s", output)
	}
	if helpHasField(output, "input_image") {
		t.Fatalf("expected Flux Kontext text-to-image help to omit input_image, got:\n%s", output)
	}
}

func TestIdeogramV3HelpUsesCanonicalSourceImageURL(t *testing.T) {
	for _, action := range []string{"edit-image", "remix-image", "reframe-image"} {
		c := newCLI()
		c.stdout = &bytes.Buffer{}
		c.stderr = &bytes.Buffer{}

		cmd := c.command()
		cmd.SetArgs([]string{"ideogram-v3", action, "--help"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		output := c.stdout.(*bytes.Buffer).String()
		if !helpHasField(output, "source_image_url") {
			t.Fatalf("expected Ideogram V3 %s help to include source_image_url, got:\n%s", action, output)
		}
		if helpHasField(output, "image_url") {
			t.Fatalf("expected Ideogram V3 %s help not to expose image_url, got:\n%s", action, output)
		}
	}
}

func TestIdeogramV3CharacterRemixHelpUsesStyleReferenceImages(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"ideogram-v3", "remix-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "style_reference_image_urls") {
		t.Fatalf("expected Ideogram V3 remix-image help to include style_reference_image_urls, got:\n%s", output)
	}
	if helpHasField(output, "image_urls") {
		t.Fatalf("expected Ideogram V3 remix-image help not to expose image_urls, got:\n%s", output)
	}
}

func TestNanoBananaHelpUsesCanonicalFields(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"nano-banana", "text-to-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	for _, expected := range []string{"aspect_ratio", "output_resolution", "reference_image_urls", "Accepted values by model"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected Nano Banana text-to-image help to include %q, got:\n%s", expected, output)
		}
	}
	for _, stale := range []string{"image_size", "image_input"} {
		if helpHasField(output, stale) {
			t.Fatalf("expected Nano Banana text-to-image help not to include stale field %q, got:\n%s", stale, output)
		}
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"nano-banana", "edit-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	for _, expected := range []string{"aspect_ratio", "source_image_urls"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected Nano Banana edit-image help to include %q, got:\n%s", expected, output)
		}
	}
	for _, stale := range []string{"image_size", "image_urls"} {
		if helpHasField(output, stale) {
			t.Fatalf("expected Nano Banana edit-image help not to include stale field %q, got:\n%s", stale, output)
		}
	}
}

func TestQwen2HelpUsesCanonicalAspectRatio(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"qwen-2", "text-to-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "aspect_ratio") {
		t.Fatalf("expected Qwen 2 text-to-image help to include aspect_ratio, got:\n%s", output)
	}
	if helpHasField(output, "image_size") {
		t.Fatalf("expected Qwen 2 text-to-image help not to include stale image_size, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"qwen-2", "edit-image", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !helpHasField(output, "aspect_ratio") {
		t.Fatalf("expected Qwen 2 edit-image help to include aspect_ratio, got:\n%s", output)
	}
	if !helpHasField(output, "source_image_url") {
		t.Fatalf("expected Qwen 2 edit-image help to include source_image_url, got:\n%s", output)
	}
	if helpHasField(output, "image_size") {
		t.Fatalf("expected Qwen 2 edit-image help not to include stale image_size, got:\n%s", output)
	}
	if helpHasField(output, "image_url") {
		t.Fatalf("expected Qwen 2 edit-image help not to expose old image field name, got:\n%s", output)
	}

}

func TestGeneratedContractValuesAreComposedIntoHelp(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"veo-3-1", "text-to-video", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "Accepted values: 4, 6, 8.") {
		t.Fatalf("expected Veo 3.1 help to append generated duration_seconds values, got:\n%s", output)
	}
	if !helpHasField(output, "input_mode") {
		t.Fatalf("expected Veo 3.1 help to include input_mode, got:\n%s", output)
	}
	for _, expected := range []string{"first_frame_image_url", "last_frame_image_url", "reference_image_urls"} {
		if !helpHasField(output, expected) {
			t.Fatalf("expected Veo 3.1 help to include %s, got:\n%s", expected, output)
		}
	}
	for _, stale := range []string{"generation_type", "image_urls"} {
		if helpHasField(output, stale) {
			t.Fatalf("expected Veo 3.1 help not to include stale field %s, got:\n%s", stale, output)
		}
	}
	if !strings.Contains(output, "Accepted values: text, first_and_last_frames, reference.") {
		t.Fatalf("expected Veo 3.1 help to append generated input_mode values, got:\n%s", output)
	}
	if !strings.Contains(output, "Accepted values: veo-3.1, veo-3.1-fast, veo-3.1-lite.") {
		t.Fatalf("expected Veo 3.1 help to append generated model values, got:\n%s", output)
	}
}

func TestDynamicImageActionsAreRegistered(t *testing.T) {
	cases := []struct {
		service string
		action  string
		fields  []string
	}{
		{service: "flux-2", action: "remix-image", fields: []string{"source_image_urls", "Accepted values: flux-2-flex-remix-image, flux-2-max-remix-image, flux-2-pro-remix-image."}},
		{service: "flux-2", action: "text-to-image", fields: []string{"output_count", "Accepted values: flux-2-flex-text-to-image, flux-2-max-text-to-image, flux-2-pro-text-to-image."}},
		{service: "flux", action: "text-to-image", fields: []string{"output_count", "Accepted values: flux-2-klein, flux-dev, flux-pro."}},
		{service: "flux", action: "remix-image", fields: []string{"source_image_url", "Accepted values: flux-dev, flux-pro."}},
		{service: "imagen-4", action: "remix-image", fields: []string{"source_image_urls", "Accepted values: imagen-4-pro-remix-image."}},
		{service: "seedream", action: "edit-image", fields: []string{"source_image_urls", "Accepted values: seedream-4.5-edit, seedream-5-lite-edit, seedream-5-pro-edit, seedream-v4-edit."}},
		{service: "seedream", action: "decompose-layers", fields: []string{"image_url", "Accepted values: seedream-5-pro-layer-decomposition."}},
	}

	for _, tc := range cases {
		t.Run(tc.service+"/"+tc.action, func(t *testing.T) {
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}

			cmd := c.command()
			cmd.SetArgs([]string{tc.service, tc.action, "--help"})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			output := c.stdout.(*bytes.Buffer).String()
			for _, field := range tc.fields {
				if !strings.Contains(output, field) {
					t.Fatalf("expected %s %s help to include %q, got:\n%s", tc.service, tc.action, field, output)
				}
			}
		})
	}
}

func TestSunoVoiceValidationPhraseCommandsAreRegistered(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"suno", "voice-to-validation-phrase", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "voice_url") || !strings.Contains(output, "vocal_start_seconds") || !strings.Contains(output, "vocal_end_seconds") {
		t.Fatalf("expected Suno voice-to-validation-phrase help to include voice fields, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"suno", "regenerate-validation-phrase", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "task_id") || !strings.Contains(output, "callback_url") {
		t.Fatalf("expected Suno regenerate-validation-phrase help to include regeneration fields, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"suno", "generate-voice", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "verify_url") || !strings.Contains(output, "voice_name") || !strings.Contains(output, "singer_skill_level") {
		t.Fatalf("expected Suno generate-voice help to include custom voice fields, got:\n%s", output)
	}

	c = newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd = c.command()
	cmd.SetArgs([]string{"suno", "check-voice", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output = c.stdout.(*bytes.Buffer).String()
	if !strings.Contains(output, "task_id") {
		t.Fatalf("expected Suno check-voice help to include task_id, got:\n%s", output)
	}
}

func TestSunoCreateMashupHelpDoesNotDuplicateModelField(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"suno", "create-mashup", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if count := helpFieldCount(output, "model"); count != 1 {
		t.Fatalf("expected Suno create-mashup help to list model once, got %d:\n%s", count, output)
	}
}
