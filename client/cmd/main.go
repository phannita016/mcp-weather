package main

import (
	"log/slog"
	"os"
	"weather/client/app"
	"weather/client/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load configuration from the config file
	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize MCP client with basic information (name & version)
	opt := mcp.ClientOptions{}
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, &opt)

	// Create a new App instance with the loaded configuration and client
	app := app.NewApp(*conf, client)

	// Run the application and handle any errors
	if err := app.Run(); err != nil {
		slog.Error("failed to run app", "error", err)
		os.Exit(1)
	}
}
