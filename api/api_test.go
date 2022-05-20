package api_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		d := &mock.CantExtAPIMock{}
		c, err := config.Get()
		if err != nil {
			t.Fail()
		}

		api := api.Setup(ctx, r, c, d)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/dataset/{datasetID}/cantabular/{cantdataset}/lang/{lang}", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
