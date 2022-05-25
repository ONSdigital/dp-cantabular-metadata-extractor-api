package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/metadata"
	"github.com/gorilla/mux"

	dphttp "github.com/ONSdigital/dp-net/http"
)

// getMetadata is the main entry point
func (api *CantabularMetadataExtractorAPI) getMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)

	mt, dimensions, _ := api.GetMetadataTable(ctx, params["datasetID"])

	md, err := api.GetMetadataDataset(ctx, params["cantdataset"], dimensions)
	if err != nil {
		log.Print(err) // XXX
	}

	m := cantabular.MetadataQueryResult{TableQueryResult: mt, DatasetQueryResult: md}

	json, _ := json.Marshal(m)
	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) GetMetadataTable(ctx context.Context, cantDataset string) (*cantabular.MetadataTableQuery, []string, error) {
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: api.Cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &metadata.Metadata{Client: cantabularClient}
	return m.GetMetadataTable(cantDataset)

}

func (api *CantabularMetadataExtractorAPI) GetMetadataDataset(ctx context.Context, cantDataset string, dims []string) (*cantabular.MetadataDatasetQuery, error) {
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: api.Cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &metadata.Metadata{Client: cantabularClient}
	return m.GetMetadataDataset(cantDataset, dims)

}

/*
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
*/
