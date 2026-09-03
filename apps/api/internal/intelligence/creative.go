package intelligence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"orbit/apps/api/config"
)

type CreativeGenerator struct {
	cfg *config.Config
}

func NewCreativeGenerator(cfg *config.Config) *CreativeGenerator {
	return &CreativeGenerator{cfg: cfg}
}

type GenerationRequest struct {
	Type        string            `json:"type"`
	ProductName string            `json:"product_name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Platform    string            `json:"platform"`
	Language    string            `json:"language"`
	Variations  int               `json:"variations"`
	Style       string            `json:"style"`
}

func (g *CreativeGenerator) GenerateText(ctx context.Context, req GenerationRequest) ([]GeneratedCreative, error) {
	if g.cfg.AI.APIKey == "" {
		return g.generateFallbackText(req), nil
	}

	creatives, err := g.callLLM(ctx, req)
	if err != nil {
		return g.generateFallbackText(req), nil
	}

	return creatives, nil
}

func (g *CreativeGenerator) GenerateImage(ctx context.Context, prompt string) (string, error) {
	if g.cfg.AI.APIKey == "" {
		return "", fmt.Errorf("LLM_API_KEY not configured")
	}

	imageURL, err := g.callImageGeneration(ctx, prompt)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

func (g *CreativeGenerator) GenerateVideo(ctx context.Context, imagePrompt string, duration int) (string, error) {
	if g.cfg.AI.APIKey == "" {
		return "", fmt.Errorf("LLM_API_KEY not configured")
	}

	videoURL, err := g.callVideoGeneration(ctx, imagePrompt, duration)
	if err != nil {
		return "", err
	}

	return videoURL, nil
}

func (g *CreativeGenerator) callLLM(ctx context.Context, req GenerationRequest) ([]GeneratedCreative, error) {
	apiKey := g.cfg.AI.APIKey
	model := g.cfg.AI.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	prompt := fmt.Sprintf(`You are an expert ad copywriter for %s. Generate %d ad creative variations for a product.

Product: %s
Description: %s
Category: %s
Platform: %s
Style: %s

For each variation, provide:
1. headline: Short, attention-grabbing headline (max 40 chars for Meta, 60 for Google)
2. description: Compelling description (max 125 chars for Meta, 90 for Google)
3. cta: Call-to-action button text (e.g., "Shop Now", "Learn More", "Buy Now")
4. prompt: Image generation prompt for this creative

Return JSON array: [{"headline": "...", "description": "...", "cta": "...", "prompt": "..."}]`,
		req.Platform, req.Variations, req.ProductName, req.Description, req.Category, req.Platform, req.Style)

	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 500,
		"temperature": 0.8,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	baseURL := g.cfg.AI.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: g.cfg.AI.Timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from LLM")
	}

	content := result.Choices[0].Message.Content

	var parsed []map[string]string
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	creatives := make([]GeneratedCreative, 0, len(parsed))
	for i, p := range parsed {
		creatives = append(creatives, GeneratedCreative{
			ID:          fmt.Sprintf("creative_%d", i+1),
			Type:        "text",
			Headline:    p["headline"],
			Description: p["description"],
			CTA:         p["cta"],
			Prompt:      p["prompt"],
			Platform:    req.Platform,
			Score:       0.7 + float64(i)*0.05,
		})
	}

	return creatives, nil
}

func (g *CreativeGenerator) callImageGeneration(ctx context.Context, prompt string) (string, error) {
	apiKey := g.cfg.AI.APIKey

	requestBody := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": prompt,
		"n":      1,
		"size":   "1024x1024",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	baseURL := g.cfg.AI.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/images/generations", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image generation API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	return result.Data[0].URL, nil
}

func (g *CreativeGenerator) callVideoGeneration(ctx context.Context, imagePrompt string, duration int) (string, error) {
	return "", fmt.Errorf("video generation not yet implemented")
}

func (g *CreativeGenerator) generateFallbackText(req GenerationRequest) []GeneratedCreative {
	variations := []struct {
		headline    string
		description string
		cta         string
	}{
		{fmt.Sprintf("Shop %s", req.ProductName), fmt.Sprintf("Discover %s - %s", req.ProductName, req.Description), "Shop Now"},
		{fmt.Sprintf("%s - Limited Time", req.ProductName), fmt.Sprintf("Get %s today. Fast shipping!", req.ProductName), "Buy Now"},
		{fmt.Sprintf("New: %s", req.ProductName), fmt.Sprintf("Upgrade your %s collection with %s", req.Category, req.ProductName), "Learn More"},
	}

	creatives := make([]GeneratedCreative, 0, len(variations))
	for i, v := range variations {
		creatives = append(creatives, GeneratedCreative{
			ID:          fmt.Sprintf("creative_%d", i+1),
			Type:        "text",
			Headline:    v.headline,
			Description: v.description,
			CTA:         v.cta,
			Prompt:      fmt.Sprintf("Professional product photo of %s, clean background, e-commerce style", req.ProductName),
			Platform:    req.Platform,
			Score:       0.7 + float64(i)*0.05,
		})
	}

	return creatives
}
