package dtos

import "github.com/openai/openai-go"

type (
	Tool struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		InputSchema map[string]any `json:"input_schema"`
	}

	Message struct {
		Role       string                                 `json:"role"`
		Content    interface{}                            `json:"content"`
		ToolCalls  []openai.ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
		ToolCallID string                                 `json:"tool_call_id,omitempty"`
	}
)
