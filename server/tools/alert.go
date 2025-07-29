package tools

import (
	"context"
	"encoding/json"
	"strings"
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

	var alerts []string
	for _, f := range data.Features {
		alerts = append(alerts, formatAlert(f))
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(alerts, "\n")}},
	}, nil
}

func formatAlert(f dtos.Feature) string {
	lines := []string{
		"Event: " + defaultString(f.AlertProperties.Event, "Unknown"),
		"Area: " + defaultString(f.AlertProperties.AreaDesc, "Unknown"),
		"Severity: " + defaultString(f.AlertProperties.Severity, "Unknown"),
		"Description: " + defaultString(f.AlertProperties.Description, "No description available"),
		"Instructions: " + defaultString(f.AlertProperties.Instruction, "No specific instructions provided"),
	}
	return strings.Join(lines, "\n")
}

func defaultString(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
