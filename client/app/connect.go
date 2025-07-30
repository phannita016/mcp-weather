package app

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverPath = "e:/GOLANG-LAB/MCP/mcp-weather/server/cmd/weather.exe"

func (a *App) Connect() error {
	// transport := mcp.NewCommandTransport(exec.Command(serverPath))
	transport := mcp.NewStreamableClientTransport("http://localhost:8080/mcp/stream", &mcp.StreamableClientTransportOptions{})
	session, err := a.client.Connect(context.Background(), transport)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	a.session = session
	return nil
}

// func (a *App) Connect(serverPath string) error {
// 	absPath, err := filepath.Abs(serverPath)
// 	if err != nil {
// 		return fmt.Errorf("cannot get absolute path: %w", err)
// 	}
// 	fmt.Println("Resolved absolute path:", absPath)

// 	transport := mcp.NewCommandTransport(exec.Command(absPath))
// 	session, err := a.client.Connect(context.Background(), transport)
// 	if err != nil {
// 		return fmt.Errorf("connect failed: %w", err)
// 	}
// 	a.session = session
// 	return nil
// }
