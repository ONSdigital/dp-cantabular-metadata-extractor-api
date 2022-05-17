// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package api

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"sync"
)

// Ensure, that DatasetAPIMock does implement DatasetAPI.
// If this is not the case, regenerate this file with moq.
var _ DatasetAPI = &DatasetAPIMock{}

// DatasetAPIMock is a mock implementation of DatasetAPI.
//
// 	func TestSomethingThatUsesDatasetAPI(t *testing.T) {
//
// 		// make and configure a mocked DatasetAPI
// 		mockedDatasetAPI := &DatasetAPIMock{
// 			GetVersionDimensionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
// 				panic("mock out the GetVersionDimensions method")
// 			},
// 		}
//
// 		// use mockedDatasetAPI in code that requires DatasetAPI
// 		// and then make assertions.
//
// 	}
type DatasetAPIMock struct {
	// GetVersionDimensionsFunc mocks the GetVersionDimensions method.
	GetVersionDimensionsFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetVersionDimensions holds details about calls to the GetVersionDimensions method.
		GetVersionDimensions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
	}
	lockGetVersionDimensions sync.RWMutex
}

// GetVersionDimensions calls GetVersionDimensionsFunc.
func (mock *DatasetAPIMock) GetVersionDimensions(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
	if mock.GetVersionDimensionsFunc == nil {
		panic("DatasetAPIMock.GetVersionDimensionsFunc: method is nil but DatasetAPI.GetVersionDimensions was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		ID:               id,
		Edition:          edition,
		Version:          version,
	}
	mock.lockGetVersionDimensions.Lock()
	mock.calls.GetVersionDimensions = append(mock.calls.GetVersionDimensions, callInfo)
	mock.lockGetVersionDimensions.Unlock()
	return mock.GetVersionDimensionsFunc(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
}

// GetVersionDimensionsCalls gets all the calls that were made to GetVersionDimensions.
// Check the length with:
//     len(mockedDatasetAPI.GetVersionDimensionsCalls())
func (mock *DatasetAPIMock) GetVersionDimensionsCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	ID               string
	Edition          string
	Version          string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}
	mock.lockGetVersionDimensions.RLock()
	calls = mock.calls.GetVersionDimensions
	mock.lockGetVersionDimensions.RUnlock()
	return calls
}
