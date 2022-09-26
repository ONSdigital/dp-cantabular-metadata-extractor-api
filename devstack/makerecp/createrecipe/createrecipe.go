package createrecipe

import (
	"crypto/rand"
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

// data tested against "UR" 20220830
// XXX only cantabular dataset ids which originally used "UR" (now mapped to "UR"
// will work)

const HackedDataSetName = "UR"

type CreateRecipe struct {
	ONSDataSetID string
	Dimensions   []string
	Host         string
	ExtApiHost   string
	ValidIDs     []string
	UUID         string
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
		ValidIDs:     validIDs,
		UUID:         uuidV4(),
	}
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
dataset(name:"UR") {
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
		"TS002": {"oa", "legal_partnership_status"},
		"TS005": {"oa", "passports_all_27a"},
		"TS007": {"msoa", "resident_age_101a"},
		"TS008": {"oa", "sex"},
		"TS009": {"ltla", "sex", "resident_age_91a"},
		"TS012": {"ltla", "country_of_birth_60a"},
		"TS013": {"msoa", "passports_all_52a"},
		"TS015": {"oa", "year_arrival_uk"},
		"TS016": {"oa", "residence_length_6b"},
		"TS021": {"oa", "ethnic_group_tb_20b"},
		"TS024": {"ltla", "main_language_detailed"},
		"TS027": {"oa", "national_identity_all"},
		"TS028": {"oa", "national_identity_detailed"},
		"TS029": {"oa", "english_proficiency"},
		"TS030": {"oa", "religion_tb"},
		"TS032": {"oa", "welsh_skills_all"},
		"TS033": {"oa", "welsh_skills_speak"},
		"TS034": {"oa", "welsh_skills_write"},
		"TS035": {"oa", "welsh_skills_read"},
		"TS036": {"oa", "welsh_skills_understand"},
		"TS037": {"oa", "health_in_general"},
		"TS038": {"oa", "disability"},
		"TS039": {"oa", "is_carer"},
		"TS056": {"oa", "alternative_address_indicator"},
		"TS058": {"oa", "workplace_travel_10a"},
		"TS059": {"oa", "hours_per_week_worked"},
		"TS060": {"msoa", "industry_current_88a"},
		"TS061": {"oa", "transport_to_workplace_12a"},
		"TS062": {"oa", "ns_sec_10a"},
		"TS063": {"oa", "occupation_current_10a"},
		"TS064": {"msoa", "occupation_current_105a"},
		"TS065": {"oa", "has_ever_worked"},
		"TS066": {"oa", "economic_activity_status_12a"},
		"TS067": {"oa", "highest_qualification"},
		"TS071": {"msoa", "uk_armed_forces"},
		"TS076": {"ltla", "welsh_skills_speak", "resident_age_86a"},
	}

	for _, v := range m {
		// XXX we need to override to be "ltla" always Fran 20220831
		v[0] = "ltla"
	}

	return m
}

func IsGeo(s string) bool {
	// UR
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

func uuidV4() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return uuid
	}

	// for version 4 (rand) uuid
	// this makes sure that the 13th character is "4"
	b[6] = (b[6] | 0x40) & 0x4F
	// this makes sure that the 17th is "8", "9", "a", or "b"
	b[8] = (b[8] | 0x80) & 0xBF

	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
