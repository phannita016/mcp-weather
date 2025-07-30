// app/app.go
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
	"weather/client/config"
	"weather/client/dtos"
	"weather/client/engine"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverPath = "e:/GOLANG-LAB/MCP/mcp-weather/server/cmd/weather.exe"

type App struct {
	client       *mcp.Client
	session      *mcp.ClientSession
	tools        []dtos.Tool
	openAIEngine *engine.OpenAIClient
}

func NewApp(conf config.Config, client *mcp.Client) *App {
	openAIengine := engine.NewOpenAIClient(conf.Anthropic.APIKey, engine.MODEL)

	return &App{
		client:       client,
		openAIEngine: openAIengine,
	}
}

func (a *App) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transport := mcp.NewStreamableClientTransport(
		"http://localhost:8080/mcp/stream",
		&mcp.StreamableClientTransportOptions{},
	)

	session, err := a.client.Connect(ctx, transport)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	a.session = session
	return nil
}

func (a *App) Run() error {
	if err := a.Connect(context.Background()); err != nil {
		slog.Error("failed to connect to server", "error", err)
		os.Exit(1)
	}

	return a.StartPromptLoop()
}
