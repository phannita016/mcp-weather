package dtos

import "github.com/openai/openai-go"

type (
	Message struct {
		Role       string                                 `json:"role"`
		Content    string                                 `json:"content"`
		ToolCallID string                                 `json:"tool_call_id,omitempty"`
		ToolCalls  []openai.ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
	}
)
