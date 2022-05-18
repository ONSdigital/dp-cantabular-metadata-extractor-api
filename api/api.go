package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/metadata"
	"github.com/gorilla/mux"
)

//go:generate moq -out mock/datasetapi.go -pkg mock . DatasetAPI

// DatasetAPI - An interface used to access the DatasetAPI
type DatasetAPI interface {
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}

// cantExtAPI
type cantExtAPI interface {
	GetMetaData(cantDataset string, dimensions []string) (resp metadata.Resp)
}

type CantabularMetadataExtractorAPI struct {
	Router     *mux.Router
	cantExtAPI cantExtAPI
	DatasetAPI DatasetAPI
	cfg        *config.Config
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, config *config.Config, d DatasetAPI) *CantabularMetadataExtractorAPI {
	api := &CantabularMetadataExtractorAPI{
		Router:     r,
		DatasetAPI: d,
		cfg:        config,
	}

	r.HandleFunc("/metadata/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", api.getMetadata).Methods("GET")
	return api
}
