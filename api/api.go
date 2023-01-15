package api

import (
	"context"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
)

const (
	CantabularMetadataDatasetReadPermission   string = "cantabularmetadatadataset:read"
)

//go:generate moq -out mock/datasetapi.go -pkg mock . CantMetaAPI

// cantExtAPI
type CantMetaAPI interface {
	MetadataDatasetQuery(ctx context.Context, req cantabular.MetadataDatasetQueryRequest) (*cantabular.MetadataDatasetQuery, error)
	MetadataTableQuery(ctx context.Context, req cantabular.MetadataTableQueryRequest) (*cantabular.MetadataTableQuery, error)
}

type CantabularMetadataExtractorAPI struct {
	Router      *mux.Router
	CantMetaAPI CantMetaAPI
	Cfg         *config.Config
	auth          authorisation.Middleware
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, 
	r *mux.Router, 
	config *config.Config, 
	c CantMetaAPI, 
	auth authorisation.Middleware) *CantabularMetadataExtractorAPI {

	api := &CantabularMetadataExtractorAPI{
		Router:      r,
		CantMetaAPI: c,
		Cfg:         config,
		auth: auth,
	}

	r.HandleFunc("/cantabular-metadata/dataset/{datasetID}/lang/{lang}", auth.Require(CantabularMetadataDatasetReadPermission, api.getMetadata)).Methods("GET")
	return api
}
