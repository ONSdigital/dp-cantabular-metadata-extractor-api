package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/gorilla/mux"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
)

// DatasetAPI - An interface used to access the DatasetAPI
type DatasetAPI interface {
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}
type CantabularMetadataExtractorAPI struct {
	Router *mux.Router
	datasetAPI	DatasetAPI
	cfg  *config.Config
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, config  *config.Config, d DatasetAPI) *CantabularMetadataExtractorAPI {
	api := &CantabularMetadataExtractorAPI{
		Router: r,
		datasetAPI: d,
		cfg:config,
	}

	r.HandleFunc("/metadata/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", api.getMetadata).Methods("GET")
	return api
}
