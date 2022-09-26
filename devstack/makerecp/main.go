package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/devstack/makerecp/createrecipe"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var id, host, extapihost, checkdims, alias string
	var check, checkall, autoalias, list bool
	flag.StringVar(&id, "id", "TS009", "specify pre-defined query id")
	flag.StringVar(&host, "host", "http://localhost:28300", "specify extractor-api url")
	flag.StringVar(&extapihost, "extapihost", "http://localhost:8492", "specify extapi url")
	flag.StringVar(&checkdims, "checkdims", "", "check list of dims, eg. \"ltla,sex\" ")
	flag.StringVar(&alias, "alias", "Testing for metadata demo v3", "set alias manually")
	flag.BoolVar(&autoalias, "setalias", false, "set alias/name automatically from metadata server label")
	flag.BoolVar(&check, "check", false, "check specified id use 'id' as well")
	flag.BoolVar(&checkall, "checkall", false, "check all ids known to this program")
	flag.BoolVar(&list, "list", false, "list ids known to this program")
	flag.Parse()

	if checkall {

		for id := range createrecipe.GetMap() {
			fmt.Printf("Testing id=%s ", id)
			cr := createrecipe.New(id, host, extapihost)

			// check against metadata server
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

			// check against cantabular server

			if !cr.OKDimsInDS() {
				log.Fatalf("dims '%#v' not fully present in '%s' dataset", cr.Dimensions, "UR") // XXX
			}

			fmt.Print("OK\n")

		}

		os.Exit(0)

	}

	cr := createrecipe.New(id, host, extapihost)

	if list {
		fmt.Println(strings.Join(cr.ValidIDs, " "))
		os.Exit(0)
	}

	// check against cantabular server
	if checkdims != "" {
		cr.Dimensions = strings.Split(checkdims, ",")
		if !cr.OKDimsInDS() {
			log.Fatalf("dims '%#v' not fully present in '%s' dataset", cr.Dimensions, "UR") // XXX
		} else {
			fmt.Println("dims OK")
		}
		os.Exit(0)
	}

	if !cr.OKDimsInDS() {
		log.Fatalf("dims '%#v' not fully present in '%s' dataset", cr.Dimensions, "UR") // XXX
	}

	if !cr.CheckID() {
		log.Fatalf("'%s' not in valid id list '%#v'", id, cr.ValidIDs)
	}

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

	if autoalias {
		tf, err := cr.GetMetaData()
		if err != nil {
			log.Fatal(err)
		}

		alias = "RTT: " + tf.TableQueryResult.Service.Tables[0].Label + "(" + id + ") " + time.Now().Format("20060102-1504:05")
	}

	r := createrecipe.Recipe{
		Alias:          alias,
		CantabularBlob: "UR", // XXX
		Format:         "cantabular_table",
		ID:             cr.UUID,
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
