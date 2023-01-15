package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authorisationMock "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"
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
		c := &mock.CantMetaAPIMock{}
		cfg, err := config.Get()
		authorisationMiddleware := &authorisationMock.MiddlewareMock{
			RequireFunc: func(permission string, handlerFunc http.HandlerFunc) http.HandlerFunc {
				return handlerFunc
			},
			CloseFunc: func(ctx context.Context) error {
				return nil
			},
		}
		if err != nil {
			t.Fail()
		}

		api := api.Setup(ctx, r, cfg, c, authorisationMiddleware)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/cantabular-metadata/dataset/{datasetID}/lang/{lang}", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
