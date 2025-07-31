package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

// StartPromptLoop initializes the session by loading available tools,
// then starts an interactive prompt loop to handle user queries.
// It ensures the session is properly closed at the end.
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

// loadTools fetches the list of available tools from the session server,
// converts their input schemas into a usable format,
// and registers them for later use in chat completions.
func (a *App) loadTools(ctx context.Context) error {
	listTools, err := a.session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return err
	}

	for _, t := range listTools.Tools {
		schemaMap := map[string]any{
			"type":       t.InputSchema.Type,
			"properties": t.InputSchema.Properties,
			"required":   t.InputSchema.Required,
		}

		a.tools = append(a.tools, openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        t.Name,
				Description: param.Opt[string]{Value: t.Description},
				Parameters:  schemaMap,
			},
			Type: "function",
		})
	}

	fmt.Println("✅ Connected to server with tools:", len(a.tools))
	return nil
}

// promptLoop runs an interactive command-line prompt,
// reading user input continuously until "quit" is entered,
// and passing each non-empty query to handleQuery for processing.
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

/*
		handleQuery manages the conversation between the user, the LLM, and external tools.
		It sends the user query to the LLM, executes any requested tools, updates the conversation,
		and repeats until no more tool calls are needed.
		Finally, it outputs the combined response from the LLM.

	  handleQuery Flow:
	  1. Start with user query as initial message.
	  2. Send current messages to the LLM.
	  3. Check if the LLM response includes any ToolCalls.
	     - If no ToolCalls, finish and output the final result.
	     - If ToolCalls exist:
	        a. Append ToolCalls to the conversation messages.
	        b. Call each tool in order with provided arguments.
	        c. Append tool responses back to the messages.
	  4. Repeat from step 2 with updated messages until no more ToolCalls.
*/
func (a *App) handleQuery(ctx context.Context, query string) error {
	var finalResult string

	// Start conversation with user input
	messages := []openai.ChatCompletionMessage{{Role: "user", Content: query}}

	for {
		// Call OpenAI engine with current conversation context
		completion, err := a.openAIEngine.OpenAiEngine(messages, a.tools)
		if err != nil {
			slog.Error("Error generating tool call", "error", err)
			return err
		}
		// Extract response content and any tool calls from the model
		completionMsg := a.groupChatCompletionChoices(completion.Choices)
		toolCalls := completionMsg.ToolCalls
		finalResult = completionMsg.Content

		// Break the loop if there are no tool calls
		if len(toolCalls) == 0 {
			break
		}

		// Add assistant message with tool calls to the history
		messages = append(messages, openai.ChatCompletionMessage{
			Role:      "assistant",
			Content:   "",
			ToolCalls: toolCalls,
		})

		// Handle each tool call from the assistant
		for _, toolCall := range toolCalls {
			fmt.Printf("🛠️  Calling tool: %s args=%v\n", toolCall.Function.Name, toolCall.Function.Arguments)

			// Execute the tool
			callToolParam := &mcp.CallToolParams{
				Name:      toolCall.Function.Name,
				Arguments: json.RawMessage(toolCall.Function.Arguments),
			}

			result, err := a.session.CallTool(ctx, callToolParam)
			if err != nil {
				slog.Error("Error calling tool", "tool", toolCall.Function.Name, "err", err)
				continue
			}

			// Add tool result to conversation history
			for _, c := range result.Content {
				text := c.(*mcp.TextContent).Text
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    "tool",
					Content: text,
					ToolCalls: []openai.ChatCompletionMessageToolCall{{
						ID: toolCall.ID,
					}},
				})
			}
		}
	}

	fmt.Println()
	fmt.Println("🎉 result:")
	fmt.Println(finalResult)
	return nil
}

// chatCompletionMessage extracts the content and tool calls from the first choice.
func (a *App) groupChatCompletionChoices(choices []openai.ChatCompletionChoice) openai.ChatCompletionMessage {
	if len(choices) == 0 {
		return openai.ChatCompletionMessage{}
	}

	return openai.ChatCompletionMessage{
		Content:   choices[0].Message.Content,
		ToolCalls: choices[0].Message.ToolCalls,
	}

	// var sb strings.Builder
	// for _, choice := range choices {
	// 	if choice.Message.Content != "" {
	// 		sb.WriteString(choice.Message.Content)
	// 		sb.WriteString("\n")
	// 	}
	// }

	// return openai.ChatCompletionMessage{
	// 	Content:   sb.String(),
	// 	ToolCalls: choices[0].Message.ToolCalls,
	// }
}
