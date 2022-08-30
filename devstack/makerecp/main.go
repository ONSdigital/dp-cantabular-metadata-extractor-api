package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/devstack/makerecp/createrecipe"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var id, host string
	var check, setalias, list bool
	flag.StringVar(&id, "id", "TS009", "specify pre-defined query id")
	flag.StringVar(&host, "host", "http://localhost:28300", "specify extractor-api url")
	flag.BoolVar(&setalias, "setalias", false, "set alias/name automatically from metadata server label")
	flag.BoolVar(&check, "check", false, "")
	flag.BoolVar(&list, "list", false, "list ids known to this program")
	flag.Parse()

	cr := createrecipe.New(id, host)

	if list {
		fmt.Println(strings.Join(cr.ValidIDs, " "))
		os.Exit(0)
	}

	if !cr.CheckID() {
		log.Fatalf("'%s' not in valid id list '%#v'", id, cr.ValidIDs)
	}

	alias := "Testing for metadata demo v3"

	if check {
		tf, err := cr.GetMetaData()
		if err != nil {
			log.Fatal(err)
		}

		datasetName := tf.TableQueryResult.Service.Tables[0].DatasetName

		if datasetName != createrecipe.HackedDataSetName {
			log.Fatalf("wrong dataset name '%s' need '%s'", datasetName, createrecipe.HackedDataSetName)
		}

		_, ourVars := createrecipe.SplitVars(cr.Dimensions)

		_, mdVars := createrecipe.SplitVars(tf.TableQueryResult.Service.Tables[0].Vars)

		if !reflect.DeepEqual(mdVars, ourVars) {
			log.Fatalf("expected vars '%#v' don't match metadata-server vars '%#v'", ourVars, mdVars)
		}

		fmt.Println("cantabular dataset names match & non geographical dimensions match OK")
		os.Exit(0)
	}

	if setalias {
		tf, err := cr.GetMetaData()
		if err != nil {
			log.Fatal(err)
		}

		alias = tf.TableQueryResult.Service.Tables[0].Label
	}

	r := createrecipe.Recipe{
		Alias:          alias,
		CantabularBlob: "dp_synth_dataset",
		Format:         "cantabular_table",
		ID:             "la2e031b-3064-427d-8fed-4b35c99bf1a0", // XXX
	}

	r.OutputInstances = []createrecipe.OutputInstance{{ // XXX
		CodeLists: cr.GetCodeLists(),
		DatasetID: id,
		Editions:  []string{"2021"},
		Title:     alias,
	}}

	u, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(u))

}
