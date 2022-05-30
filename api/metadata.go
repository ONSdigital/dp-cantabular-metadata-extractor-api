package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Temporary Hack (TM) to convert CamelCase to snake_case
// TODO use proper JSON structs with 2021 metadata - remove this!

var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

type conventionalMarshaller struct {
	Value interface{}
}

func (c conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)

	return converted, err
}

// getMetadata is the main entry point
func (api *CantabularMetadataExtractorAPI) getMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)

	if params["lang"] == "" {
		params["lang"] = "en"
	}

	mt, dimensions, err := api.GetMetadataTable(ctx, params["datasetID"], params["lang"])
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	md, err := api.GetMetadataDataset(ctx, params["cantdataset"], dimensions, params["lang"])
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := cantabular.MetadataQueryResult{TableQueryResult: mt, DatasetQueryResult: md}

	// REMOVE IN PHASE 2
	json, err := json.MarshalIndent(conventionalMarshaller{m}, "", "  ")
	//json, err := json.Marshal(m)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (api *CantabularMetadataExtractorAPI) GetMetadataTable(ctx context.Context, cantDataset string, lang string) (*cantabular.MetadataTableQuery, []string, error) {

	req := cantabular.MetadataTableQueryRequest{Variables: []string{cantDataset}, Lang: lang}

	var dims []string
	mt, err := api.CantExtAPI.MetadataTableQuery(context.Background(), req)
	if err != nil {
		return mt, dims, err
	}

	if len(mt.Service.Tables) == 0 {
		return mt, dims, errors.New("no dims/vars") // XXX
	}

	for _, v := range mt.Service.Tables[0].Vars {
		dims = append(dims, string(v))
	}

	return mt, dims, err
}

func (api *CantabularMetadataExtractorAPI) GetMetadataDataset(ctx context.Context, cantDataset string, dimensions []string, lang string) (*cantabular.MetadataDatasetQuery, error) {

	req := cantabular.MetadataDatasetQueryRequest{}
	req.Dataset = cantDataset
	req.Variables = dimensions
	req.Lang = lang

	md, err := api.CantExtAPI.MetadataDatasetQuery(ctx, req)

	return md, err

}
