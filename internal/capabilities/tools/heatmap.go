package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	HeatmapToolName        = "get_air_quality_heatmap_tile"
	HeatmapToolDescription = "Get air quality heatmap tile image for visualization. Returns a PNG image tile for the specified map type and coordinates."
)

var HeatmapToolSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"mapType": map[string]interface{}{
			"type":        "string",
			"description": "Type of heatmap (UAQI_RED_GREEN UAQI_INDIGO_PERSIAN PM25_INDIGO_PERSIAN GBR_DEFRA DEU_UBA CAN_EC FRA_ATMO US_AQI)",
		},
		"zoom": map[string]interface{}{
			"type":        "integer",
			"description": "Zoom level (0-16)",
		},
		"x": map[string]interface{}{
			"type":        "integer",
			"description": "East-west tile coordinate",
		},
		"y": map[string]interface{}{
			"type":        "integer",
			"description": "North-south tile coordinate",
		},
	},
	"required": []interface{}{"mapType", "zoom", "x", "y"},
}

// HeatmapInput defines the input for the heatmap tile tool
type HeatmapInput struct {
	MapType string `json:"mapType" jsonschema:"required,description=Type of heatmap (UAQI_RED_GREEN UAQI_INDIGO_PERSIAN PM25_INDIGO_PERSIAN GBR_DEFRA DEU_UBA CAN_EC FRA_ATMO US_AQI)"`
	Zoom    int    `json:"zoom" jsonschema:"required,description=Zoom level (0-16)"`
	X       int    `json:"x" jsonschema:"required,description=East-west tile coordinate"`
	Y       int    `json:"y" jsonschema:"required,description=North-south tile coordinate"`
}

// HeatmapOutput defines the output for the heatmap tile tool
type HeatmapOutput struct {
	ImageData string `json:"imageData" jsonschema:"description=Base64 encoded PNG image data"`
}

// NewHeatmapHandler creates a new heatmap handler with the API key
func NewHeatmapHandler(apiKey string) func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse input from request arguments
		var input HeatmapInput
		if request.Params.Arguments != nil {
			// Convert map to JSON and back to struct
			jsonData, err := json.Marshal(request.Params.Arguments)
			if err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to marshal arguments: %v", err)}},
				}, nil
			}
			if err := json.Unmarshal(jsonData, &input); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid input: %v", err)}},
				}, nil
			}
		}

		// Validate zoom level
		if input.Zoom < 0 || input.Zoom > 16 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: "zoom must be between 0 and 16"}},
			}, nil
		}

		// Call API
		client := NewClient(apiKey)
		imageData, err := client.GetHeatmapTile(MapType(input.MapType), input.Zoom, input.X, input.Y)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get heatmap tile: %v", err)}},
			}, nil
		}

		// Encode image as base64
		base64Data := base64.StdEncoding.EncodeToString(imageData)
		result := fmt.Sprintf("data:image/png;base64,%s", base64Data)

		// Return success result with base64 image
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, nil
	}
}
