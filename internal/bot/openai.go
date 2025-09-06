package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// openAIResponse represents the response from OpenAI Chat Completions API
type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
}

// openAIChoice represents a single choice in the OpenAI response
type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

// openAIMessage represents a message in the OpenAI response
type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIEngine calls OpenAI's Chat Completions API (simple, cheap, effective).
type OpenAIEngine struct {
	apiKey string
	model  string
	url    string
	client *http.Client
}

func NewOpenAIEngine(apiKey, model string) *OpenAIEngine {
	return &OpenAIEngine{
		apiKey: apiKey,
		model:  model,
		url:    "https://api.openai.com/v1/chat/completions",
		client: &http.Client{Timeout: 22 * time.Second},
	}
}

func (e *OpenAIEngine) Generate(ctx context.Context, topic, stance string, history []HistoryItem, userMessage string) (string, error) {
	if e.apiKey == "" {
		return "", errors.New("OPENAI_API_KEY is missing")
	}

	messages := buildMessages(topic, stance, history, userMessage)

	payload := map[string]any{
		"model":       e.model,
		"messages":    messages,
		"temperature": 0.9,
		"max_tokens":  400,
	}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai http %d: %s", resp.StatusCode, string(body))
	}

	var out openAIResponse

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	if len(out.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return out.Choices[0].Message.Content, nil
}
