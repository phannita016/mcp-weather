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
	"weather/client/engine"

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

		if err := a.handleQuery(ctx, query); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) handleQuery(ctx context.Context, query string) error {
	messages := []dtos.Message{{Role: "user", Content: query}}

	toolCalls, _, err := engine.OpenAiEngine(messages, a.tools)
	if err != nil {
		slog.Error("Error generating tool call", "error", err)
		return err
	}

	for _, toolCall := range toolCalls {
		name := toolCall["name"].(string)
		args := toolCall["arguments"].(map[string]interface{})

		fmt.Printf("Calling tool: %s with args: %+v\n", name, args)

		result, err := a.session.CallTool(ctx, &mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		})
		if err != nil {
			slog.Error("Error calling tool", "error", err)
			continue
		}

		b, err := json.MarshalIndent(result.Content, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling result:", err)
		} else {
			fmt.Println("Tool result:", string(b))
			fmt.Println()
		}
	}

	if len(toolCalls) == 0 {
		_, content, err := engine.OpenAiEngine(messages, a.tools)
		if err != nil {
			slog.Error("Error getting AI response", "error", err)
			return err
		}
		fmt.Println("Tool result:", content)
	}

	return nil
}
