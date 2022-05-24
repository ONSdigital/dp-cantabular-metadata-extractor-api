package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
)

//go:generate moq -out mock/datasetapi.go -pkg mock . DatasetAPI

/*
// DatasetAPI - An interface used to access the DatasetAPI
type DatasetAPI interface {
	GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (m dataset.VersionDimensions, err error)
}
*/

// cantExtAPI
type CantExtAPI interface {
	// does interface need var names? XXX
	GetMetadataDataset(cantDataset string, dimensions []string) *cantabular.MetadataDatasetQuery
	GetMetadataTable(datasetID string) (*cantabular.MetadataTableQuery, []string, error)
}

type CantabularMetadataExtractorAPI struct { // XXX
	Router     *mux.Router
	CantExtAPI CantExtAPI
	Cfg        *config.Config
}

//Setup function sets up the api and returns an api
// STM TODO add my stuff here
func Setup(ctx context.Context, r *mux.Router, config *config.Config, c CantExtAPI) *CantabularMetadataExtractorAPI {
	//func Setup(ctx context.Context, r *mux.Router, config *config.Config, d DatasetAPI) *CantabularMetadataExtractorAPI {
	api := &CantabularMetadataExtractorAPI{
		Router: r,
		//		DatasetAPI: d,
		CantExtAPI: c,
		Cfg:        config,
	}

	r.HandleFunc("/dataset/{datasetID}/cantabular/{cantdataset}/lang/{lang}", api.getMetadata).Methods("GET")
	//r. // XXXHandleFunc("/metadata/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", api.getMetadata).Methods("GET")
	return api
}
