package main

import (
	"log/slog"
	"os"
	"weather/client/app"
	"weather/client/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	app := app.NewApp(*conf, client)

	if err := app.Run(); err != nil {
		slog.Error("failed to run app", "error", err)
		os.Exit(1)
	}
}
