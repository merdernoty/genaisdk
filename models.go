package genai

import "google.golang.org/genai"

// ClientConfig конфигурация клиента
type ClientConfig struct {
	APIKey    string
	Project   string
	Location  string
	Model     string
	UseVertex bool
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig(apiKey string) *ClientConfig {
	return &ClientConfig{
		APIKey:    apiKey,
		Model:     "gemini-1.5-flash",
		Location:  "us-central1",
		UseVertex: false,
	}
}

// GenerateOptions опции для генерации контента
type GenerateOptions struct {
	Temperature       *float32
	MaxOutputTokens   *int32
	TopP              *float32
	TopK              *int32
	SystemInstruction string
}

// Response структура ответа от модели
type Response struct {
	Text     string
	Original *genai.GenerateContentResponse
}

// StreamHandler интерфейс для обработки потоковых ответов
type StreamHandler interface {
	OnChunk(text string)
	OnComplete(fullText string)
	OnError(err error)
}

// DefaultStreamHandler базовая реализация StreamHandler
type DefaultStreamHandler struct {
	ChunkFunc    func(string)
	CompleteFunc func(string)
	ErrorFunc    func(error)
}

func (h *DefaultStreamHandler) OnChunk(text string) {
	if h.ChunkFunc != nil {
		h.ChunkFunc(text)
	}
}

func (h *DefaultStreamHandler) OnComplete(fullText string) {
	if h.CompleteFunc != nil {
		h.CompleteFunc(fullText)
	}
}

func (h *DefaultStreamHandler) OnError(err error) {
	if h.ErrorFunc != nil {
		h.ErrorFunc(err)
	}
}
