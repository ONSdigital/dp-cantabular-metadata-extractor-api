package api

import (
	"context"
	"encoding/json"
	"net/http"

	datasetAPI "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
)

func (api *CantabularMetadataExtractorAPI) getMetadata(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	
	dimensions, err := api.getDimensions(ctx)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	json, _ := json.Marshal(dimensions)
	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) getDimensions(ctx context.Context) (*datasetAPI.VersionDimensions, error) {

	dimensions, err := api.datasetAPI.GetVersionDimensions(ctx, "", "", "", "", "", "")
	if err != nil {
		return nil, err
	}

	return &dimensions, nil
}

