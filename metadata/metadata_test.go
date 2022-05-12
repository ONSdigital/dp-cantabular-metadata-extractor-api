package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/mock"
	"github.com/ryboe/q"

	//dphttp "github.com/ONSdigital/dp-net/http"
	. "github.com/smartystreets/goconvey/convey"
)

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

			q.Q(md.Dataset.Meta.Source.Contact)

			Convey("Then the expected metadata information should be returned", func() {
				So(md.Dataset.Meta.Source.Contact.ContactEmail, ShouldEqual, "census.customerservices@ons.gov.uk")
				/*
					So(md.Codebook, ShouldHaveLength, 5)
					So(md.Codebook[0].Name, ShouldEqual, "city")
					So(md.Codebook[1].Labels[0], ShouldEqual, "England")
					So(md.Codebook[2].Labels, ShouldHaveLength, 2)
					So(md.Codebook[3].Codes[2], ShouldEqual, "2")
					So(md.Codebook[4].MapFrom[0].SourceNames[0], ShouldEqual, "siblings")
					So(md.Codebook[4].MapFrom[0].Code[1], ShouldEqual, "1-2")
				*/
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

func TestGetCantabularMetaData(t *testing.T) {
	dims := []string{"Age", "Country"}

	resp := getCantabularMetaData("Teaching-Dataset", dims)

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

	/*
		// debugging
		bs, err := json.Marshal(resp)

		if err != nil {
			t.Fail()
		}

		println(jsonpp(bs))

	*/

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
