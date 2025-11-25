package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	HistoryToolName        = "get_air_quality_history"
	HistoryToolDescription = "Get historical air quality data for a specific location. Returns past hourly air quality measurements."
)

var HistoryToolSchema = map[string]interface{}{
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
			"description": "Max hourly records per page (default: 72 max: 168)",
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
			"description": "Specific historical time (ISO 8601 format)",
		},
		"hours": map[string]interface{}{
			"type":        "integer",
			"description": "Number of hours of history",
		},
		"periodStartTime": map[string]interface{}{
			"type":        "string",
			"description": "History period start time (ISO 8601 format)",
		},
		"periodEndTime": map[string]interface{}{
			"type":        "string",
			"description": "History period end time (ISO 8601 format)",
		},
	},
	"required": []interface{}{"latitude", "longitude"},
}

// HistoryInput defines the input for the history tool
type HistoryInput struct {
	Latitude          float64  `json:"latitude" jsonschema:"required,description=Location latitude"`
	Longitude         float64  `json:"longitude" jsonschema:"required,description=Location longitude"`
	PageSize          int      `json:"pageSize,omitempty" jsonschema:"description=Max hourly records per page (default: 72 max: 168)"`
	PageToken         string   `json:"pageToken,omitempty" jsonschema:"description=Pagination token for next page"`
	ExtraComputations []string `json:"extraComputations,omitempty" jsonschema:"description=Additional features to compute"`
	UaqiColorPalette  string   `json:"uaqiColorPalette,omitempty" jsonschema:"description=Color palette for UAQI"`
	UniversalAqi      *bool    `json:"universalAqi,omitempty" jsonschema:"description=Include Universal AQI (default: true)"`
	LanguageCode      string   `json:"languageCode,omitempty" jsonschema:"description=Response language code (default: en)"`
	DateTime          string   `json:"dateTime,omitempty" jsonschema:"description=Specific historical time (ISO 8601 format)"`
	Hours             int      `json:"hours,omitempty" jsonschema:"description=Number of hours of history"`
	PeriodStartTime   string   `json:"periodStartTime,omitempty" jsonschema:"description=History period start time (ISO 8601 format)"`
	PeriodEndTime     string   `json:"periodEndTime,omitempty" jsonschema:"description=History period end time (ISO 8601 format)"`
}

// HistoryOutput defines the output for the history tool
type HistoryOutput struct {
	Response HistoryResponse `json:"response"`
}

// NewHistoryHandler creates a new history handler with the API key
func NewHistoryHandler(apiKey string) func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse input from request arguments
		var input HistoryInput
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
		req := HistoryRequest{
			Location: LatLng{
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
			PageSize:     input.PageSize,
			PageToken:    input.PageToken,
			LanguageCode: input.LanguageCode,
			UniversalAqi: input.UniversalAqi,
			DateTime:     input.DateTime,
			Hours:        input.Hours,
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
		resp, err := client.GetHistory(req)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get history: %v", err)}},
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
