package main

import (
	"log/slog"
	"os"
	"weather/client/app"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	app := app.NewApp(client)

	if err := app.Run(); err != nil {
		slog.Error("failed to run app", "error", err)
		os.Exit(1)
	}
}
