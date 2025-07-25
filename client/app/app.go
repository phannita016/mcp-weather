// app/app.go
package app

import (
	"log/slog"
	"os"
	"weather/client/dtos"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type App struct {
	client  *mcp.Client
	session *mcp.ClientSession
	tools   []dtos.Tool
}

func NewApp(client *mcp.Client) *App {
	return &App{client: client}
}

func (a *App) Run() error {
	if err := a.Connect(); err != nil {
		slog.Error("failed to connect to server", "error", err)
		os.Exit(1)
	}

	return a.StartPromptLoop()
}
