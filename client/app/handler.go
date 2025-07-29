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

	toolCalls, content, err := engine.OpenAiEngine(messages, a.tools)
	if err != nil {
		slog.Error("Error generating tool call", "error", err)
		return err
	}

	finalResults := []string{}
	for _, toolCall := range toolCalls {
		name := toolCall.Name
		args := toolCall.Arguments

		var parsedArgs map[string]interface{}
		if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
			slog.Error("Error unmarshaling args", "error", err)
			return err
		}

		fmt.Printf("🛠️  Calling tool: name=%s args=%v\n", name, parsedArgs)

		result, err := a.session.CallTool(ctx, &mcp.CallToolParams{
			Name:      name,
			Arguments: parsedArgs,
		})
		if err != nil {
			slog.Error("Error calling tool", "error", err)
			continue
		}

		for _, c := range result.Content {
			text := c.(*mcp.TextContent).Text
			finalResults = append(finalResults, text)
		}
	}

	if content != "" {
		finalResults = append(finalResults, content)
	}

	fmt.Println("🎉 Final combined result:")
	for i, r := range finalResults {
		fmt.Printf("[Chunk %d]\n", i+1)
		fmt.Println(r)
		fmt.Println("------------------------------------------------")
	}

	return nil
}
