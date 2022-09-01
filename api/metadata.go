package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

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

	// XXX Are all vars in the same (cantabular) dataset?
	cantdataset := string(mt.Service.Tables[0].DatasetName)

	md, err := api.GetMetadataDataset(ctx, cantdataset, dimensions, params["lang"])
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// we always override the geographical code (index 0) to ltla since the
	// "base" in the metadata is "oa" which isn't suitable - Fran 20220831

	// XXX first attempt for MVP but this should be more robust
	// TODO query the MD again for ltla content if ltla not in use (usual case)
	// look for geo variable properly and replace whole var

	mt.Service.Tables[0].Vars[0] = "ltla"

	md.Dataset.Vars[0].Name = "ltla"
	md.Dataset.Vars[0].Description = "As of 2022 there are 309 lower tier local authorities in England, comprising non-metropolitan districts (181), unitary authorities (59), metropolitan districts (36) and London boroughs (33, including City of London). There are 22 lower tier local authorities in Wales, comprising 22 unitary authorities"
	md.Dataset.Vars[0].Label = "Lower Tier Local Authorities"

	md.Dataset.Vars[0].Meta.ONSVariable.VariableTitle = "Lower Tier Local Authorities"
	md.Dataset.Vars[0].Meta.ONSVariable.GeographicTheme = "Administrative"
	md.Dataset.Vars[0].Meta.ONSVariable.VariableMnemonic = "ltla"
	md.Dataset.Vars[0].Meta.ONSVariable.VariableTitle = "Lower Tier Local Authorities"

	m := cantabular.MetadataQueryResult{TableQueryResult: mt, DatasetQueryResult: md}

	json, err := json.Marshal(m)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		log.Error(ctx, err.Error(), err)
	}
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
