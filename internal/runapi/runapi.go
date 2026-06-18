// Package runapi provides the official Go SDK for [RunAPI.ai].
//
// Use [NewClient] to create an aggregate client with access to all services,
// or import individual service packages (suno, veo31, nanobanana, runway, runwayaleph) directly.
//
//	client, err := runapi.NewClient(option.WithAPIKey("sk-your-api-key"))
//	result, err := client.Suno.TextToMusic.Run(ctx, suno.TextToMusicParams{...})
//
// [RunAPI.ai]: https://runapi.ai
package runapi

import (
	"github.com/runapi-ai/core-sdk/go/base"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/elevenlabs-sdk/go/elevenlabs"
	"github.com/runapi-ai/flux-2-sdk/go/flux2"
	"github.com/runapi-ai/flux-kontext-sdk/go/fluxkontext"
	"github.com/runapi-ai/gemini-omni-sdk/go/geminiomni"
	"github.com/runapi-ai/gpt-4o-image-sdk/go/gpt4oimage"
	"github.com/runapi-ai/gpt-image-sdk/go/gptimage"
	"github.com/runapi-ai/gpt-image-2-sdk/go/gptimage2"
	"github.com/runapi-ai/grok-imagine-sdk/go/grokimagine"
	"github.com/runapi-ai/hailuo-sdk/go/hailuo"
	"github.com/runapi-ai/happyhorse-sdk/go/happyhorse"
	"github.com/runapi-ai/ideogram-v3-sdk/go/ideogramv3"
	"github.com/runapi-ai/imagen-4-sdk/go/imagen4"
	"github.com/runapi-ai/infinitetalk-sdk/go/infinitetalk"
	"github.com/runapi-ai/kling-sdk/go/kling"
	"github.com/runapi-ai/luma-sdk/go/luma"
	"github.com/runapi-ai/nano-banana-sdk/go/nanobanana"
	"github.com/runapi-ai/core-sdk/go/option"
	"github.com/runapi-ai/qwen-2-sdk/go/qwen2"
	"github.com/runapi-ai/recraft-sdk/go/recraft"
	"github.com/runapi-ai/runway-sdk/go/runway"
	"github.com/runapi-ai/runway-aleph-sdk/go/runwayaleph"
	"github.com/runapi-ai/seedance-sdk/go/seedance"
	"github.com/runapi-ai/seedream-sdk/go/seedream"
	"github.com/runapi-ai/suno-sdk/go/suno"
	"github.com/runapi-ai/topaz-sdk/go/topaz"
	"github.com/runapi-ai/veo-3.1-sdk/go/veo31"
	"github.com/runapi-ai/wan-sdk/go/wan"
	"github.com/runapi-ai/z-image-sdk/go/zimage"
)

// Client is the aggregate RunAPI client combining all services.
// All sub-clients share a single HTTP transport, so [option.ClientOption] values
// (API key, base URL, timeouts) only need to be set once via [NewClient].
type Client struct {
	// Base provides the Universal Resources (Files, Account) on every client.
	base.Base
	// Suno creates and manipulates music: text-to-music, extend, cover, remix,
	// stem separation, lyrics, MIDI, sound effects, persona voices, and more.
	Suno *suno.Client
	// Veo31 generates, extends, and upscales video with Veo 3.1 models.
	Veo31 *veo31.Client
	// NanoBanana generates and edits images with Nano Banana models.
	NanoBanana *nanobanana.Client
	// Imagen4 generates images and remixes existing ones with Imagen 4 models.
	Imagen4 *imagen4.Client
	// Seedance generates video from text or image prompts with Seedance models.
	Seedance *seedance.Client
	// Seedream generates images and applies prompt-guided edits with Seedream models.
	Seedream *seedream.Client
	// Runway generates and extends video with Runway Gen-4 models.
	Runway *runway.Client
	// RunwayAleph edits existing video with Runway Aleph models.
	RunwayAleph *runwayaleph.Client
	// Kling generates video from text or images and supports AI avatars and
	// motion control with Kling models.
	Kling *kling.Client
	// FluxKontext generates images with Flux Kontext models.
	FluxKontext *fluxkontext.Client
	// Flux2 generates images and remixes existing ones with Flux 2 models.
	Flux2 *flux2.Client
	// GeminiOmni creates audio, character voices, and video with Gemini models.
	GeminiOmni *geminiomni.Client
	// Qwen2 generates, remixes, and edits images with Qwen 2 models.
	Qwen2 *qwen2.Client
	// Recraft upscales images and removes backgrounds using Recraft models.
	Recraft *recraft.Client
	// ZImage generates images from text with Z-Image models.
	ZImage *zimage.Client
	// IdeogramV3 generates, edits, remixes, and reframes images with Ideogram V3 models.
	IdeogramV3 *ideogramv3.Client
	// Elevenlabs synthesizes speech, dialogue, and sound effects, and performs
	// speech-to-text transcription and audio isolation with ElevenLabs models.
	Elevenlabs *elevenlabs.Client
	// InfiniteTalk generates lip-synced video from an audio track and a face image.
	InfiniteTalk *infinitetalk.Client
	// Wan generates video, images, and animations from text, images, or speech
	// with Wan models.
	Wan *wan.Client
	// Luma modifies existing video using Luma models.
	Luma *luma.Client
	// Hailuo generates video from text or images with Hailuo models.
	Hailuo *hailuo.Client
	// HappyHorse generates video from text, images, or character prompts, and
	// edits existing video with HappyHorse models.
	HappyHorse *happyhorse.Client
	// GptImage generates and edits images with GPT Image 1.5 models.
	GptImage *gptimage.Client
	// GptImage2 generates and edits images with GPT Image 2 models.
	GptImage2 *gptimage2.Client
	// Gpt4oImage generates images with GPT-4o Image models.
	Gpt4oImage *gpt4oimage.Client
	// GrokImagine generates video and images, edits images, extends and
	// upscales outputs with Grok-Imagine models.
	GrokImagine *grokimagine.Client
	// Topaz upscales images and video to higher resolutions using Topaz models.
	Topaz *topaz.Client
}

// NewClient creates an aggregate client with a shared HTTP transport.
// All service sub-clients inherit the resolved options (API key, base URL,
// timeouts, retry policy). If no API key option is provided, the RUNAPI_API_KEY
// environment variable is used.
func NewClient(opts ...option.ClientOption) (*Client, error) {
	resolved, err := option.ResolveClientOptions(opts...)
	if err != nil {
		return nil, err
	}
	if resolved.UserAgent == "" {
		resolved.UserAgent = core.SDKUserAgent(Version)
	}
	httpClient, err := core.NewHTTPClient(resolved)
	if err != nil {
		return nil, err
	}
	return &Client{
		Base:         base.New(httpClient),
		Suno:         suno.NewClientWithHTTP(httpClient),
		Veo31:        veo31.NewClientWithHTTP(httpClient),
		NanoBanana:   nanobanana.NewClientWithHTTP(httpClient),
		Imagen4:      imagen4.NewClientWithHTTP(httpClient),
		Seedance:     seedance.NewClientWithHTTP(httpClient),
		Seedream:     seedream.NewClientWithHTTP(httpClient),
		Runway:       runway.NewClientWithHTTP(httpClient),
		RunwayAleph:  runwayaleph.NewClientWithHTTP(httpClient),
		Kling:        kling.NewClientWithHTTP(httpClient),
		FluxKontext:  fluxkontext.NewClientWithHTTP(httpClient),
		Flux2:        flux2.NewClientWithHTTP(httpClient),
		GeminiOmni:   geminiomni.NewClientWithHTTP(httpClient),
		Qwen2:        qwen2.NewClientWithHTTP(httpClient),
		Recraft:      recraft.NewClientWithHTTP(httpClient),
		ZImage:       zimage.NewClientWithHTTP(httpClient),
		IdeogramV3:   ideogramv3.NewClientWithHTTP(httpClient),
		Elevenlabs:   elevenlabs.NewClientWithHTTP(httpClient),
		InfiniteTalk: infinitetalk.NewClientWithHTTP(httpClient),
		Wan:          wan.NewClientWithHTTP(httpClient),
		Luma:         luma.NewClientWithHTTP(httpClient),
		Hailuo:       hailuo.NewClientWithHTTP(httpClient),
		HappyHorse:   happyhorse.NewClientWithHTTP(httpClient),
		GptImage:     gptimage.NewClientWithHTTP(httpClient),
		GptImage2:    gptimage2.NewClientWithHTTP(httpClient),
		Gpt4oImage:   gpt4oimage.NewClientWithHTTP(httpClient),
		GrokImagine:  grokimagine.NewClientWithHTTP(httpClient),
		Topaz:        topaz.NewClientWithHTTP(httpClient),
	}, nil
}
