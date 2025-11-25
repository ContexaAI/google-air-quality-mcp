package prompts

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all prompts with the MCP server
func RegisterAll(server *mcp.Server) {
	// Register prompt for air quality heatmap tile tool (direct coordinates)
	server.AddPrompt(&mcp.Prompt{
		Name:        "air_quality_heatmap_prompt",
		Description: "Get heatmap tile image for air quality visualization",
		Arguments: []*mcp.PromptArgument{
			{Name: "mapType", Description: "Type of heatmap (e.g., UAQI_RED_GREEN)", Required: true},
			{Name: "x", Description: "Tile X coordinate", Required: true},
			{Name: "y", Description: "Tile Y coordinate", Required: true},
			{Name: "zoom", Description: "Zoom level (0-16)", Required: true},
		},
	}, AirQualityHeatmapHandler)

	// Prompt for current air quality by location name
	server.AddPrompt(&mcp.Prompt{
		Name:        "current_air_quality_by_location_prompt",
		Description: "Get current air quality conditions by providing a location name (LLM will convert to latitude/longitude)",
		Arguments:   []*mcp.PromptArgument{{Name: "location", Description: "Human readable location name (e.g., 'Paris, France')", Required: true}},
	}, CurrentAirQualityByLocationHandler)

	// Prompt for air quality forecast by location name
	server.AddPrompt(&mcp.Prompt{
		Name:        "air_quality_forecast_by_location_prompt",
		Description: "Get air quality forecast for a location name (LLM will convert to latitude/longitude)",
		Arguments:   []*mcp.PromptArgument{{Name: "location", Description: "Human readable location name", Required: true}, {Name: "pageSize", Description: "Number of forecast hours to return (optional)", Required: false}},
	}, AirQualityForecastByLocationHandler)

	// Prompt for air quality history by location name
	server.AddPrompt(&mcp.Prompt{
		Name:        "air_quality_history_by_location_prompt",
		Description: "Get historical air quality data for a location name (LLM will convert to latitude/longitude)",
		Arguments:   []*mcp.PromptArgument{{Name: "location", Description: "Human readable location name", Required: true}, {Name: "hours", Description: "Number of past hours to retrieve (optional)", Required: false}},
	}, AirQualityHistoryByLocationHandler)

	// Prompt for air quality heatmap tile by location name
	server.AddPrompt(&mcp.Prompt{
		Name:        "air_quality_heatmap_by_location_prompt",
		Description: "Get heatmap tile image for a location name (LLM will convert to latitude/longitude and tile coordinates)",
		Arguments:   []*mcp.PromptArgument{{Name: "location", Description: "Human readable location name", Required: true}, {Name: "mapType", Description: "Type of heatmap (e.g., UAQI_RED_GREEN)", Required: true}, {Name: "zoom", Description: "Zoom level (0-16)", Required: true}},
	}, AirQualityHeatmapByLocationHandler)
}
