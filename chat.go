package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Chat представляет чат сессию
type Chat struct {
	client  *genai.Client
	model   string
	history []*genai.Content
}

// SendMessage отправляет сообщение в чат
func (c *Chat) SendMessage(ctx context.Context, message string) (*Response, error) {
	if c.client == nil {
		return nil, ErrClientNotInitialized
	}

	if strings.TrimSpace(message) == "" {
		return nil, ErrEmptyPrompt
	}

	// Добавляем сообщение пользователя в историю
	userContent := &genai.Content{
		Parts: []*genai.Part{{Text: message}},
	}
	c.history = append(c.history, userContent)

	// Генерируем ответ с учетом всей истории
	response, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		c.history,
		nil,
	)
	if err != nil {
		return nil, NewGenerationError(c.model, message, err)
	}

	// Добавляем ответ модели в историю
	assistantContent := &genai.Content{
		Parts: []*genai.Part{{Text: response.Text()}},
	}
	c.history = append(c.history, assistantContent)

	return &Response{
		Text:     response.Text(),
		Original: response,
	}, nil
}

// SendMessageWithJSON отправляет сообщение с дополнительными JSON данными
func (c *Chat) SendMessageWithJSON(ctx context.Context, message string, data interface{}) (*Response, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать JSON: %w", err)
	}

	fullMessage := fmt.Sprintf("%s\n\nДополнительные данные в формате JSON:\n%s", message, string(jsonBytes))
	return c.SendMessage(ctx, fullMessage)
}

// SendMessageWithOptions отправляет сообщение с дополнительными опциями
func (c *Chat) SendMessageWithOptions(ctx context.Context, message string, options *GenerateOptions) (*Response, error) {
	if c.client == nil {
		return nil, ErrClientNotInitialized
	}

	if strings.TrimSpace(message) == "" {
		return nil, ErrEmptyPrompt
	}

	// Добавляем сообщение пользователя в историю
	userContent := &genai.Content{
		Parts: []*genai.Part{{Text: message}},
	}
	c.history = append(c.history, userContent)

	// Создаем конфигурацию
	var config *genai.GenerateContentConfig
	if options != nil {
		config = c.buildConfig(options)
	}

	// Генерируем ответ с учетом всей истории и опций
	response, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		c.history,
		config,
	)
	if err != nil {
		return nil, NewGenerationError(c.model, message, err)
	}

	// Добавляем ответ модели в историю
	assistantContent := &genai.Content{
		Parts: []*genai.Part{{Text: response.Text()}},
	}
	c.history = append(c.history, assistantContent)

	return &Response{
		Text:     response.Text(),
		Original: response,
	}, nil
}

// ClearHistory очищает историю чата
func (c *Chat) ClearHistory() {
	c.history = make([]*genai.Content, 0)
}

// GetHistory возвращает историю чата
func (c *Chat) GetHistory() []*genai.Content {
	// Возвращаем копию, чтобы избежать случайного изменения
	history := make([]*genai.Content, len(c.history))
	copy(history, c.history)
	return history
}

// buildConfig создает конфигурацию для API из опций
func (c *Chat) buildConfig(options *GenerateOptions) *genai.GenerateContentConfig {
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
