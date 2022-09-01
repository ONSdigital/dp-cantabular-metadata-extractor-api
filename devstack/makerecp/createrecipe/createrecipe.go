package createrecipe

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Recipe struct {
	Alias           string `json:"alias"`
	CantabularBlob  string `json:"cantabular_blob"`
	Format          string `json:"format"`
	ID              string `json:"id"`
	OutputInstances `json:"output_instances"`
}

type OutputInstance struct {
	CodeLists `json:"code_lists"`
	DatasetID string   `json:"dataset_id"`
	Editions  []string `json:"editions"`
	Title     string   `json:"title"`
}

type OutputInstances []OutputInstance
type CodeLists []CodeList

type CodeList struct {
	Href                         string `json:"href"`
	ID                           string `json:"id"`
	IsHierarchy                  bool   `json:"is_hierarchy"`
	Name                         string `json:"name"`
	IsCantabularGeography        bool   `json:"is_cantabular_geography"`
	IsCantabularDefaultGeography bool   `json:"is_cantabular_default_geography"`
}

// enough of the table response we use
type TableFrag struct {
	TableQueryResult struct {
		Service struct {
			Tables []struct {
				Name        string   `json:"name"`
				DatasetName string   `json:"dataset_name"`
				Label       string   `json:"label"`
				Vars        []string `json:"vars"`
				Meta        struct {
					AlternateGeographicVariables []string `json:"alternate_geographic_variables"`
				} `json:"meta"`
			} `json:"tables"`
		} `json:"service"`
	} `json:"table_query_result"`
}

// data tested against "dp_synth_dataset" 20220830
// XXX only cantabular dataset ids which originally used "UR" (now mapped to "dp_synth_dataset"
// will work)

const HackedDataSetName = "dp_synth_dataset"

type CreateRecipe struct {
	ONSDataSetID string
	Dimensions   []string
	Host         string
	ExtApiHost   string
	ValidIDs     []string
}

func New(id, host, extApiHost string) *CreateRecipe {
	var validIDs []string
	for k := range GetMap() {
		validIDs = append(validIDs, k)
	}
	sort.Strings(validIDs)
	return &CreateRecipe{
		ONSDataSetID: id,
		Host:         host,
		ExtApiHost:   extApiHost,
		Dimensions:   GetMap()[id],
		ValidIDs:     validIDs}
}

func (cr *CreateRecipe) GetMetaData() (TableFrag, error) {
	var tf TableFrag
	resp, err := http.Get(fmt.Sprintf("%s/cantabular-metadata/dataset/%s/lang/en", cr.Host, cr.ONSDataSetID))
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf, err
	}

	if err := json.Unmarshal(body, &tf); err != nil {
		return tf, err
	}

	return tf, err
}

func (cr *CreateRecipe) GetCodeLists() (cls CodeLists) {

	for _, v := range cr.Dimensions {
		cl := CodeList{
			Href:        fmt.Sprintf("http://localhost:22400/code-lists/%s", v),
			ID:          v,
			Name:        v,
			IsHierarchy: false,
		}

		if IsGeo(v) {
			cl.IsCantabularDefaultGeography = true
			cl.IsCantabularGeography = true
		}
		cls = append(cls, cl)
	}

	return cls
}

func (cr *CreateRecipe) CheckID() bool {
	return InSlice(cr.ONSDataSetID, cr.ValidIDs)
}

func (cr *CreateRecipe) OKDimsInDS() bool {

	query := `variables={}&query={
dataset(name:"dp_synth_dataset") {
	variables(names:["%s"]) {
	  edges {
		node {
		  name 
		  description
		  label
		}
	  }
	}
   }
 } `

	payload := url.PathEscape(fmt.Sprintf(query, strings.Join(cr.Dimensions, "\",\"")))

	resp, err := http.Get(fmt.Sprintf("%s/graphql?%s", cr.ExtApiHost, payload))
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return !strings.Contains(string(body), "does not exist")

}

// GetMap returns our definitions as the Golden Source of Truth (maybe)
func GetMap() map[string][]string {
	// these are the real values and we overside the geo ones, eg. oa
	m := map[string][]string{
		"AP001": {"oa", "sex"},                                      // WORKS
		"AP002": {"oa", "resident_age_11a"},                         // broken in MD release
		"AP003": {"oa", "legal_partnership_status_6a"},              // broken in MD release
		"AP011": {"oa", "main_language_11a"},                        // broken in MD release
		"AP025": {"oa", "industry_current_9a"},                      // broken in MD release
		"AP026": {"oa", "occupation_current_10a"},                   // broken in MD release
		"AP027": {"ltla", "transport_to_work"},                      // broken dataset lacks dims
		"RM014": {"ltla", "workplace_travel_5a", "resident_age_6a"}, // broken in MD release
		"TS009": {"ltla", "sex"},                                    // WORKS
	}

	for _, v := range m {
		// XXX we need to override to be "ltla" always Fran 20220831
		v[0] = "ltla"
	}

	return m
}

func IsGeo(s string) bool {
	// dp_synth_dataset
	isGeo := map[string]bool{
		"ctry": true, // England & Wales + Scotland (2022) and NI (2021)
		"lsoa": true, // Lower Layer Super Output Areas
		"ltla": true, // Lower Tier Local Authorities ~330
		"msoa": true, // Middle Layer Super Output Areas
		"nat":  true, // England & Wales
		"oa":   true, // Output Areas ~180K
		"rgn":  true, // 9 UK regions (+ Wales?)
		"utla": true, // Upper Tier Local Authorities
	}

	return isGeo[s]

}
func SplitVars(totalVars []string) (geoVar string, vars []string) {
	for _, v := range totalVars {
		if IsGeo(v) {
			// assume only one geo var at this point
			geoVar = v
		} else {
			vars = append(vars, v)
		}
	}

	return geoVar, vars
}

func InSlice(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}

	return false
}
