package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/gorilla/mux"
)

// DatasetAPI - An interface used to access the DatasetAPI
type DatasetAPI interface {
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}
type CantabularMetadataExtractorAPI struct {
	Router *mux.Router
	datasetAPI	DatasetAPI
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, d DatasetAPI) *CantabularMetadataExtractorAPI {
	api := &CantabularMetadataExtractorAPI{
		Router: r,
		datasetAPI: d,
	}

	r.HandleFunc("/metadata", api.getMetadata).Methods("GET")
	return api
}
