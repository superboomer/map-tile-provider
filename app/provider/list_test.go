package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapList_Register(t *testing.T) {
	list := createProviderList()

	mockProvider := &ProviderMock{
		IDFunc:   func() string { return "provider" },
		NameFunc: func() string { return "provider_name" },
	}

	err := list.Register(mockProvider)
	assert.NoError(t, err)

	// Verify that the provider was added
	_, err = list.Get(mockProvider.ID())
	assert.NoError(t, err)

	err = list.Register(mockProvider)
	assert.Error(t, err)
	assert.EqualError(t, err, "provider provider_name (provider) already exist")
}

func TestMapList_Get(t *testing.T) {
	list := createProviderList()

	mockProvider := &ProviderMock{
		IDFunc:   func() string { return "provider" },
		NameFunc: func() string { return "provider_name" },
	}

	err := list.Register(mockProvider)
	assert.NoError(t, err)

	// Test getting the provider by ID
	retrievedProvider, err := list.Get(mockProvider.ID())
	assert.NoError(t, err)
	assert.Equal(t, mockProvider.ID(), retrievedProvider.ID())
	assert.Equal(t, mockProvider.Name(), retrievedProvider.Name())

	// Test getting a non-existing provider
	_, err = list.Get("nonExistentId")
	assert.Error(t, err)
	assert.EqualError(t, err, "provider nonExistentId not found")
}

func TestMapList_GetAllID(t *testing.T) {
	list := createProviderList()

	mockProvider := &ProviderMock{
		IDFunc:   func() string { return "provider" },
		NameFunc: func() string { return "provider_name" },
	}

	err := list.Register(mockProvider)
	assert.NoError(t, err)

	mockProvider2 := &ProviderMock{
		IDFunc:   func() string { return "provider2" },
		NameFunc: func() string { return "provider_name2" },
	}

	err = list.Register(mockProvider2)
	assert.NoError(t, err)

	// Test getting all IDs
	ids := list.GetAllID()
	assert.Contains(t, ids, mockProvider.ID())
	assert.Contains(t, ids, mockProvider2.ID())
}
