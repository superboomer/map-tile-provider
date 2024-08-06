package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/tile"
)

// MockProviderSchema simulates a provider schema for testing
var MockProviderSchema = schema{
	Name:       "MockProvider",
	ID:         "mp",
	MaxJobs:    100,
	MaxZoom:    19,
	Projection: "wgs84",
	Request: reqSchema{
		URL: "https://example.com/{x}/{y}/{z}.png",
	},
}

// TestCreateProvider tests the creation of a MapProvider instance, including error cases
func TestCreateProvider(t *testing.T) {
	// Test successful creation
	provider, err := createProvider(&MockProviderSchema)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, MockProviderSchema.Name, provider.Name())
	assert.Equal(t, MockProviderSchema.ID, provider.ID())

	// invalidProjectionSchema simulates a provider schema with an unsupported projection
	var invalidProjectionSchema = schema{
		Name:       "InvalidProjectionProvider",
		ID:         "ip",
		MaxJobs:    100,
		MaxZoom:    19,
		Projection: "unsupported_projection_type",
		Request: reqSchema{
			URL: "https://example.com/data",
		},
	}

	// Test unsupported projection
	_, err = createProvider(&invalidProjectionSchema)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "projection")
}

// TestGetTile tests the GetTile method of MapProvider for various edge cases
func TestGetTile(t *testing.T) {
	provider, _ := createProvider(&MockProviderSchema)

	// Test with valid inputs
	lat, long, scale := 0.0, 0.0, 1.0 // Example coordinates and scale
	testTile := provider.GetTile(lat, long, scale)
	assert.NotNil(t, testTile)
	assert.NotZero(t, testTile.X)
	assert.NotZero(t, testTile.Y)
	assert.NotZero(t, testTile.Z)

	// Test with maximum zoom level
	lat, long, scale = 0.0, 0.0, float64(MockProviderSchema.MaxZoom)
	testTile = provider.GetTile(lat, long, scale)
	assert.NotNil(t, testTile)
	assert.NotZero(t, testTile.X)
	assert.NotZero(t, testTile.Y)

	// Test with extremely high scale (to check for overflow or incorrect calculations)
	lat, long, scale = 0.0, 0.0, float64(^uint(0)>>1) // Half of uint max value
	testTile = provider.GetTile(lat, long, scale)
	assert.NotNil(t, testTile)
	assert.NotZero(t, testTile.X)
	assert.NotZero(t, testTile.Y)
	assert.NotZero(t, testTile.Z)
}

func TestGetRequest(t *testing.T) {
	provider, _ := createProvider(&MockProviderSchema)

	// Test with valid tile coordinates
	testTileValid := &tile.Tile{X: 123, Y: 456, Z: 7} // Example tile within valid range
	req := provider.GetRequest(testTileValid)
	assert.NotNil(t, req)
	assert.Equal(t, "https://example.com/123/456/7.png", req.URL.String())
}

func TestGetRequestWithHeaders(t *testing.T) {
	// MockProviderSchemaWithHeaders is assumed to be modified to include Headers
	var MockProviderSchemaWithHeaders = schema{
		Name:       "MockProvider",
		ID:         "mp",
		MaxJobs:    100,
		MaxZoom:    19,
		Projection: "wgs84",
		Request: reqSchema{
			URL: "https://example.com/{x}/{y}/{z}.png",
			Headers: []headersSchema{
				{"Authorization", "Bearer YOUR_TOKEN_HERE"},
				{"Content-Type", "image/png"},
			},
		},
	}
	provider, _ := createProvider(&MockProviderSchemaWithHeaders)

	// Test with valid tile coordinates
	testTileValid := &tile.Tile{X: 123, Y: 456, Z: 7} // Example tile within valid range
	req := provider.GetRequest(testTileValid)

	assert.NotNil(t, req)
	assert.Equal(t, req.Header.Get("Authorization"), MockProviderSchemaWithHeaders.Request.Headers[0].Value)
	assert.Equal(t, req.Header.Get("Content-Type"), MockProviderSchemaWithHeaders.Request.Headers[1].Value)
}
