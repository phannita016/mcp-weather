// app/app.go
package app

import (
	"context"
	"fmt"
	"os/exec"
	"time"
	"weather/client/config"
	"weather/client/engine"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go"
)

const serverPath = "e:/GOLANG-LAB/MCP/mcp-weather/server/cmd/weather.exe"

type App struct {
	conf         config.Config
	client       *mcp.Client
	openAIEngine *engine.OpenAIClient
	mcps         map[string]*mcp.ClientSession
	tools        []openai.ChatCompletionToolParam
}

func NewApp(conf config.Config, client *mcp.Client) *App {
	return &App{
		conf:   conf,
		client: client,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use SSE client transport
	// transport := mcp.NewSSEClientTransport("http://localhost:8080/mcp/stream", &mcp.SSEClientTransportOptions{})

	// Create a command-based transport to run the MCP server as a child process via stdin/stdout
	cmd := exec.Command("npx", "-y", "mcp-echarts")
	echartsTransport := mcp.NewCommandTransport(cmd)

	echartsSession, err := a.client.Connect(ctx, echartsTransport)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	// Create a streamable client transport to communicate with the MCP server
	waetherTransport := mcp.NewStreamableClientTransport(
		"http://localhost:8080/mcp/stream",
		&mcp.StreamableClientTransportOptions{},
	)

	weatherSession, err := a.client.Connect(ctx, waetherTransport)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	a.mcps = map[string]*mcp.ClientSession{
		"echarts": echartsSession,
		"weather": weatherSession,
	}

	// Initialize OpenAI client using the API key from config
	openAIengine := engine.NewOpenAIClient(a.conf.Anthropic.APIKey, engine.MODEL)
	a.openAIEngine = openAIengine

	// Start the main prompt loop
	return a.StartPromptLoop()
}
