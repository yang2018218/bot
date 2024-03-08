package v1

import "time"

type ChatGPTCompletionsUsage struct {
	PromptTokens    int64 `json:"prompt_tokens"`
	CompletionToken int64 `json:"completion_token"`
	TotalTokens     int64 `json:"total_tokens"`
}

type ChatGPTChatCompletionsMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"-"`
}

type ChatGPTChatCompletionsOpts struct {
	Model     string                          `json:"model"`
	Messages  []ChatGPTChatCompletionsMessage `json:"messages"`
	CreatedAt *time.Time                      `json:"-"`
}

type ChatGPTChatCompletionsChoice struct {
	Index        *int                          `json:"index"`
	Message      ChatGPTChatCompletionsMessage `json:"message"`
	Logprobs     interface{}                   `json:"logprobs"`
	FinishReason string                        `json:"finish_reason"`
}

type ChatGPTChatCompletionsResponse struct {
	ID      string                         `json:"id"`
	Object  string                         `json:"object"`
	Created int64                          `json:"created"`
	Model   string                         `json:"model"`
	Choices []ChatGPTChatCompletionsChoice `json:"choices"`
	Usage   ChatGPTCompletionsUsage        `json:"usage"`
}
