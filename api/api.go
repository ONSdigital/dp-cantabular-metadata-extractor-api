package api

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
	"github.com/ryboe/q"
)

//go:generate moq -out mock/datasetapi.go -pkg mock . CantExtAPI

// cantExtAPI
type CantExtAPI interface {
	MetadataDatasetQuery(ctx context.Context, req cantabular.MetadataDatasetQueryRequest) (*cantabular.MetadataDatasetQuery, error)
	MetadataTableQuery(ctx context.Context, req cantabular.MetadataTableQueryRequest) (*cantabular.MetadataTableQuery, error)
}

type CantabularMetadataExtractorAPI struct { // XXX
	Router     *mux.Router
	CantExtAPI CantExtAPI
	Cfg        *config.Config
}

//Setup function sets up the api and returns an api
// STM TODO add my stuff here
func Setup(ctx context.Context, r *mux.Router, config *config.Config, c CantExtAPI) *CantabularMetadataExtractorAPI {
	q.Q("SETUP")

	api := &CantabularMetadataExtractorAPI{
		Router:     r,
		CantExtAPI: c,
		Cfg:        config,
	}

	r.HandleFunc("/dataset/{datasetID}/cantabular/{cantdataset}/lang/{lang}", api.getMetadata).Methods("GET")
	return api
}
