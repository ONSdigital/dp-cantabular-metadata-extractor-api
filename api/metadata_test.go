package api_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api/mock"
	"github.com/ryboe/q"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersionDimensions(t *testing.T) {

	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	ctx := context.Background()
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
			q.Q(ctx)
			actual, err := cantMetadataExtractorApi.GetDimensions(context.Background(), mockDataset)
			if err != nil {
				t.Fail()
			}
			So(actual, ShouldEqual, expected)

		})
	})
}
