package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"weather/client/dtos"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (a *App) StartPromptLoop() error {
	defer a.session.Close()
	ctx := context.Background()

	if err := a.loadTools(ctx); err != nil {
		return err
	}

	if err := a.promptLoop(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) loadTools(ctx context.Context) error {
	listTools, err := a.session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return err
	}

	for _, t := range listTools.Tools {
		var schemaMap map[string]any
		if t.InputSchema != nil {
			b, err := json.Marshal(t.InputSchema)
			if err == nil {
				_ = json.Unmarshal(b, &schemaMap)
			}
		}
		a.tools = append(a.tools, dtos.Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: schemaMap,
		})
	}

	fmt.Println("✅ Connected to server with tools:")
	for _, t := range a.tools {
		fmt.Println("-", t.Name)
	}
	return nil
}

func (a *App) promptLoop(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n🤖 Query: ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "quit" {
			break
		}

		if query == "" {
			continue
		}

		if err := a.handleQuery(ctx, query); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) handleQuery(ctx context.Context, query string) error {
	messages := []dtos.Message{{Role: "user", Content: query}}

	finalResult := ""
	for {
		// Call OpenAI engine with current conversation context
		completion, err := a.openAIEngine.OpenAiEngine(messages, a.tools)
		if err != nil {
			slog.Error("Error generating tool call", "error", err)
			return err
		}
		completionMsg := completion.Choices[0].Message

		// pretty, err := json.MarshalIndent(completion, "", "  ")
		// if err != nil {
		// 	slog.Error("Failed to marshal JSON", "error", err)
		// } else {
		// 	fmt.Println(string(pretty))
		// }

		// Extract tool call details of type "function" from the assistant's response message
		toolCalls := a.openAIEngine.ExtractToolCalls(completionMsg.ToolCalls)

		finalResult = completionMsg.Content
		if len(toolCalls) == 0 {
			break
		}

		// Add assistant message that includes the tool calls to the message history.
		// This preserves the fact that the assistant requested tool usage.
		messages = append(messages, dtos.Message{
			Role:      "assistant",
			Content:   "",
			ToolCalls: toolCalls,
		})

		// Iterate through each tool call returned by the assistant
		for _, toolCall := range toolCalls {
			// Execute the tool with parsed arguments
			var parsedArgs map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &parsedArgs); err != nil {
				slog.Error("Invalid tool args", "tool", toolCall.Function.Name, "err", err)
				continue
			}

			fmt.Printf("🛠️  Calling tool: %s args=%v\n", toolCall.Function.Name, parsedArgs)

			// Execute the tool with parsed arguments
			result, err := a.session.CallTool(ctx, &mcp.CallToolParams{
				Name:      toolCall.Function.Name,
				Arguments: parsedArgs,
			})
			if err != nil {
				// messages = append(messages, dtos.Message{
				// 	Role:       "tool",
				// 	Content:    fmt.Sprintf("Error: %v", err),
				// 	ToolCallID: toolCall.ID,
				// })
				slog.Error("Error calling tool", "tool", toolCall.Function.Name, "err", err)
				continue
			}

			// Add each piece of content returned by the tool to the conversation
			for _, c := range result.Content {
				text := c.(*mcp.TextContent).Text
				messages = append(messages, dtos.Message{
					Role:       "tool",
					Content:    text,
					ToolCallID: toolCall.ID,
				})
			}
		}
	}

	fmt.Println("🎉 Final combined result:")
	fmt.Println(finalResult)
	return nil
}
