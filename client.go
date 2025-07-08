package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Client основной клиент SDK
type Client struct {
	client *genai.Client
	config *ClientConfig
}

// New создает новый клиент с API ключом
func New(apiKey string) *Client {
	return &Client{
		config: DefaultConfig(apiKey),
	}
}

// NewWithConfig создает клиент с кастомной конфигурацией
func NewWithConfig(config *ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// Connect подключается к API
func (c *Client) Connect(ctx context.Context) error {
	if c.config.APIKey == "" && !c.config.UseVertex {
		return ErrAPIKeyRequired
	}

	if c.config.UseVertex && c.config.Project == "" {
		return ErrProjectRequired
	}

	var clientConfig *genai.ClientConfig

	if c.config.UseVertex {
		clientConfig = &genai.ClientConfig{
			Project:  c.config.Project,
			Location: c.config.Location,
			Backend:  genai.BackendVertexAI,
		}
	} else {
		clientConfig = &genai.ClientConfig{
			APIKey:  c.config.APIKey,
			Backend: genai.BackendGeminiAPI,
		}
	}

	client, err := genai.NewClient(ctx, clientConfig)
	if err != nil {
		return NewClientError("подключение", err)
	}

	c.client = client
	return nil
}

// Generate генерирует текст из промпта
func (c *Client) Generate(ctx context.Context, prompt string) (*Response, error) {
	return c.GenerateWithOptions(ctx, prompt, nil)
}

// GenerateWithOptions генерирует текст с дополнительными опциями
func (c *Client) GenerateWithOptions(ctx context.Context, prompt string, options *GenerateOptions) (*Response, error) {
	if c.client == nil {
		return nil, ErrClientNotInitialized
	}

	if strings.TrimSpace(prompt) == "" {
		return nil, ErrEmptyPrompt
	}

	parts := []*genai.Part{{Text: prompt}}
	content := &genai.Content{Parts: parts}

	var config *genai.GenerateContentConfig
	if options != nil {
		config = c.buildConfig(options)
	}

	response, err := c.client.Models.GenerateContent(
		ctx,
		c.config.Model,
		[]*genai.Content{content},
		config,
	)
	if err != nil {
		return nil, NewGenerationError(c.config.Model, prompt, err)
	}

	return &Response{
		Text:     response.Text(),
		Original: response,
	}, nil
}

// GenerateWithJSON генерирует текст с дополнительными JSON данными
func (c *Client) GenerateWithJSON(ctx context.Context, prompt string, data interface{}) (*Response, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать JSON: %w", err)
	}

	fullPrompt := fmt.Sprintf("%s\n\nДополнительные данные в формате JSON:\n%s", prompt, string(jsonBytes))
	return c.Generate(ctx, fullPrompt)
}

// Chat представляет сессию чата с историей сообщений

// NewChat создает новую чат сессию
func (c *Client) NewChat() *Chat {
	return &Chat{
		client:  c.client,
		model:   c.config.Model,
		history: make([]*genai.Content, 0),
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	c.client = nil
	return nil
}

// buildConfig создает конфигурацию для API из опций
func (c *Client) buildConfig(options *GenerateOptions) *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{}

	if options.Temperature != nil {
		config.Temperature = options.Temperature
	}

	if options.MaxOutputTokens != nil {
		config.MaxOutputTokens = *options.MaxOutputTokens
	}

	if options.TopP != nil {
		config.TopP = options.TopP
	}

	if options.TopK != nil {
		topK := float32(*options.TopK)
		config.TopK = &topK
	}

	if options.SystemInstruction != "" {
		config.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: options.SystemInstruction}},
		}
	}

	return config
}
