// app/app.go
package app

import (
	"log/slog"
	"os"
	"weather/client/dtos"
	"weather/client/engine"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type App struct {
	client       *mcp.Client
	session      *mcp.ClientSession
	tools        []dtos.Tool
	openAIEngine *engine.OpenAIClient
}

func NewApp(client *mcp.Client) *App {
	openAIengine := engine.NewOpenAIClient(engine.ANTHROPIC_API_KEY, engine.MODEL)

	return &App{
		client:       client,
		openAIEngine: openAIengine,
	}
}

func (a *App) Run() error {
	if err := a.Connect(); err != nil {
		slog.Error("failed to connect to server", "error", err)
		os.Exit(1)
	}

	return a.StartPromptLoop()
}
