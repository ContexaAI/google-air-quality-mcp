# Google Air Quality MCP Server

A Model Context Protocol (MCP) server that provides access to Google's Air Quality API, enabling AI assistants to retrieve real-time air quality data, forecasts, historical information, and heatmap visualizations.

## Features

- 🌍 **Real-time Air Quality Data** - Current conditions for any location worldwide
- 📊 **Forecasts** - Hourly air quality predictions
- 📈 **Historical Data** - Access past air quality measurements
- 🗺️ **Heatmap Visualizations** - Generate air quality heatmap tiles
- 🤖 **LLM-Friendly Prompts** - Natural language location queries
- 📦 **MCP Resources** - Structured data access via URI templates

## Prerequisites

- Go 1.21 or higher
- Google Air Quality API key (see [API Key Setup](#api-key-setup))

## Installation

1. **Clone the repository**
```bash
git clone git@github.com:ContexaAI/google-air-quality-mcp.git
cd google-air-quality-mcp
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment variables**

Create a `.env` file in the project root:
```env
GOOGLE_AIR_QUALITY_API_KEY=your_api_key_here
MCP_SERVER_NAME=Google Air Quality MCP Server
PORT=8080
```

## API Key Setup

### Step 1: Enable the API

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to **APIs & Services** > **Library**
4. Search for "Air Quality API"
5. Click **Enable**

### Step 2: Create API Credentials

1. Go to **APIs & Services** > **Credentials**
2. Click **Create Credentials** > **API Key**
3. Copy the generated API key
4. (Recommended) Click **Restrict Key** and:
   - Under **API restrictions**, select "Restrict key"
   - Choose "Air Quality API" from the dropdown
   - Under **Application restrictions**, you can restrict by IP, HTTP referrer, etc.

### Step 3: Configure the Server

Add your API key to the `.env` file:
```env
GOOGLE_AIR_QUALITY_API_KEY=AIzaSyD...your-key-here
```

> **Note**: The Air Quality API may require billing to be enabled on your Google Cloud project. Check the [pricing page](https://developers.google.com/maps/documentation/air-quality/usage-and-billing) for details.

## Usage

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` with the MCP endpoint at `/mcp`.

### Connecting with MCP Inspector

1. Open [MCP Inspector](https://inspector.modelcontextprotocol.io/)
2. Connect to `http://localhost:8080/mcp`
3. Explore available tools, prompts, and resources

## API Reference

### Tools

| Tool Name | Description | Required Parameters | Optional Parameters |
|-----------|-------------|---------------------|---------------------|
| `get_current_air_quality` | Get current air quality conditions for a specific location | `latitude` (float)<br>`longitude` (float) | `universalAqi` (bool)<br>`languageCode` (string)<br>`extraComputations` (array)<br>`uaqiColorPalette` (string) |
| `get_air_quality_forecast` | Get hourly air quality forecast predictions | `latitude` (float)<br>`longitude` (float) | `pageSize` (int)<br>`pageToken` (string)<br>`universalAqi` (bool)<br>`languageCode` (string)<br>`extraComputations` (array) |
| `get_air_quality_history` | Get historical air quality data | `latitude` (float)<br>`longitude` (float) | `hours` (int)<br>`pageSize` (int)<br>`pageToken` (string)<br>`universalAqi` (bool)<br>`languageCode` (string) |
| `get_air_quality_heatmap_tile` | Get heatmap tile image for visualization | `mapType` (string)<br>`zoom` (int)<br>`x` (int)<br>`y` (int) | - |

#### Valid Map Types for Heatmap
- `UAQI_RED_GREEN` - Universal AQI with red-green color palette
- `UAQI_INDIGO_PERSIAN` - Universal AQI with indigo-persian palette
- `PM25_INDIGO_PERSIAN` - PM2.5 concentration heatmap
- `GBR_DEFRA` - UK DEFRA index
- `DEU_UBA` - German UBA index
- `CAN_EC` - Canadian AQHI
- `FRA_ATMO` - French ATMO index
- `US_AQI` - US EPA AQI

### Prompts

| Prompt Name | Description | Arguments |
|-------------|-------------|-----------|
| `current_air_quality_by_location_prompt` | Get current air quality using a human-readable location name (LLM resolves to coordinates) | `location` (string) - e.g., "Paris, France" |
| `air_quality_forecast_by_location_prompt` | Get air quality forecast for a location name | `location` (string)<br>`pageSize` (int, optional) |
| `air_quality_history_by_location_prompt` | Get historical air quality for a location name | `location` (string)<br>`hours` (int, optional) |
| `air_quality_heatmap_by_location_prompt` | Get heatmap tile for a location name | `location` (string)<br>`mapType` (string)<br>`zoom` (int) |
| `air_quality_heatmap_prompt` | Get heatmap tile using direct coordinates | `mapType` (string)<br>`x` (int)<br>`y` (int)<br>`zoom` (int) |

### Resources

| Resource URI | Type | Description | Content Type |
|--------------|------|-------------|--------------|
| `example://server-info` | Static | Basic server information and available resources | `text/plain` |

## Examples

### Using Tools

**Get current air quality for San Francisco:**
```json
{
  "name": "get_current_air_quality",
  "arguments": {
    "latitude": 37.7749,
    "longitude": -122.4194,
    "universalAqi": true
  }
}
```

**Get 24-hour forecast for New York:**
```json
{
  "name": "get_air_quality_forecast",
  "arguments": {
    "latitude": 40.7128,
    "longitude": -74.0060,
    "pageSize": 24
  }
}
```

**Get historical data for London:**
```json
{
  "name": "get_air_quality_history",
  "arguments": {
    "latitude": 51.5074,
    "longitude": -0.1278,
    "hours": 72
  }
}
```

### Using Prompts

Prompts allow you to use natural language location names. The LLM will automatically resolve them to coordinates.

**Example prompt usage:**
```
Use the current_air_quality_by_location_prompt with location "Tokyo, Japan"
```

The LLM will:
1. Resolve "Tokyo, Japan" to coordinates (35.6762, 139.6503)
2. Call the `get_current_air_quality` tool with those coordinates
3. Return the air quality data

### Using Resources

**Read server information:**
```
URI: example://server-info
```

This returns basic information about the server and available resources.

## Project Structure

```
.
├── cmd/
│   └── server/              # Server entry point
├── internal/
│   ├── capabilities/        # MCP capabilities
│   │   ├── prompts/        # Prompt definitions and handlers
│   │   ├── resources/      # Resource definitions and handlers
│   │   └── tools/          # Tool implementations
│   │       ├── client.go   # Google Air Quality API client
│   │       ├── types.go    # Shared types and structures
│   │       ├── current_conditions.go
│   │       ├── forecast.go
│   │       ├── history.go
│   │       └── heatmap.go
│   ├── config/             # Configuration management
│   └── mcp/                # MCP server setup
├── .env                    # Environment variables (not in git)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Development

### Building

```bash
go build -o bin/server cmd/server/main.go
```

### Running

```bash
./bin/server
```

### Testing

```bash
go test ./...
```

## Configuration

Environment variables (set in `.env`):

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GOOGLE_AIR_QUALITY_API_KEY` | Your Google Air Quality API key | - | ✅ Yes |
| `MCP_SERVER_NAME` | Display name for the MCP server | `Google Air Quality MCP Server` | No |
| `PORT` | HTTP server port | `8080` | No |

## Troubleshooting

### API Key Issues

**Error: "API key not valid"**
- Verify your API key is correct in `.env`
- Ensure the Air Quality API is enabled in Google Cloud Console
- Check that API restrictions (if any) allow requests from your IP

**Error: "This API method requires billing to be enabled"**
- Enable billing on your Google Cloud project
- The Air Quality API may have usage costs

### Connection Issues

**Error: "address already in use"**
- Another process is using port 8080
- Kill the process: `npx kill-port 8080`
- Or change the `PORT` in `.env`

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/)
- [Google Air Quality API Documentation](https://developers.google.com/maps/documentation/air-quality)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Inspector](https://inspector.modelcontextprotocol.io/)

## Support

For issues and questions:
- Open an issue on [GitHub](https://github.com/ContexaAI/google-air-quality-mcp/issues)
- Check the [Google Air Quality API documentation](https://developers.google.com/maps/documentation/air-quality)
