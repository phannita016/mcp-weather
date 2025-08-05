package srv

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RunStdio starts the server and communicates over standard input/output.
// This is the typical mode for a command-line MCP plugin.
// It blocks until the server is stopped.
func (s *Server) RunStdio() error {
	slog.Info("Running server with stdio transport")
	if err := s.mcpServer.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		return err
	}

	return nil
}

// RunHTTP starts the server and listens for connections on the given HTTP address.
// It uses Server-Sent Events (SSE) for transport.
// This is useful for web-based clients.
func (s *Server) RunHTTP(addr string) error {
	transport := mcp.NewStreamableServerTransport("")
	http.Handle("/mcp/stream", transport)
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	// Run the MCP server in a separate goroutine to not block the HTTP server.
	go func() {
		if err := s.mcpServer.Run(context.Background(), transport); err != nil {
			slog.Error("MCP server run failed", "error", err)
			// os.Exit(1)
			// return err
		}
	}()

	slog.Info("Starting HTTP server", "address", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("failed to start HTTP server", "error", err)
		// os.Exit(1)
		return err
	}
	return nil
}

// RunSSE starts the SSE server that can handle multiple MCP server instances
func (s *Server) RunSSE() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path != "/mcp/stream" {
			slog.Warn("No MCP server registered for path", "path", path)
			http.Error(w, "invalid MCP path", http.StatusNotFound)
			return
		}

		sseHandler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			slog.Info("New MCP connection", "path", path)
			return s.MCP()
		})

		sseT := mcp.NewSSEServerTransport("/mcp/stream", nil)
		sseT.ServeHTTP(w, r)

		sseHandler.ServeHTTP(w, r)
	})

	if handler == nil {
		return fmt.Errorf("failed to create SSE handler")
	}

	slog.Info("Starting SSE server", "address", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("Failed to start SSE server", "error", err)
		return err
	}

	return nil
}

// func (s *Server) RunHTTP(addr string) error {
// 	transport := mcp.NewStreamableServerTransport("")
// 	mux := http.NewServeMux()
// 	mux.Handle("/mcp/stream", transport)

// 	httpServer := &http.Server{
// 		Addr:    addr,
// 		Handler: mux,
// 	}

// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	errChan := make(chan error, 1)

// 	// Run the MCP server in a separate goroutine to not block the HTTP server.
// 	go func() {
// 		if err := s.mcpServer.Run(ctx, transport); err != nil {
// 			errChan <- err
// 		}
// 	}()

// 	// Run the HTTP server in a separate goroutine.
// 	go func() {
// 		slog.Info("Starting HTTP server", "address", addr)
// 		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			errChan <- err
// 		}
// 	}()

// 	// Wait for an interrupt signal or an error from one of the goroutines.
// 	select {
// 	case err := <-errChan:
// 		return err
// 	case <-ctx.Done():
// 		slog.Info("Shutting down server...")
// 		// Create a new context for shutdown to allow for a timeout.
// 		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		if err := httpServer.Shutdown(shutdownCtx); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
