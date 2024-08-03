package provider

import (
	"net/http"
	"testing"

	"github.com/superboomer/map-tile-provider/app/tile"
)

// MockProvider is a mock implementation of the Provider interface for testing.
type MockProvider struct{}

func (mp *MockProvider) GetTile(lat, long, scale float64) tile.Tile {
	return tile.Tile{} // Return an empty tile for testing purposes.
}

func (mp *MockProvider) MaxJobs() int {
	return 1 // Return a fixed number of jobs for testing.
}

func (mp *MockProvider) MaxZoom() int {
	return 20 // Return a fixed number of jobs for testing.
}

func (mp *MockProvider) Name() string {
	return "name"
}

func (mp *MockProvider) GetRequest(t *tile.Tile) *http.Request {
	return &http.Request{} // Return a dummy request for testing.
}

// TestCreateProviderList tests the CreateProviderList function.
func TestCreateProviderList(t *testing.T) {
	pl := CreateProviderList()
	if pl == nil {
		t.Error("Expected non-nil ProviderList")
	}
}

// TestRegister tests the Register method of ProviderList.
func TestRegister(t *testing.T) {
	pl := CreateProviderList()
	mockProvider := &MockProvider{}

	// Test registering a new provider
	err := pl.Register(mockProvider)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test registering the same provider again
	err = pl.Register(mockProvider)
	if err == nil {
		t.Error("Expected error when registering existing provider, got nil")
	} else if err.Error() != "provider provider1 already exist" {
		t.Errorf("Expected error 'provider provider1 already exist', got %v", err)
	}
}

// TestGet tests the Get method of ProviderList.
func TestGet(t *testing.T) {
	pl := CreateProviderList()
	mockProvider := &MockProvider{}

	// Register a provider
	pl.Register(mockProvider)

	// Test getting an existing provider
	provider, err := pl.Get("provider1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected non-nil provider")
	}

	// Test getting a non-existing provider
	_, err = pl.Get("non_existing_provider")
	if err == nil {
		t.Error("Expected error when getting non-existing provider, got nil")
	} else if err.Error() != "provider non_existing_provider not found" {
		t.Errorf("Expected error 'provider non_existing_provider not found', got %v", err)
	}
}

func TestGetAllNames(t *testing.T) {
	// Create a new ProviderList
	pl := CreateProviderList()

	// Register some mock providers
	pl.Register(&MockProvider{})
	pl.Register(&MockProvider{})
	pl.Register(&MockProvider{})

	// Call GetAllNames
	names := pl.GetAllNames()

	// Define the expected names
	expectedNames := []string{"ProviderA", "ProviderB", "ProviderC"}

	// Check if the result matches the expected names
	if len(names) != len(expectedNames) {
		t.Errorf("Expected %d names, got %d", len(expectedNames), len(names))
	}

	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expectedName := range expectedNames {
		if !nameMap[expectedName] {
			t.Errorf("Expected name %s not found in result", expectedName)
		}
	}
}

func TestGetAllNamesEmpty(t *testing.T) {
	// Create a new ProviderList
	pl := CreateProviderList()

	// Call GetAllNames on an empty ProviderList
	names := pl.GetAllNames()

	// Check if the result is an empty slice
	if len(names) != 0 {
		t.Error("Expected GetAllNames to return an empty slice for an empty ProviderList")
	}
}
