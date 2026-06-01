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
	"github.com/runapi-ai/cli/internal/account"
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
type Client struct {
	// Account provides account info and balance operations.
	Account *account.Client
	// Suno provides music generation operations.
	Suno *suno.Client
	// Veo31 provides Veo 3.1 video operations.
	Veo31 *veo31.Client
	// NanoBanana provides image generation operations.
	NanoBanana *nanobanana.Client
	// Imagen4 provides Imagen 4 image generation operations.
	Imagen4 *imagen4.Client
	// Seedance provides video generation operations.
	Seedance *seedance.Client
	// Seedream provides image generation operations.
	Seedream *seedream.Client
	// Runway provides video generation operations.
	Runway *runway.Client
	// RunwayAleph provides video editing operations.
	RunwayAleph *runwayaleph.Client
	// Kling provides video generation operations.
	Kling *kling.Client
	// FluxKontext provides Flux Kontext image generation operations.
	FluxKontext *fluxkontext.Client
	// Flux2 provides Flux 2 image generation operations.
	Flux2 *flux2.Client
	// GeminiOmni provides Gemini Omni audio, character, and video operations.
	GeminiOmni *geminiomni.Client
	// Qwen2 provides Qwen2 image edit operations.
	Qwen2 *qwen2.Client
	// Recraft provides image post-processing operations.
	Recraft *recraft.Client
	// ZImage provides Z-Image text-to-image operations.
	ZImage *zimage.Client
	// IdeogramV3 provides Ideogram V3 image generation operations.
	IdeogramV3 *ideogramv3.Client
	// Elevenlabs provides speech and audio operations.
	Elevenlabs *elevenlabs.Client
	// InfiniteTalk provides lip-sync video operations.
	InfiniteTalk *infinitetalk.Client
	// Wan provides Wan video and image generation operations.
	Wan *wan.Client
	// Luma provides video modification operations.
	Luma *luma.Client
	// Hailuo provides Hailuo video generation operations.
	Hailuo *hailuo.Client
	// HappyHorse provides text, image, character-guided text, and edit-video operations.
	HappyHorse *happyhorse.Client
	// GptImage provides GPT Image 1.5 image generation and editing operations.
	GptImage *gptimage.Client
	// GptImage2 provides GPT Image 2 image generation and editing operations.
	GptImage2 *gptimage2.Client
	// Gpt4oImage provides Gpt4o Image generation operations.
	Gpt4oImage *gpt4oimage.Client
	// GrokImagine provides Grok-Imagine multimodal generation operations.
	GrokImagine *grokimagine.Client
	// Topaz provides image and video upscale operations.
	Topaz *topaz.Client
}

// NewClient creates an aggregate client with a shared HTTP transport.
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
		Account:      account.NewClientWithHTTP(httpClient),
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
