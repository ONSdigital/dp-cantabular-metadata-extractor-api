package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

	dimensions, err := api.GetDimensions(ctx, dataset)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	json, _ := json.Marshal(dimensions)
	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) GetDimensions(ctx context.Context, d Dataset) ([]string, error) {
	fullDimensions, err := api.datasetAPI.GetVersionDimensions(ctx, "", api.cfg.ServiceAuthToken, "", d.ID, d.Edition, d.Version)

	if err != nil {
		return nil, fmt.Errorf("failed to get version dimensions: %w", err)
	}

	dimensionsSlice := []string{}

	for _, dimension := range fullDimensions.Items {
		dimensionsSlice = append(dimensionsSlice, dimension.Name)
	}

	return dimensionsSlice, nil
}
