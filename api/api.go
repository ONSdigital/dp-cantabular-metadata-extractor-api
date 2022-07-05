package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
)

//go:generate moq -out mock/datasetapi.go -pkg mock . CantExtAPI

// cantExtAPI
type CantExtAPI interface {
	MetadataDatasetQuery(ctx context.Context, req cantabular.MetadataDatasetQueryRequest) (*cantabular.MetadataDatasetQuery, error)
	MetadataTableQuery(ctx context.Context, req cantabular.MetadataTableQueryRequest) (*cantabular.MetadataTableQuery, error)
}

type CantabularMetadataExtractorAPI struct {
	Router     *mux.Router
	CantExtAPI CantExtAPI
	Cfg        *config.Config
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, config *config.Config, c CantExtAPI) *CantabularMetadataExtractorAPI {

	api := &CantabularMetadataExtractorAPI{
		Router:     r,
		CantExtAPI: c,
		Cfg:        config,
	}

	r.HandleFunc("/cantabular-metadata/dataset/{datasetID}/lang/{lang}", api.getMetadata).Methods("GET")
	return api
}
