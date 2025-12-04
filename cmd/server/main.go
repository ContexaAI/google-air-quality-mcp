package main

import (
	"log"

	"github.com/akshaygalande/google-air-quality-mcp/internal/config"
	"github.com/akshaygalande/google-air-quality-mcp/internal/mcp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Gin
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "mcp-session-id"},
		ExposeHeaders:   []string{"mcp-session-id"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Google Air Quality MCP Server",
		})
	})

	// Initialize MCP Server and setup Streamable HTTP
	mcpServer := mcp.NewMCPServer(cfg.MCPServerName, "0.1.0")
	mcpServer.SetupStreamableHTTP(r)

	log.Printf("Starting Gin server with MCP Streamable HTTP on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Gin server error: %v", err)
	}
}
