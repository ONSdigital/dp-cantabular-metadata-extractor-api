package api_test

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/shurcooL/graphql"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMetadataTable(t *testing.T) {

	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorApi.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorApi.CantExtAPI = &mock.CantExtAPIMock{
			GetMetadataTableFunc: func(datasetID string) (*cantabular.MetadataTableQuery, []string, error) {
				mt := &cantabular.MetadataTableQuery{}
				return mt, []string{"Age", "Sex"}, nil
			},
		}
	})

	Convey("getDimensions method should return correct dimensions", t, func() {
		expected := []string{"Age", "Sex"}
		_, dims, err := cantMetadataExtractorApi.CantExtAPI.GetMetadataTable("Teaching-Dataset")
		if err != nil {
			t.Fail()
		}
		So(dims, ShouldResemble, expected)
	})

	// ... moar
}

func TestGetMetadataDataset(t *testing.T) {

	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorApi.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorApi.CantExtAPI = &mock.CantExtAPIMock{
			GetMetadataDatasetFunc: func(cantDataset string, dimensions []string) (*cantabular.MetadataDatasetQuery, error) {
				md := &cantabular.MetadataDatasetQuery{}
				md.Dataset.Description = graphql.String("This is some summary test...")
				return md, nil
			},
		}
	})

	Convey("getDimensions method should return correct dimensions", t, func() {
		md, err := cantMetadataExtractorApi.CantExtAPI.GetMetadataDataset("Teaching-Dataset", []string{"Age", "Sex"})
		if err != nil {
			t.Fail()
		}
		So(md.Dataset.Description, ShouldResemble, graphql.String("This is some summary test..."))

	})
}

/*
func TestGetVersionDimensions(t *testing.T) {

	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorApi.Cfg, _ = config.Get()

	Convey("Given a mock DatasetAPI client and dataset", t, func() {
		cantMetadataExtractorApi.DatasetAPI = &mock.DatasetAPIMock{
			GetVersionDimensionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
				mockJson := `{"items":[{"id":"","name":"Age","links":{"access_rights":{"href":""},"dataset":{"href":""},"dimensions":{"href":""},"edition":{"href":""},"editions":{"href":""},"latest_version":{"href":""},"versions":{"href":""},"self":{"href":""},"code_list":{"href":"http://localhost:22400/code-lists/Age","id":"Age"},"options":{"href":"http://localhost:22000/datasets/initial-metadata-poc-demo-v1/editions/2021/versions/1/dimensions/Age/options","id":"Age"},"version":{"href":"http://localhost:22000/datasets/initial-metadata-poc-demo-v1/editions/2021/versions/1"},"code":{"href":""},"taxonomy":{"href":""},"job":{"href":""}},"description":"","label":""},{"id":"","name":"Country","links":{"access_rights":{"href":""},"dataset":{"href":""},"dimensions":{"href":""},"edition":{"href":""},"editions":{"href":""},"latest_version":{"href":""},"versions":{"href":""},"self":{"href":""},"code_list":{"href":"http://localhost:22400/code-lists/Country","id":"Country"},"options":{"href":"http://localhost:22000/datasets/initial-metadata-poc-demo-v1/editions/2021/versions/1/dimensions/Country/options","id":"Country"},"version":{"href":"http://localhost:22000/datasets/initial-metadata-poc-demo-v1/editions/2021/versions/1"},"code":{"href":""},"taxonomy":{"href":""},"job":{"href":""}},"description":"","label":"Country"}]}`
				var mockReturn dataset.VersionDimensions
				json.Unmarshal([]byte(mockJson), &mockReturn)

				return mockReturn, nil
			},
		}
		mockDataset := api.Dataset{
			ID:      "test_id",
			Edition: "test_edition",
			Version: "test_version",
		}

		Convey("getDimensions method should return correct dimensions", func() {
			expected := []string{"Age", "Country"}
			actual, err := cantMetadataExtractorApi.GetDimensions(context.Background(), mockDataset)
			if err != nil {
				t.Fail()
			}
			So(actual, ShouldResemble, expected)

		})
	})
	Convey("Given a mock DatasetAPI client with GetVersionDimensions returning an error", t, func() {
		cantMetadataExtractorApi.DatasetAPI = &mock.DatasetAPIMock{
			GetVersionDimensionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (dataset.VersionDimensions, error) {
				var mockReturn dataset.VersionDimensions
				return mockReturn, errors.New("error")
			},
		}
		mockDataset := api.Dataset{
			ID:      "test_id",
			Edition: "test_edition",
			Version: "test_version",
		}

		Convey("getDimensions method should return an error", func() {
			expectedErr := errors.New("failed to get version dimensions: error")
			_, actualErr := cantMetadataExtractorApi.GetDimensions(context.Background(), mockDataset)
			So(actualErr.Error(), ShouldResemble, expectedErr.Error())

		})
	})
}
*/
