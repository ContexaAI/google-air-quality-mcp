package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ForecastToolName        = "get_air_quality_forecast"
	ForecastToolDescription = "Get air quality forecast for a specific location. Returns hourly forecasts with air quality indexes and pollutant predictions."
)

var ForecastToolSchema = map[string]interface{}{
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
		"pageSize": map[string]interface{}{
			"type":        "integer",
			"description": "Number of forecast hours to return",
		},
		"pageToken": map[string]interface{}{
			"type":        "string",
			"description": "Pagination token for next page",
		},
		"extraComputations": map[string]interface{}{
			"type":        "array",
			"description": "Additional features to compute",
			"items": map[string]interface{}{
				"type": "string",
			},
		},
		"uaqiColorPalette": map[string]interface{}{
			"type":        "string",
			"description": "Color palette for UAQI",
		},
		"universalAqi": map[string]interface{}{
			"type":        "boolean",
			"description": "Include Universal AQI (default: true)",
		},
		"languageCode": map[string]interface{}{
			"type":        "string",
			"description": "Response language code (default: en)",
		},
		"dateTime": map[string]interface{}{
			"type":        "string",
			"description": "Specific forecast time (ISO 8601 format)",
		},
		"periodStartTime": map[string]interface{}{
			"type":        "string",
			"description": "Forecast period start time (ISO 8601 format)",
		},
		"periodEndTime": map[string]interface{}{
			"type":        "string",
			"description": "Forecast period end time (ISO 8601 format)",
		},
	},
	"required": []interface{}{"latitude", "longitude"},
}

// ForecastInput defines the input for the forecast tool
type ForecastInput struct {
	Latitude          float64  `json:"latitude" jsonschema:"required,description=Location latitude"`
	Longitude         float64  `json:"longitude" jsonschema:"required,description=Location longitude"`
	PageSize          int      `json:"pageSize,omitempty" jsonschema:"description=Number of forecast hours to return"`
	PageToken         string   `json:"pageToken,omitempty" jsonschema:"description=Pagination token for next page"`
	ExtraComputations []string `json:"extraComputations,omitempty" jsonschema:"description=Additional features to compute"`
	UaqiColorPalette  string   `json:"uaqiColorPalette,omitempty" jsonschema:"description=Color palette for UAQI"`
	UniversalAqi      *bool    `json:"universalAqi,omitempty" jsonschema:"description=Include Universal AQI (default: true)"`
	LanguageCode      string   `json:"languageCode,omitempty" jsonschema:"description=Response language code (default: en)"`
	DateTime          string   `json:"dateTime,omitempty" jsonschema:"description=Specific forecast time (ISO 8601 format)"`
	PeriodStartTime   string   `json:"periodStartTime,omitempty" jsonschema:"description=Forecast period start time (ISO 8601 format)"`
	PeriodEndTime     string   `json:"periodEndTime,omitempty" jsonschema:"description=Forecast period end time (ISO 8601 format)"`
}

// ForecastOutput defines the output for the forecast tool
type ForecastOutput struct {
	Response ForecastResponse `json:"response"`
}

// NewForecastHandler creates a new forecast handler with the API key
func NewForecastHandler(apiKey string) func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse input from request arguments
		var input ForecastInput
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
		req := ForecastRequest{
			Location: LatLng{
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
			PageSize:     input.PageSize,
			PageToken:    input.PageToken,
			LanguageCode: input.LanguageCode,
			UniversalAqi: input.UniversalAqi,
			DateTime:     input.DateTime,
		}

		// Convert extra computations
		for _, comp := range input.ExtraComputations {
			req.ExtraComputations = append(req.ExtraComputations, ExtraComputation(comp))
		}

		// Set color palette if provided
		if input.UaqiColorPalette != "" {
			req.UaqiColorPalette = ColorPalette(input.UaqiColorPalette)
		}

		// Set period if provided
		if input.PeriodStartTime != "" || input.PeriodEndTime != "" {
			req.Period = &Interval{
				StartTime: input.PeriodStartTime,
				EndTime:   input.PeriodEndTime,
			}
		}

		// Call API
		client := NewClient(apiKey)
		resp, err := client.GetForecast(req)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get forecast: %v", err)}},
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
