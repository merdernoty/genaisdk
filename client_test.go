package genai

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	apiKey := "test-api-key"
	client := New(apiKey)

	if client == nil {
		t.Fatal("ожидался не nil клиент")
	}

	if client.config.APIKey != apiKey {
		t.Errorf("ожидался API ключ %s, получен %s", apiKey, client.config.APIKey)
	}

	if client.config.Model != "gemini-1.5-flash" {
		t.Errorf("ожидалась модель по умолчанию gemini-1.5-flash, получена %s", client.config.Model)
	}
}

func TestNewWithConfig(t *testing.T) {
	config := &ClientConfig{
		APIKey:    "test-key",
		Model:     "gemini-1.5-pro",
		UseVertex: true,
		Project:   "test-project",
		Location:  "us-east1",
	}

	client := NewWithConfig(config)

	if client == nil {
		t.Fatal("ожидался не nil клиент")
	}

	if client.config.APIKey != config.APIKey {
		t.Errorf("неверный API ключ")
	}

	if client.config.Model != config.Model {
		t.Errorf("неверная модель")
	}

	if client.config.UseVertex != config.UseVertex {
		t.Errorf("неверная настройка Vertex")
	}
}

func TestConnectValidation(t *testing.T) {
	ctx := context.Background()

	// Тест без API ключа
	client := New("")
	err := client.Connect(ctx)
	if err != ErrAPIKeyRequired {
		t.Errorf("ожидалась ошибка ErrAPIKeyRequired, получена %v", err)
	}

	// Тест Vertex без проекта
	config := &ClientConfig{
		UseVertex: true,
		Project:   "",
	}
	client = NewWithConfig(config)
	err = client.Connect(ctx)
	if err != ErrProjectRequired {
		t.Errorf("ожидалась ошибка ErrProjectRequired, получена %v", err)
	}
}

func TestGenerateValidation(t *testing.T) {
	ctx := context.Background()
	client := New("test-key")

	// Тест без инициализации
	_, err := client.Generate(ctx, "test")
	if err != ErrClientNotInitialized {
		t.Errorf("ожидалась ошибка ErrClientNotInitialized, получена %v", err)
	}

	// Симулируем инициализированный клиент для тестирования валидации
	// (в реальных тестах нужно мокать API)
	client.client = nil // Оставляем nil для проверки валидации

	_, err = client.Generate(ctx, "")
	if err != ErrEmptyPrompt {
		t.Errorf("ожидалась ошибка ErrEmptyPrompt, получена %v", err)
	}

	_, err = client.Generate(ctx, "   ")
	if err != ErrEmptyPrompt {
		t.Errorf("ожидалась ошибка ErrEmptyPrompt для пробелов, получена %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	apiKey := "test-key"
	config := DefaultConfig(apiKey)

	if config.APIKey != apiKey {
		t.Errorf("неверный API ключ в конфигурации по умолчанию")
	}

	if config.Model != "gemini-1.5-flash" {
		t.Errorf("неверная модель по умолчанию")
	}

	if config.Location != "us-central1" {
		t.Errorf("неверная локация по умолчанию")
	}

	if config.UseVertex != false {
		t.Errorf("UseVertex должен быть false по умолчанию")
	}
}

func TestBuildConfig(t *testing.T) {
	client := New("test-key")

	temperature := float32(0.7)
	maxTokens := int32(1000)
	topP := float32(0.9)
	topK := int32(40)
	systemInstruction := "Ты помощник"

	options := &GenerateOptions{
		Temperature:       &temperature,
		MaxOutputTokens:   &maxTokens,
		TopP:              &topP,
		TopK:              &topK,
		SystemInstruction: systemInstruction,
	}

	config := client.buildConfig(options)

	if config.Temperature == nil || *config.Temperature != temperature {
		t.Errorf("неверная температура в конфигурации")
	}

	if config.MaxOutputTokens != maxTokens {
		t.Errorf("неверное количество токенов в конфигурации")
	}

	if config.TopP == nil || *config.TopP != topP {
		t.Errorf("неверный TopP в конфигурации")
	}

	if config.TopK == nil || int32(*config.TopK) != topK {
		t.Errorf("неверный TopK в конфигурации")
	}

	if config.SystemInstruction == nil {
		t.Error("системная инструкция не должна быть nil")
	}
}
