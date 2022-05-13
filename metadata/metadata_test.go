package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	dphttp "github.com/ONSdigital/dp-net/http"

	//  dphttp "github.com/ONSdigital/dp-net/http"
	. "github.com/smartystreets/goconvey/convey"
)

var intFlag = flag.Bool("int", false, "perform int tests")

// this probably belongs in dp-api-clients-go but is here as a stopgap
func TestMockGetCantabularMetaDataHappy(t *testing.T) {

	Convey("Given a correct response from the Metadata Server", t, func() {
		testCtx := context.Background()

		mockGQLClient := &mock.GraphQLClientMock{QueryFunc: func(ctx context.Context, query interface{}, vars map[string]interface{}) error {
			md := query.(*cantabular.MetadataQuery)
			md.Dataset.Meta.Source.Contact.ContactEmail = "census.customerservices@ons.gov.uk"
			return nil
		},
		}
		cantabularClient := cantabular.NewClient(
			cantabular.Config{
				ExtApiHost: "cantabular.ext.host",
			},
			nil,
			mockGQLClient,
		)

		Convey("When the MetadataQuery method is called", func() {
			req := cantabular.MetadataQueryRequest{}
			md, err := cantabularClient.MetadataQuery(testCtx, req)
			So(err, ShouldBeNil)

			Convey("Then the expected metadata information should be returned", func() {
				So(md.Dataset.Meta.Source.Contact.ContactEmail, ShouldEqual, "census.customerservices@ons.gov.uk")
			})
		})
	})
}

func testMetadataResponse() ([]byte, error) {
	b, err := ioutil.ReadFile("metadata_test.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	return b, nil
}

func Response(body []byte, statusCode int) *http.Response {
	reader := bytes.NewBuffer(body)
	readCloser := ioutil.NopCloser(reader)

	return &http.Response{
		StatusCode: statusCode,
		Body:       readCloser,
	}
}

func TestMockGetCantabularMetaData(t *testing.T) {

	if !*intFlag {
		t.Skip("not doing int tests")
	}

	cfg, _ := config.Get()
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &Metadata{Client: cantabularClient}

	dims := []string{"Age", "Country"}
	resp := m.GetMetaData("Teaching-Dataset", dims)

	if resp.Dataset.Contact.Email != "census.customerservices@ons.gov.uk" {
		t.Fail()
	}

	var respDims []string

	for _, v := range resp.Version.Dimensions {
		respDims = append(respDims, v.Name)
	}
	sort.Strings(respDims)
	sort.Strings(dims)

	if !reflect.DeepEqual(dims, respDims) {
		t.Error("didn't get the same dims back as we sent!")
	}

	// debugging
	bs, err := json.Marshal(resp)

	if err != nil {
		t.Fail()
	}

	println(jsonpp(bs))

}

func jsonpp(b []byte) (s string) {
	var out bytes.Buffer
	if err := json.Indent(&out, b, " ", " "); err != nil {
		log.Print(err)
	} else {
		s = out.String()
	}
	return s
}

func TestGetCantabularMetaData(t *testing.T) {

	mockGQLClient := &mock.GraphQLClientMock{QueryFunc: func(ctx context.Context, query interface{}, vars map[string]interface{}) error {
		md := query.(*cantabular.MetadataQuery)
		md.Dataset.Meta.Source.Contact.ContactEmail = "census.customerservices@ons.gov.uk"
		return nil
	},
	}
	cantabularClient := cantabular.NewClient(
		cantabular.Config{
			ExtApiHost: "cantabular.ext.host",
		},
		nil,
		mockGQLClient,
	)

	dims := []string{"Age", "Country"}

	m := &Metadata{Client: cantabularClient}
	resp := m.GetMetaData("Teaching-Dataset", dims)

	if resp.Dataset.Contact.Email != "census.customerservices@ons.gov.uk" {
		t.Fail()
	}

	// TODO more coverage...
}
