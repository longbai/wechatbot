package gpt

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/longbai/wechatbot/config"
)

const Backend = "https://api.openai.com/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
	Index        int         `json:"index"`
	Message      RoleContent `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type RoleContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string        `json:"model"`
	Prompt           string        `json:"prompt,omitempty"`
	MaxTokens        int           `json:"max_tokens"`
	Temperature      float32       `json:"temperature"`
	TopP             int           `json:"top_p"`
	FrequencyPenalty int           `json:"frequency_penalty"`
	PresencePenalty  float32       `json:"presence_penalty"`
	Stop             []string      `json:"stop"`
	Messages         []RoleContent `json:"messages,omitempty"`
}

func HasKeyWords(msg string) bool {
	keywords := config.LoadConfig().Keys
	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}
	return false
}

func Completions(msg string) (string, error) {
	if HasKeyWords(msg) {
		return "这个话题对我太沉重, 我还小", nil
	}
	cfg := openai.DefaultConfig(config.LoadConfig().ApiKey)
	backend := Backend

	if config.LoadConfig().Backend != "" {
		backend = config.LoadConfig().Backend
	}
	cfg.BaseURL = backend + "v1"
	client := openai.NewClientWithConfig(cfg)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
			Temperature:      0.8,
			TopP:             1,
			Stop:             []string{"股票", "投资", "基金", "理财"},
			PresencePenalty:  0.6,
			FrequencyPenalty: 0,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
