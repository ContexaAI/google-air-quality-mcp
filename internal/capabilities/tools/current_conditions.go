package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	CurrentConditionsToolName        = "get_current_air_quality"
	CurrentConditionsToolDescription = "Get current air quality conditions for a specific location. Returns air quality indexes, pollutant levels, and health recommendations."
)

var CurrentConditionsToolSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"latitude": map[string]interface{}{
			"type":        "number",
			"description": "Location latitude",
		},
		"longitude": map[string]interface{}{
			"type":        "number",
			"description": "Location longitude",
		},
		"extraComputations": map[string]interface{}{
			"type":        "array",
			"description": "Additional features to compute (LOCAL_AQI HEALTH_RECOMMENDATIONS POLLUTANT_ADDITIONAL_INFO DOMINANT_POLLUTANT_CONCENTRATION POLLUTANT_CONCENTRATION)",
			"items": map[string]interface{}{
				"type": "string",
			},
		},
		"uaqiColorPalette": map[string]interface{}{
			"type":        "string",
			"description": "Color palette for UAQI (RED_GREEN INDIGO_PERSIAN NUMERIC)",
		},
		"universalAqi": map[string]interface{}{
			"type":        "boolean",
			"description": "Include Universal AQI (default: true)",
		},
		"languageCode": map[string]interface{}{
			"type":        "string",
			"description": "Response language code (default: en)",
		},
	},
	"required": []interface{}{"latitude", "longitude"},
}

// CurrentConditionsInput defines the input for the current conditions tool
type CurrentConditionsInput struct {
	Latitude          float64  `json:"latitude" jsonschema:"required,description=Location latitude"`
	Longitude         float64  `json:"longitude" jsonschema:"required,description=Location longitude"`
	ExtraComputations []string `json:"extraComputations,omitempty" jsonschema:"description=Additional features to compute (LOCAL_AQI HEALTH_RECOMMENDATIONS POLLUTANT_ADDITIONAL_INFO DOMINANT_POLLUTANT_CONCENTRATION POLLUTANT_CONCENTRATION)"`
	UaqiColorPalette  string   `json:"uaqiColorPalette,omitempty" jsonschema:"description=Color palette for UAQI (RED_GREEN INDIGO_PERSIAN NUMERIC)"`
	UniversalAqi      *bool    `json:"universalAqi,omitempty" jsonschema:"description=Include Universal AQI (default: true)"`
	LanguageCode      string   `json:"languageCode,omitempty" jsonschema:"description=Response language code (default: en)"`
}

// CurrentConditionsOutput defines the output for the current conditions tool
type CurrentConditionsOutput struct {
	Response CurrentConditionsResponse `json:"response"`
}

// NewCurrentConditionsHandler creates a new current conditions handler with the API key
func NewCurrentConditionsHandler(apiKey string) func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse input from request arguments
		var input CurrentConditionsInput
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

		// Build request
		req := CurrentConditionsRequest{
			Location: LatLng{
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
			LanguageCode: input.LanguageCode,
			UniversalAqi: input.UniversalAqi,
		}

		// Convert extra computations
		for _, comp := range input.ExtraComputations {
			req.ExtraComputations = append(req.ExtraComputations, ExtraComputation(comp))
		}

		// Set color palette if provided
		if input.UaqiColorPalette != "" {
			req.UaqiColorPalette = ColorPalette(input.UaqiColorPalette)
		}

		// Call API
		client := NewClient(apiKey)
		resp, err := client.GetCurrentConditions(req)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get current conditions: %v", err)}},
			}, nil
		}

		// Convert response to JSON string
		jsonResp, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to marshal response: %v", err)}},
			}, nil
		}

		// Return success result with JSON response
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(jsonResp)}},
		}, nil
	}
}
