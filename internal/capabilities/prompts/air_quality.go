package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CurrentAirQualityByLocationHandler handles the prompt for current air quality by location
func CurrentAirQualityByLocationHandler(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := "unknown location"
	if request.Params.Arguments != nil {
		if loc, ok := request.Params.Arguments["location"]; ok {
			location = loc
		}
	}

	return &mcp.GetPromptResult{
		Description: "Prompt to get current air quality for a location",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf("Please get the current air quality for %s. You should first determine the latitude and longitude for this location, and then use the 'get_current_air_quality' tool.", location),
				},
			},
		},
	}, nil
}

// AirQualityForecastByLocationHandler handles the prompt for air quality forecast by location
func AirQualityForecastByLocationHandler(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := "unknown location"
	pageSize := ""
	if request.Params.Arguments != nil {
		if loc, ok := request.Params.Arguments["location"]; ok {
			location = loc
		}
		if ps, ok := request.Params.Arguments["pageSize"]; ok {
			pageSize = ps
		}
	}

	promptText := fmt.Sprintf("Please get the air quality forecast for %s.", location)
	if pageSize != "" {
		promptText += fmt.Sprintf(" Please retrieve %s hours of forecast data.", pageSize)
	}
	promptText += " You should first determine the latitude and longitude for this location, and then use the 'get_air_quality_forecast' tool."

	return &mcp.GetPromptResult{
		Description: "Prompt to get air quality forecast for a location",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: promptText,
				},
			},
		},
	}, nil
}

// AirQualityHistoryByLocationHandler handles the prompt for air quality history by location
func AirQualityHistoryByLocationHandler(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := "unknown location"
	hours := ""
	if request.Params.Arguments != nil {
		if loc, ok := request.Params.Arguments["location"]; ok {
			location = loc
		}
		if h, ok := request.Params.Arguments["hours"]; ok {
			hours = h
		}
	}

	promptText := fmt.Sprintf("Please get the historical air quality data for %s.", location)
	if hours != "" {
		promptText += fmt.Sprintf(" Please retrieve data for the past %s hours.", hours)
	}
	promptText += " You should first determine the latitude and longitude for this location, and then use the 'get_air_quality_history' tool."

	return &mcp.GetPromptResult{
		Description: "Prompt to get historical air quality data for a location",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: promptText,
				},
			},
		},
	}, nil
}

// AirQualityHeatmapByLocationHandler handles the prompt for air quality heatmap by location
func AirQualityHeatmapByLocationHandler(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	location := "unknown location"
	mapType := "UAQI_RED_GREEN"
	zoom := "10"

	if request.Params.Arguments != nil {
		if loc, ok := request.Params.Arguments["location"]; ok {
			location = loc
		}
		if mt, ok := request.Params.Arguments["mapType"]; ok {
			mapType = mt
		}
		if z, ok := request.Params.Arguments["zoom"]; ok {
			zoom = z
		}
	}

	promptText := fmt.Sprintf("Please get the air quality heatmap tile for %s using map type '%s' at zoom level %s.", location, mapType, zoom)
	promptText += " You should first determine the latitude and longitude for this location, convert them to the appropriate X and Y tile coordinates for the given zoom level, and then use the 'get_air_quality_heatmap_tile' tool."

	return &mcp.GetPromptResult{
		Description: "Prompt to get air quality heatmap tile for a location",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: promptText,
				},
			},
		},
	}, nil
}

// AirQualityHeatmapHandler handles the prompt for air quality heatmap tile (direct coordinates)
func AirQualityHeatmapHandler(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	mapType := "UAQI_RED_GREEN"
	x := "0"
	y := "0"
	zoom := "0"

	if request.Params.Arguments != nil {
		if mt, ok := request.Params.Arguments["mapType"]; ok {
			mapType = mt
		}
		if val, ok := request.Params.Arguments["x"]; ok {
			x = val
		}
		if val, ok := request.Params.Arguments["y"]; ok {
			y = val
		}
		if val, ok := request.Params.Arguments["zoom"]; ok {
			zoom = val
		}
	}

	return &mcp.GetPromptResult{
		Description: "Prompt to get air quality heatmap tile",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf("Please get the air quality heatmap tile for map type '%s' at coordinates x=%s, y=%s, zoom=%s using the 'get_air_quality_heatmap_tile' tool.", mapType, x, y, zoom),
				},
			},
		},
	}, nil
}
