package mcp

import (
	"log"
	"net/http"

	"github.com/akshaygalande/google-air-quality-mcp/internal/capabilities"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPServer struct {
	server *mcp.Server
}

func NewMCPServer(name string, version string) *MCPServer {
	// Create a server
	s := mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil)

	// Register all features (tools, prompts, resources)
	capabilities.RegisterAll(s)

	return &MCPServer{
		server: s,
	}
}

func (s *MCPServer) SetupStreamableHTTP(r *gin.Engine) {
	// Create the Streamable HTTP handler
	// We pass a function that returns the server instance.
	// Assuming the server instance is safe to reuse or handles sessions internally.
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return s.server
	}, nil)

	// Register the handler for the /mcp endpoint (or whatever path is preferred)
	// Streamable HTTP typically uses a single endpoint or a base path.
	// The handler likely handles the sub-paths or query params.
	// We'll register it at /mcp/*path to catch everything if needed, or just /mcp
	// But usually it's a single endpoint.
	// Let's assume it handles standard HTTP requests.

	r.Any("/mcp", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})

	log.Println("Streamable HTTP endpoint registered at /mcp")
}
