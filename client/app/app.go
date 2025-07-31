// app/app.go
package app

import (
	"context"
	"fmt"
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
	session      *mcp.ClientSession
	tools        []openai.ChatCompletionToolParam
	openAIEngine *engine.OpenAIClient
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

	// Create a command-based transport to run the MCP server as a child process via stdin/stdout
	// transport := mcp.NewCommandTransport(exec.Command(serverPath))

	// Create a streamable client transport to communicate with the MCP server
	transport := mcp.NewStreamableClientTransport(
		"http://localhost:8080/mcp/stream",
		&mcp.StreamableClientTransportOptions{},
	)

	// Use SSE client transport
	// transport := mcp.NewSSEClientTransport("http://localhost:8080/mcp/streamm", &mcp.SSEClientTransportOptions{})

	// Connect to the server and start a session
	session, err := a.client.Connect(ctx, transport)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	a.session = session

	// Initialize OpenAI client using the API key from config
	openAIengine := engine.NewOpenAIClient(a.conf.Anthropic.APIKey, engine.MODEL)
	a.openAIEngine = openAIengine

	// Start the main prompt loop
	return a.StartPromptLoop()
}
