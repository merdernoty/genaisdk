package genai

import (
	"errors"
	"fmt"
)

// Предопределенные ошибки
var (
	ErrClientNotInitialized = errors.New("клиент не инициализирован")
	ErrAPIKeyRequired       = errors.New("API ключ обязателен")
	ErrProjectRequired      = errors.New("проект обязателен для Vertex AI")
	ErrEmptyPrompt          = errors.New("промпт не может быть пустым")
	ErrInvalidModel         = errors.New("указана неверная модель")
)

// ClientError ошибка клиента
type ClientError struct {
	Operation string
	Err       error
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("ошибка в операции %s: %v", e.Operation, e.Err)
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError создает новую ошибку клиента
func NewClientError(operation string, err error) *ClientError {
	return &ClientError{
		Operation: operation,
		Err:       err,
	}
}

// GenerationError ошибка генерации контента
type GenerationError struct {
	Model  string
	Prompt string
	Err    error
}

func (e *GenerationError) Error() string {
	return fmt.Sprintf("ошибка генерации для модели %s: %v", e.Model, e.Err)
}

func (e *GenerationError) Unwrap() error {
	return e.Err
}

// NewGenerationError создает новую ошибку генерации
func NewGenerationError(model, prompt string, err error) *GenerationError {
	return &GenerationError{
		Model:  model,
		Prompt: prompt,
		Err:    err,
	}
}
