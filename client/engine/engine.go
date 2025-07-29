package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"weather/client/dtos"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// Anthropics config
const (
	ANTHROPIC_API_KEY = ""
	MODEL             = shared.ChatModelGPT4o
)

func OpenAiEngine(messages []dtos.Message, tools []dtos.Tool) ([]openai.ChatCompletionMessageToolCallFunction, string, error) {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithAPIKey(ANTHROPIC_API_KEY),
		option.WithHeader("Content-Type", "application/json"),
	)

	var openaiMessages []openai.ChatCompletionMessageParamUnion
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessageParamUnion{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: param.Opt[string]{Value: msg.Content.(string)},
				},
				Role: constant.User(msg.Role),
			},
		})
	}

	var toolParams []openai.ChatCompletionToolParam
	for _, t := range tools {
		toolParams = append(toolParams, openai.ChatCompletionToolParam{
			Type: "function",
			Function: shared.FunctionDefinitionParam{
				Name:        t.Name,
				Description: param.Opt[string]{Value: t.Description},
				Parameters:  shared.FunctionParameters(t.InputSchema),
			},
		})
	}

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     MODEL,
		Messages:  openaiMessages,
		Tools:     toolParams,
		MaxTokens: param.Opt[int64]{Value: 1000},
	})
	if err != nil {
		return nil, "", err
	}

	pretty, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal JSON", "error", err)
	} else {
		fmt.Println(string(pretty))
	}

	var extractedCalls []openai.ChatCompletionMessageToolCallFunction
	for _, toolCall := range resp.Choices[0].Message.ToolCalls {
		if toolCall.Type == "function" {

			extractedCalls = append(extractedCalls, openai.ChatCompletionMessageToolCallFunction{
				Name:      toolCall.Function.Name,
				Arguments: toolCall.Function.Arguments,
			})
		}
	}

	return extractedCalls, resp.Choices[0].Message.Content, nil
}
