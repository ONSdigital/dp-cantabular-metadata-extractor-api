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
				if err != nil {
					t.Error(err)
				}
				return &mt, nil
			},
		}
		Convey("GetMetadataTable method should return correct dimensions", func() { // XXX
			expected := []string{"Region", "Occupation", "Age"}
			_, dims, err := cantMetadataExtractorAPI.GetMetadataTable(ctx, "Teaching-Dataset", "en")
			if err != nil {
				t.Error(err)
			}
			So(dims, ShouldResemble, expected)
		})
	})

	// ... moar
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
			if err != nil {
				t.Fail()
			}
			So(md.Dataset.Description, ShouldResemble, graphql.String("This is some summary test..."))

		})
	})
}

func getMT() (cantabular.MetadataTableQuery, error) {

	j := `{
    "service": {
      "tables": [
        {
          "name": "LC6112EW",
          "label": "Occupation by age",
          "description": "This dataset provides 2011 Census estimates that classify all usual residents in employment the week before the census in England and Wales by occupation and by age. The estimates are as at census day, 27 March 2011.",
          "vars": [
            "Region",
            "Occupation",
            "Age"
          ],
          "meta": {
            "contact": {
              "contact_name": "Census Customer Services",
              "contact_email": "census.customerservices@ons.gov.uk",
              "contact_phone": "01329 444 972",
              "contact_website": "https://www.ons.gov.uk/census/censuscustomerservices"
            },
            "census_releases": [
              {
                "census_release_description": "Example release: labour market, housing and qualifications",
                "census_release_number": "3",
                "release_date": "26/02/2014"
              }
            ],
            "dataset_mnemonic2011": "LC6112EW",
            "dataset_population": "All usual residents",
            "dissemination_source": "Census 2011",
            "geographic_coverage": "England and Wales",
            "geographic_variable_mnemonic": "Region",
            "last_updated": "26/02/2014",
            "keywords": [],
            "publications": [],
            "related_datasets": [
              "LC6107EW"
            ],
            "release_frequency": "",
            "statistical_unit": {
              "statistical_unit": "People",
              "statistical_unit_description": "People living in England and Wales"
            },
            "unique_url": "",
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
