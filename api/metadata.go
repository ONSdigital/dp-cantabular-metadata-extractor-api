package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/metadata"
	dphttp "github.com/ONSdigital/dp-net/http"
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

	// TODO err
	// XXX Cantabular dataset name hardcoded
	resp := api.getCantMeta(ctx, "Teaching-Dataset", dimensions)

	json, _ := json.Marshal(resp)
	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) GetDimensions(ctx context.Context, d Dataset) ([]string, error) {
	fullDimensions, err := api.DatasetAPI.GetVersionDimensions(ctx, "", api.Cfg.ServiceAuthToken, "", d.ID, d.Edition, d.Version)

	if err != nil {
		return nil, fmt.Errorf("failed to get version dimensions: %w", err)
	}

	dimensionsSlice := []string{}

	for _, dimension := range fullDimensions.Items {
		dimensionsSlice = append(dimensionsSlice, dimension.Name)
	}

	return dimensionsSlice, nil
}

func (api *CantabularMetadataExtractorAPI) getCantMeta(ctx context.Context, cantDataset string, dims []string) metadata.Resp {
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: api.cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	// TODO return error
	m := &metadata.Metadata{Client: cantabularClient}
	resp := m.GetMetaData(cantDataset, dims)

	return resp
}
