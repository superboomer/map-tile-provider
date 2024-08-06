// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package cache

import (
	"github.com/superboomer/map-tile-provider/app/tile"
	"sync"
)

// Ensure, that CacheMock does implement Cache.
// If this is not the case, regenerate this file with moq.
var _ Cache = &CacheMock{}

// CacheMock is a mock implementation of Cache.
//
//	func TestSomethingThatUsesCache(t *testing.T) {
//
//		// make and configure a mocked Cache
//		mockedCache := &CacheMock{
//			LoadTileFunc: func(vendor string, t *tile.Tile) ([]byte, error) {
//				panic("mock out the LoadTile method")
//			},
//			SaveTileFunc: func(vendor string, t *tile.Tile) error {
//				panic("mock out the SaveTile method")
//			},
//		}
//
//		// use mockedCache in code that requires Cache
//		// and then make assertions.
//
//	}
type CacheMock struct {
	// LoadTileFunc mocks the LoadTile method.
	LoadTileFunc func(vendor string, t *tile.Tile) ([]byte, error)

	// SaveTileFunc mocks the SaveTile method.
	SaveTileFunc func(vendor string, t *tile.Tile) error

	// calls tracks calls to the methods.
	calls struct {
		// LoadTile holds details about calls to the LoadTile method.
		LoadTile []struct {
			// Vendor is the vendor argument value.
			Vendor string
			// T is the t argument value.
			T *tile.Tile
		}
		// SaveTile holds details about calls to the SaveTile method.
		SaveTile []struct {
			// Vendor is the vendor argument value.
			Vendor string
			// T is the t argument value.
			T *tile.Tile
		}
	}
	lockLoadTile sync.RWMutex
	lockSaveTile sync.RWMutex
}

// LoadTile calls LoadTileFunc.
func (mock *CacheMock) LoadTile(vendor string, t *tile.Tile) ([]byte, error) {
	if mock.LoadTileFunc == nil {
		panic("CacheMock.LoadTileFunc: method is nil but Cache.LoadTile was just called")
	}
	callInfo := struct {
		Vendor string
		T      *tile.Tile
	}{
		Vendor: vendor,
		T:      t,
	}
	mock.lockLoadTile.Lock()
	mock.calls.LoadTile = append(mock.calls.LoadTile, callInfo)
	mock.lockLoadTile.Unlock()
	return mock.LoadTileFunc(vendor, t)
}

// LoadTileCalls gets all the calls that were made to LoadTile.
// Check the length with:
//
//	len(mockedCache.LoadTileCalls())
func (mock *CacheMock) LoadTileCalls() []struct {
	Vendor string
	T      *tile.Tile
} {
	var calls []struct {
		Vendor string
		T      *tile.Tile
	}
	mock.lockLoadTile.RLock()
	calls = mock.calls.LoadTile
	mock.lockLoadTile.RUnlock()
	return calls
}

// SaveTile calls SaveTileFunc.
func (mock *CacheMock) SaveTile(vendor string, t *tile.Tile) error {
	if mock.SaveTileFunc == nil {
		panic("CacheMock.SaveTileFunc: method is nil but Cache.SaveTile was just called")
	}
	callInfo := struct {
		Vendor string
		T      *tile.Tile
	}{
		Vendor: vendor,
		T:      t,
	}
	mock.lockSaveTile.Lock()
	mock.calls.SaveTile = append(mock.calls.SaveTile, callInfo)
	mock.lockSaveTile.Unlock()
	return mock.SaveTileFunc(vendor, t)
}

// SaveTileCalls gets all the calls that were made to SaveTile.
// Check the length with:
//
//	len(mockedCache.SaveTileCalls())
func (mock *CacheMock) SaveTileCalls() []struct {
	Vendor string
	T      *tile.Tile
} {
	var calls []struct {
		Vendor string
		T      *tile.Tile
	}
	mock.lockSaveTile.RLock()
	calls = mock.calls.SaveTile
	mock.lockSaveTile.RUnlock()
	return calls
}
