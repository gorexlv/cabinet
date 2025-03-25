package kimi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client struct {
	apiKey     string
	apiBaseURL string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("KIMI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("KIMI_API_KEY environment variable is not set")
	}

	return &Client{
		apiKey:     apiKey,
		apiBaseURL: "https://api.moonshot.cn/v1/chat/completions",
	}, nil
}

func (c *Client) GenerateSummary(content string) (string, error) {
	prompt := fmt.Sprintf("请为以下文章生成一个简洁的摘要（不超过200字）：\n\n%s", content)
	return c.callAPI(prompt)
}

func (c *Client) GenerateTags(content string) (string, error) {
	prompt := fmt.Sprintf("请为以下文章提取5-8个关键词标签，用逗号分隔：\n\n%s", content)
	return c.callAPI(prompt)
}

func (c *Client) callAPI(prompt string) (string, error) {
	reqBody := Request{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.apiBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return apiResp.Choices[0].Message.Content, nil
}
