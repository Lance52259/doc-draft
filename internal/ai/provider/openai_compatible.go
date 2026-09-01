package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ChatMessage is one chat turn.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResult is a model reply.
type CompletionResult struct {
	Content string
	Model   string
}

// Provider is the AI chat-completion contract.
type Provider interface {
	Complete(ctx context.Context, messages []ChatMessage, temperature float64, responseFormatJSON bool) (*CompletionResult, error)
}

// OpenAICompatible talks to DeepSeek (or any OpenAI-compatible) /chat/completions.
type OpenAICompatible struct {
	APIKey         string
	BaseURL        string
	Model          string
	Timeout        time.Duration
	MaxRetries     int
	HTTPClient     *http.Client
}

func NewDeepSeek(apiKey, baseURL, model string, timeoutSec, maxRetries int) *OpenAICompatible {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	if model == "" {
		model = "deepseek-chat"
	}
	if timeoutSec <= 0 {
		timeoutSec = 120
	}
	return &OpenAICompatible{
		APIKey:     apiKey,
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Model:      model,
		Timeout:    time.Duration(timeoutSec) * time.Second,
		MaxRetries: maxRetries,
		HTTPClient: &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (p *OpenAICompatible) endpoint() string {
	base := p.BaseURL
	if strings.HasSuffix(base, "/v1") {
		return base + "/chat/completions"
	}
	return base + "/chat/completions"
}

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []ChatMessage `json:"messages"`
	Temperature    float64       `json:"temperature"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}

type chatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func (p *OpenAICompatible) Complete(ctx context.Context, messages []ChatMessage, temperature float64, responseFormatJSON bool) (*CompletionResult, error) {
	if p.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	reqBody := chatRequest{
		Model:       p.Model,
		Messages:    messages,
		Temperature: temperature,
	}
	if responseFormatJSON {
		reqBody.ResponseFormat = &struct {
			Type string `json:"type"`
		}{Type: "json_object"}
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= p.MaxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint(), bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := p.HTTPClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("AI HTTP %d: %s", resp.StatusCode, truncate(string(body), 500))
			continue
		}
		var parsed chatResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			lastErr = err
			continue
		}
		if len(parsed.Choices) == 0 {
			lastErr = fmt.Errorf("empty choices from AI")
			continue
		}
		return &CompletionResult{Content: parsed.Choices[0].Message.Content, Model: parsed.Model}, nil
	}
	return nil, fmt.Errorf("AI completion failed after retries: %w", lastErr)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
