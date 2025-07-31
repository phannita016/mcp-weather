package main

import (
	"log/slog"
	"os"
	"weather/server/srv"
)

func main() {
	slog.Info("Starting weather MCP server...")

	// Create a new server instance. Tool registration is now handled within NewServer.
	srv := srv.NewServer()

	// To run with stdio transport for command-line usage, comment out the line
	// above and uncomment the one below:
	// if err := srv.RunStdio(); err != nil {
	// 	slog.Error("failed to start server", "error", err)
	// }

	// To run with HTTP transport, use this line:
	if err := srv.RunHTTP(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	// To run with SSE HTTP transport
	// if err := srv.RunSSE(); err != nil {
	// 	os.Exit(1)
	// }
}
