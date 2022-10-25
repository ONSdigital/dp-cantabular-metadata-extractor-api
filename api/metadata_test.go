package api_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/shurcooL/graphql"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMetadataTable(t *testing.T) {
	ctx := context.Background()

	cantMetadataExtractorAPI := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorAPI.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorAPI.CantExtAPI = &mock.CantExtAPIMock{

			MetadataTableQueryFunc: func(ctx context.Context, req cantabular.MetadataTableQueryRequest) (*cantabular.MetadataTableQuery, error) {
				mt, err := getMT()
				So(err, ShouldBeNil)
				return &mt, nil
			},
		}
		Convey("GetMetadataTable method should return correct dimensions", func() { // XXX
			expected := []string{"oa", "sex"}
			_, dims, err := cantMetadataExtractorAPI.GetMetadataTable(ctx, "UR", "en")
			So(err, ShouldBeNil)
			So(dims, ShouldResemble, expected)
		})
	})

}

func TestGetMetadataDataset(t *testing.T) {

	ctx := context.Background()
	cantMetadataExtractorAPI := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorAPI.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorAPI.CantExtAPI = &mock.CantExtAPIMock{

			MetadataDatasetQueryFunc: func(ctx context.Context, req cantabular.MetadataDatasetQueryRequest) (*cantabular.MetadataDatasetQuery, error) {
				md := &cantabular.MetadataDatasetQuery{}
				md.Dataset.Description = graphql.String("This is some summary test...")
				return md, nil
			},
		}

		Convey("getDimensions method should return correct dimensions", func() { // XXX
			md, err := cantMetadataExtractorAPI.GetMetadataDataset(ctx, "Teaching-Dataset", []string{"Age", "Sex"}, "en")
			So(err, ShouldBeNil)
			So(md.Dataset.Description, ShouldResemble, graphql.String("This is some summary test..."))

		})
	})
}

func TestOverrideMetadataTable(t *testing.T) {
	Convey("Given setup of dimensions and a MetadataTableQuery", t, func() {
		m, err := getMT()
		So(err, ShouldBeNil)

		mt := &m
		dims := []string{"oa", "sex"}

		Convey("When we call api.OverrideMetadataDataset with dimensions and a MetadataTableQuery ", func() {
			api.OverrideMetadataTable(dims, mt)

			Convey("Then we get the correct ltla overrides for the dimensions and the MetadataTableQuery", func() {
				So(dims[0], ShouldEqual, "ltla")
				So(mt.Service.Tables[0].Vars[0], ShouldEqual, "ltla")
			})
		})
	})

}

func getMT() (cantabular.MetadataTableQuery, error) {
	j := `{
    "service": {
      "tables": [
        {
          "name": "TS008",
          "dataset_name": "UR",
          "label": "Sex",
          "description": "This dataset provides Census 2021 estimates that classify usual residents in England and Wales by sex. The estimates are as at census day, 21 March 2021.",
          "vars": [
            "oa",
            "sex"
          ],
          "meta": {
            "alternate_geographic_variables": [
              "ctry",
              "lsoa",
              "ltla",
              "msoa",
              "nat",
              "rgn",
              "utla"
            ],
            "contact": {
              "contact_name": "",
              "contact_email": "",
              "contact_phone": "",
              "contact_website": ""
            },
            "census_releases": [],
            "dataset_mnemonic2011": "",
            "dataset_population": "All usual residents",
            "geographic_coverage": "England and Wales",
            "last_updated": "",
            "observation_type": {
              "observation_type_description": "Count",
              "observation_type_label": "Count",
              "decimal_places": "0",
              "prefix": "",
              "suffix": "",
              "fill_trailing_spaces": "Y",
              "negative_sign": "",
              "observation_type_code": "CT"
            },
            "publications": [],
            "related_datasets": [],
            "statistical_unit": {
              "statistical_unit": "Person",
              "statistical_unit_description": "Person"
            },
            "version": "1"
          }
        }
      ]
      }
    }`

	mtq := cantabular.MetadataTableQuery{}
	if err := json.Unmarshal([]byte(j), &mtq); err != nil {
		return mtq, err
	}

	return mtq, nil
}
