package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	datasetAPI "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
)

type Dataset struct {
	ID      string
	Edition string
	Version string
}

func (api *CantabularMetadataExtractorAPI) getMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	dataset := Dataset{ID: params["datasetID"], Edition: params["editionID"], Version: params["versionID"]}

	dimensions, err := api.getDimensions(ctx, dataset)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	json, _ := json.Marshal(dimensions)
	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) getDimensions(ctx context.Context, d Dataset) (*datasetAPI.VersionDimensions, error) {

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("error getting configuration: %w", err)
	}
	dimensions, err := api.datasetAPI.GetVersionDimensions(ctx, "", cfg.ServiceAuthToken, "", d.ID, d.Edition, d.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version dimensions: %w", err)
	}

	return &dimensions, nil
}
