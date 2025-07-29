package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"weather/server/dtos"
	"weather/server/nws"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GetForecast(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[dtos.ForecastParams]) (*mcp.CallToolResultFor[any], error) {
	pointsURL := nws.GetForecastURL(params.Arguments.Latitude, params.Arguments.Longitude)

	pointsBody, err := nws.MakeNWSRequest(pointsURL)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch forecast data for this location."}},
		}, nil
	}

	pointsData := dtos.PointsData{}
	if err := json.Unmarshal(pointsBody, &pointsData); err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Unable to parse forecast point data."}},
		}, nil
	}

	forecastBody, err := nws.MakeNWSRequest(pointsData.Properties.ForecastURL)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch detailed forecast."}},
		}, nil
	}

	forecastData := dtos.ForecastData{}
	if err := json.Unmarshal(forecastBody, &forecastData); err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Unable to parse forecast data."}},
		}, nil
	}

	var forecasts []string
	for i, period := range forecastData.Properties.Periods {
		if i >= 3 {
			break
		}
		forecasts = append(forecasts, formatPeriod(period))
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(forecasts, "\n")}},
	}, nil
}

func formatPeriod(period dtos.ForecastPeriod) string {
	lines := []string{
		"- " + defaultString(period.Name, "Unknown") + ":",
		fmt.Sprintf("   Temperature: %d°%s", period.Temperature, defaultString(period.TemperatureUnit, "")),
		"   Wind: " + defaultString(period.WindSpeed, "") + " " + defaultString(period.WindDirection, ""),
		"   Forecast: " + defaultString(period.DetailedForecast, "No forecast available"),
	}
	return strings.Join(lines, "\n")
}
