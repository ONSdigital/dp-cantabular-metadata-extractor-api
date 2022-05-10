package metadata

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"sort"
	"testing"
)

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
