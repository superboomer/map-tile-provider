package provider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/provider"
)

func TestLoadProviderList_Success(t *testing.T) {
	list, err := provider.LoadProviderList("./../../example/providers.json")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, list)
}

func TestLoadProviderList_ErrorLoadingJSON(t *testing.T) {
	list, err := provider.LoadProviderList("invalid/path")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error occurred when loading providers schema")
	assert.Nil(t, list)
}

func TestLoadProviderList_ErrorCreatingProvider(t *testing.T) {
	// Execute
	list, err := provider.LoadProviderList("./testdata/providers_invalid.json")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "projection unknown not found for provider")
	assert.Nil(t, list)
}

func TestLoadProviderList_ErrorRegisteringProvider(t *testing.T) {
	// Execute
	list, err := provider.LoadProviderList("./testdata/providers_duplicate.json")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error occurred when registering new provider")
	assert.Nil(t, list)
}
