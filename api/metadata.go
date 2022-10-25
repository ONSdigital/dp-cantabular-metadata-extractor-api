package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/shurcooL/graphql"
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

	OverrideMetadataTable(dimensions, mt)

	// XXX Are all vars in the same (cantabular) dataset?
	cantdataset := string(mt.Service.Tables[0].DatasetName)

	md, err := api.GetMetadataDataset(ctx, cantdataset, dimensions, params["lang"])
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

// OverrideMetadataTable modifies the dimensions and results of the MetadataTableQuery
// to always use "ltla".  This is the geocode used in the recipe and we need to ensure
// the result from the metadata server matches the recipe.  This ensures also the
// following GetMetadataDataset uses "ltla".
func OverrideMetadataTable(dims []string, mt *cantabular.MetadataTableQuery) {
	geoCodeOverride := "ltla" //  Fran 20220831
	validGeo := []string{"ctry", "lsoa", "ltla", "msoa", "nat", "oa", "rgn", "utla"}

	for i, v := range dims {
		if inSlice(v, validGeo) {
			dims[i] = geoCodeOverride

			break
		}
	}

outer:
	for i, v := range mt.Service.Tables {
		for j, c := range v.Vars {
			if inSlice(string(c), validGeo) {
				mt.Service.Tables[i].Vars[j] = graphql.String(geoCodeOverride)

				break outer
			}
		}
	}
}

func inSlice(x string, xs []string) bool {
	for _, v := range xs {
		if x == v {
			return true
		}
	}

	return false
}
