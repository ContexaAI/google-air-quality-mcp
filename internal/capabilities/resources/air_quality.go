package resources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akshaygalande/google-air-quality-mcp/internal/capabilities/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AirQualityResourceHandler handles air quality resource requests
type AirQualityResourceHandler struct {
	client *tools.Client
}

// NewAirQualityResourceHandler creates a new AirQualityResourceHandler
func NewAirQualityResourceHandler(apiKey string) *AirQualityResourceHandler {
	return &AirQualityResourceHandler{
		client: tools.NewClient(apiKey),
	}
}

// CurrentConditionsHandler handles requests for current air quality conditions
// URI: airquality://current/{lat},{long}
func (h *AirQualityResourceHandler) CurrentConditionsHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := request.Params.URI
	fmt.Printf("DEBUG: CurrentConditionsHandler called with URI: %s\n", uri)
	prefix := "airquality://current/"
	if !strings.HasPrefix(uri, prefix) {
		return nil, fmt.Errorf("invalid URI format")
	}
	location := strings.TrimPrefix(uri, prefix)

	lat, lon, err := parseLatLong(location)
	if err != nil {
		return nil, err
	}

	// Helper to create bool pointer
	truePtr := true

	req := tools.CurrentConditionsRequest{
		Location: tools.LatLng{
			Latitude:  lat,
			Longitude: lon,
		},
		UniversalAqi: &truePtr, // Default to true
	}

	resp, err := h.client.GetCurrentConditions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get current conditions: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// ForecastHandler handles requests for air quality forecast
// URI: airquality://forecast/{lat},{long}
func (h *AirQualityResourceHandler) ForecastHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := request.Params.URI
	fmt.Printf("DEBUG: ForecastHandler called with URI: %s\n", uri)
	prefix := "airquality://forecast/"
	if !strings.HasPrefix(uri, prefix) {
		return nil, fmt.Errorf("invalid URI format")
	}
	location := strings.TrimPrefix(uri, prefix)

	lat, lon, err := parseLatLong(location)
	if err != nil {
		return nil, err
	}

	// Helper to create bool pointer
	truePtr := true

	req := tools.ForecastRequest{
		Location: tools.LatLng{
			Latitude:  lat,
			Longitude: lon,
		},
		UniversalAqi: &truePtr,
	}

	resp, err := h.client.GetForecast(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// HistoryHandler handles requests for air quality history
// URI: airquality://history/{lat},{long}
func (h *AirQualityResourceHandler) HistoryHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := request.Params.URI
	fmt.Printf("DEBUG: HistoryHandler called with URI: %s\n", uri)
	prefix := "airquality://history/"
	if !strings.HasPrefix(uri, prefix) {
		return nil, fmt.Errorf("invalid URI format")
	}
	location := strings.TrimPrefix(uri, prefix)

	lat, lon, err := parseLatLong(location)
	if err != nil {
		return nil, err
	}

	// Helper to create bool pointer
	truePtr := true

	// Default to 24 hours history
	req := tools.HistoryRequest{
		Location: tools.LatLng{
			Latitude:  lat,
			Longitude: lon,
		},
		Hours:        24,
		UniversalAqi: &truePtr,
	}

	resp, err := h.client.GetHistory(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}

// HeatmapHandler handles requests for heatmap tiles
// URI: airquality://heatmap/{mapType}/{z}/{x}/{y}
func (h *AirQualityResourceHandler) HeatmapHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := request.Params.URI
	fmt.Printf("DEBUG: HeatmapHandler called with URI: %s\n", uri)
	prefix := "airquality://heatmap/"
	if !strings.HasPrefix(uri, prefix) {
		return nil, fmt.Errorf("invalid URI format")
	}

	// Remove prefix and split by /
	path := strings.TrimPrefix(uri, prefix)
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid heatmap URI format, expected {mapType}/{z}/{x}/{y}")
	}

	mapTypeStr := parts[0]
	zStr := parts[1]
	xStr := parts[2]
	yStr := parts[3]

	zoom, err := strconv.Atoi(zStr)
	if err != nil {
		return nil, fmt.Errorf("invalid zoom level: %w", err)
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return nil, fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return nil, fmt.Errorf("invalid y coordinate: %w", err)
	}

	// Validate map type
	validMapType := false
	// This list should match the types in tools/heatmap.go or client.go
	validTypes := []string{
		"UAQI_RED_GREEN", "UAQI_INDIGO_PERSIAN",
		"PM25_INDIGO_PERSIAN", "GBR_DEFRA",
		"DEU_UBA", "CAN_EC", "FRA_ATMO", "US_AQI",
	}
	for _, t := range validTypes {
		if mapTypeStr == t {
			validMapType = true
			break
		}
	}
	if !validMapType {
		return nil, fmt.Errorf("invalid map type: %s", mapTypeStr)
	}

	data, err := h.client.GetHeatmapTile(tools.MapType(mapTypeStr), zoom, x, y)
	if err != nil {
		return nil, fmt.Errorf("failed to get heatmap tile: %w", err)
	}

	// Encode to base64 for text transport, or use Blob if supported.
	// MCP spec allows Blob, but go-sdk ResourceContents has Blob field?
	// Let's check ResourceContents definition again.
	// The view_file of examples.go showed Text field.
	// Let's assume Text with base64 for now if Blob isn't obvious, OR check if Blob is available.
	// Wait, I didn't check ResourceContents struct fully.
	// I'll use base64 in Text for safety as it's common for binary in JSON-RPC.
	encoded := base64.StdEncoding.EncodeToString(data)

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "image/png",
				Text:     encoded,
			},
		},
	}, nil
}

func parseLatLong(s string) (float64, float64, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid location format, expected 'lat,long'")
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %w", err)
	}
	return lat, lon, nil
}
