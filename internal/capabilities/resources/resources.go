package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all resources with the MCP server
func RegisterAll(server *mcp.Server) {
	// Register a simple static resource as an example
	server.AddResource(&mcp.Resource{
		URI:         "example://server-info",
		Name:        "Server Information",
		Description: "Basic information about this MCP server",
		MIMEType:    "text/plain",
	}, ServerInfoHandler)

}

// ServerInfoHandler provides basic server information as a simple example resource
func ServerInfoHandler(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	info := `Google Air Quality MCP Server
==============================

This server provides air quality data through the Google Air Quality API.

Available Resources:
- airquality://current/{lat},{long} - Current air quality conditions
- airquality://forecast/{lat},{long} - Air quality forecast
- airquality://history/{lat},{long} - Historical air quality data
- airquality://heatmap/{mapType}/{zoom}/{x}/{y} - Heatmap tiles

Example: airquality://current/37.7749,-122.4194
`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     info,
			},
		},
	}, nil
}
