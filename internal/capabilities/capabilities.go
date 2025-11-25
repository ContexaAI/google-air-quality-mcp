package capabilities

import (
	"github.com/akshaygalande/google-air-quality-mcp/internal/capabilities/prompts"
	"github.com/akshaygalande/google-air-quality-mcp/internal/capabilities/resources"
	"github.com/akshaygalande/google-air-quality-mcp/internal/capabilities/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all MCP features (tools, prompts, resources) with the server
func RegisterAll(server *mcp.Server) {
	// Register all tools
	tools.RegisterAll(server)

	// Register all prompts
	prompts.RegisterAll(server)

	// Register all resources
	resources.RegisterAll(server)
}
