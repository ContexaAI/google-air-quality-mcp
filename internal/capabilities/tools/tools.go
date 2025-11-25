package tools

import (
	"github.com/akshaygalande/google-air-quality-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all tools with the MCP server
func RegisterAll(server *mcp.Server) {
	// Load configuration to get API key
	cfg := config.LoadConfig()

	// Register Air Quality API tools
	server.AddTool(&mcp.Tool{
		Name:        CurrentConditionsToolName,
		Description: CurrentConditionsToolDescription,
		InputSchema: CurrentConditionsToolSchema,
	}, NewCurrentConditionsHandler(cfg.APIKey))

	server.AddTool(&mcp.Tool{
		Name:        ForecastToolName,
		Description: ForecastToolDescription,
		InputSchema: ForecastToolSchema,
	}, NewForecastHandler(cfg.APIKey))

	server.AddTool(&mcp.Tool{
		Name:        HistoryToolName,
		Description: HistoryToolDescription,
		InputSchema: HistoryToolSchema,
	}, NewHistoryHandler(cfg.APIKey))

	server.AddTool(&mcp.Tool{
		Name:        HeatmapToolName,
		Description: HeatmapToolDescription,
		InputSchema: HeatmapToolSchema,
	}, NewHeatmapHandler(cfg.APIKey))
}
