package srv

import (
	"weather/server/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcpServer *mcp.Server
}

// NewServer creates and initializes a new Server instance.
// It sets up the underlying MCP server with the necessary implementation details.
func NewServer() *Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "weather",
		Version: "v1.0.0",
	}, nil)

	s := &Server{
		mcpServer: mcpServer,
	}

	s.registerTools()

	return s
}

func (s *Server) MCP() *mcp.Server {
	return s.mcpServer
}

// registerTools registers all the available tools with the MCP server.
func (s *Server) registerTools() {
	// Tool: get_alerts
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_alerts",
		Description: "Get active weather alerts for a given US state",
	}, tools.GetAlerts)

	// Tool: get_forecast
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get weather forecast for a given location",
	}, tools.GetForecast)
}
