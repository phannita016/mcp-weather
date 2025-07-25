package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"weather/server/dtos"
	"weather/server/nws"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetAlerts fetches active weather alerts for a given US state from the NWS API.
func GetAlerts(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[dtos.AlertsParams]) (*mcp.CallToolResultFor[any], error) {
	url := nws.GetAlertsURL(params.Arguments.State)

	body, err := nws.MakeNWSRequest(url)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch alerts or no alerts found."}},
		}, nil
	}

	data := dtos.FeatureCollection{}
	if err := json.Unmarshal(body, &data); err != nil || len(data.Features) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No active alerts for this state."}},
		}, nil
	}

	var result string
	for _, f := range data.Features {
		result += fmt.Sprintf("Event: %s || Area: %s \n---\n", f.AlertProperties.Event, f.AlertProperties.AreaDesc)
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result}},
	}, nil
}
