package llm

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/typescript-any/llm-playground/internal/config"
)

func NewClient() *openai.Client {
	cfg := config.LoadConfig()
	client := openai.NewClient(option.WithBaseURL(cfg.OpenRouterApiEndpoint), option.WithAPIKey(cfg.OpenRouterApiKey))
	return &client
}
