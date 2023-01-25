package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/shurcooL/graphql"
)

var (
	geoCodeOverride   = "ltla"                                                               // Fran 20220831
	validGeo          = []string{"ctry", "lsoa", "ltla", "msoa", "nat", "oa", "rgn", "utla"} // allowlist of codes
	errNotOneGeocode  = errors.New("invalid data - expected exactly one geocode")
	errUnexpectedResp = errors.New("unexpected JSON response")
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
		err = fmt.Errorf("%s : %w", "api.GetMetadataTable", err)
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := OverrideMetadataTable(dimensions, mt); err != nil {
		err = fmt.Errorf("%s : %w", "OverrideMetadataTable", err)
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(mt.Service.Tables) == 0 {
		err := fmt.Errorf("%s : %w", "mt.Service.Tables", errUnexpectedResp)
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cantdataset := string(mt.Service.Tables[0].DatasetName)

	md, err := api.CantMetaAPI.MetadataDatasetQuery(ctx, cantabular.MetadataDatasetQueryRequest{
		Dataset:   cantdataset,
		Variables: dimensions,
		Lang:      params["lang"],
	})
	if err != nil {
		err = fmt.Errorf("%s : %w", "api.CantMetaAPI.MetadataDatasetQuery", err)
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
	mt, err := api.CantMetaAPI.MetadataTableQuery(context.Background(), req)
	if err != nil {
		return mt, dims, err
	}

	if len(mt.Service.Tables) == 0 {

		return mt, dims, fmt.Errorf("%s : %w", "mt.Service.Tables", errUnexpectedResp)
	}

	if len(mt.Service.Tables[0].Vars) == 0 {

		return mt, dims, fmt.Errorf("%s : %w", "mt.Service.Tables.Vars", errUnexpectedResp)
	}

	for _, v := range mt.Service.Tables[0].Vars {
		dims = append(dims, string(v))
	}

	return mt, dims, err
}

// OverrideMetadataTable modifies the dimensions and results of the MetadataTableQuery
// to always use "ltla".  This is the geocode used in the recipe and we need to ensure
// the result from the metadata server matches the recipe.  This ensures also the
// following GetMetadataDataset uses "ltla".
func OverrideMetadataTable(dims []string, mt *cantabular.MetadataTableQuery) error {
	substituted := 0
	for i, v := range dims {
		if inSlice(v, validGeo) {
			dims[i] = geoCodeOverride
			substituted++
		}
	}

	if substituted != 1 {
		return fmt.Errorf("dimensions : %w", errNotOneGeocode)
	}

	substituted = 0
	for i, v := range mt.Service.Tables {
		for j, c := range v.Vars {
			if inSlice(string(c), validGeo) {
				mt.Service.Tables[i].Vars[j] = graphql.String(geoCodeOverride)
				substituted++
			}
		}
	}

	if substituted != 1 {
		return fmt.Errorf("service tables : %w", errNotOneGeocode)
	}

	return nil
}

func inSlice(x string, xs []string) bool {
	for _, v := range xs {
		if x == v {
			return true
		}
	}

	return false
}
