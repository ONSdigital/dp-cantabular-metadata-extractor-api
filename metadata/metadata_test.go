package metadata

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ryboe/q"

	//  dphttp "github.com/ONSdigital/dp-net/http"
	. "github.com/smartystreets/goconvey/convey"
)

var intFlag = flag.Bool("int", false, "perform int tests")

// this probably belongs in dp-api-clients-go but is here as a stopgap
func TestMockGetCantabularMetaDataHappy(t *testing.T) {

	Convey("Given a correct response from the Metadata Server", t, func() {
		testCtx := context.Background()

		mockGQLClient := &mock.GraphQLClientMock{QueryFunc: func(ctx context.Context, query interface{}, vars map[string]interface{}) error {
			md := query.(*cantabular.MetadataDatasetQuery)
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
			req := cantabular.MetadataDatasetQueryRequest{}
			md, err := cantabularClient.MetadataDatasetQuery(testCtx, req)
			So(err, ShouldBeNil)

			Convey("Then the expected metadata information should be returned", func() {
				So(md.Dataset.Meta.Source.Contact.ContactEmail, ShouldEqual, "census.customerservices@ons.gov.uk")
			})
		})
	})
}

// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

type conventionalMarshaller struct {
	Value interface{}
}

func (c conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)

	return converted, err
}

func TestMetadataQueryResult(t *testing.T) {
	// int - [ ]
	cfg, _ := config.Get()
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &Metadata{Client: cantabularClient}

	mt, dims, err := m.GetMetadataTable("LC1117EW")
	if err != nil {
		t.Error(err)
	}

	cm, err := m.GetMetadataDataset("Teaching-Dataset", dims) // XXXXXXXXXXXXXXXXXXXXXXX

	if err != nil {
		t.Error(err)
	}

	s := cantabular.MetadataQueryResult{DatasetQueryResult: cm, TableQueryResult: mt}

	s.TableQueryResult.Service.Tables[0].Meta.Keywords = nil
	s.TableQueryResult.Service.Tables[0].Meta.Publications = nil
	s.TableQueryResult.Service.Tables[0].Meta.RelatedDatasets = nil
	s.DatasetQueryResult.Dataset.Variables.Edges[0].Node.Meta.ONSVariable.Keywords = nil

	//serialize(s)

	expected, err := deserialize()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("not equal")

	}

	//	encoded, _ := json.MarshalIndent(conventionalMarshaller{expected}, "", "  ")
	//	fmt.Println(string(encoded))

	//encoded, _ := json.MarshalIndent(conventionalMarshaller{s}, "", "  ")
	//fmt.Println(string(encoded))

	/*

		q.Q(s)
		bs, err := json.Marshal(s)

		println(string(bs))

		if err != nil {
			t.Fail()
		}

		println(jsonpp(bs))

	*/
}

func TestIntGetCantabularMetaData2(t *testing.T) {

	// INT
	cfg, _ := config.Get()
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &Metadata{Client: cantabularClient}

	mt, dims, err := m.GetMetadataTable("LC1117EW")
	if err != nil {
		t.Error(err)
	}

	q.Q(dims) // XXXXXXXXXX

	// debugging
	bs, err := json.Marshal(mt)

	if err != nil {
		t.Fail()
	}

	println(jsonpp(bs))
}

func TestIntGetMetadataDataset(t *testing.T) {

	if !*intFlag {
		t.Skip("not doing int tests")
	}

	cfg, _ := config.Get()
	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: cfg.CantabularExtURL}, dphttp.NewClient(), nil)

	m := &Metadata{Client: cantabularClient}

	dims := []string{"Age", "Country"}
	resp, err := m.GetMetadataDataset("Teaching-Dataset", dims) // XXXXXXXXXXXXXXXXXXXXXXX
	if err != nil {
		t.Fail()
	}

	if resp.Dataset.Meta.Source.Contact.ContactEmail != "census.customerservices@ons.gov.uk" {
		t.Fail()

	}

	var respDims []string
	for _, v := range resp.Dataset.Variables.Edges {
		respDims = append(respDims, string(v.Node.Name))
	}

	sort.Strings(respDims)
	sort.Strings(dims)

	if !reflect.DeepEqual(dims, respDims) {
		t.Error("didn't get the same dims back as we sent!")
	}
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
		md := query.(*cantabular.MetadataDatasetQuery)
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
	resp, err := m.GetMetadataDataset("Teaching-Dataset", dims)

	if err != nil {
		t.Fail()
	}

	if resp.Dataset.Meta.Source.Contact.ContactEmail != "census.customerservices@ons.gov.uk" {
		t.Fail()
	}

	// TODO more coverage...
}

func serialize(s cantabular.MetadataQueryResult) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(s)
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile("serialized.gob", b.Bytes(), 0644); err != nil {
		log.Print(err)
	}
}

func deserialize() (cantabular.MetadataQueryResult, error) {
	s := cantabular.MetadataQueryResult{}
	b, err := ioutil.ReadFile("serialized.gob")
	if err != nil {
		return s, err
	}
	d := gob.NewDecoder(bytes.NewReader(b))
	// Decoding the serialised data
	err = d.Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
