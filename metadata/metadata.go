package metadata

import (
	"context"
	"log"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	dphttp "github.com/ONSdigital/dp-net/http"
)

type Resp struct {
	Dataset struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		Description   string `json:"description"` // Summary on page?
		UnitOfMeasure string `json:"unit_of_measure"`

		Contact struct { // slice? Original js was array (!?)
			Name      string `json:"name"`
			Telephone string `json:"telephone"`
			Email     string `json:"email"`
		} `json:"contact"`

		Keywords          []string `json:"keywords"`
		License           string   `json:"license"`
		NationalStatistic bool     `json:"national_statistic"`

		Qmi struct {
			Description string `json:"description"`
			Href        string `json:"href"`
			Title       string `json:"title"`
		} `json:"qmi"`
	} `json:"dataset"`

	Version struct {
		ID           string      `json:"id"`
		Version      int         `json:"version"`
		Edition      string      `json:"edition"`
		CollectionID string      `json:"collection_id"`
		Dimensions   []Dimension `json:"dimensions"`
		ReleaseDate  string      `json:"release_date"`
		// NextReleaseDate?
		// ReleaseFrequency?
	} `json:"version"`
}

type Dimension struct {
	Name        string `json:"name"` // Title?
	Description string `json:"description"`
	Label       string `json:"label"` // Title?
}

// TODO add lang: cy
func getCantabularMetaData(cantDataset string, dimensions []string) (resp Resp) {
	cfg, _ := config.Get()

	cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: cfg.CantabularExtURL}, dphttp.NewClient(), nil)
	//cantabularClient := cantabular.NewClient(cantabular.Config{ExtApiHost: "http://127.0.0.1:9090"}, dphttp.NewClient(), nil)
	req := cantabular.MetadataQueryRequest{}
	req.Dataset = cantDataset
	req.Variables = dimensions

	r, err := cantabularClient.MetadataQuery(context.Background(), req)

	if err != nil {
		log.Print(err)
	}

	resp.Dataset.Title = string(r.Dataset.Label)             // ???
	resp.Dataset.Description = string(r.Dataset.Description) // summary?
	resp.Dataset.Contact.Name = string(r.Dataset.Meta.Source.Contact.ContactName)
	resp.Dataset.Contact.Email = string(r.Dataset.Meta.Source.Contact.ContactEmail)
	resp.Dataset.Contact.Telephone = string(r.Dataset.Meta.Source.Contact.ContactPhone)
	resp.Dataset.License = string(r.Dataset.Meta.Source.Licence)
	resp.Dataset.Qmi.Href = string(r.Dataset.Meta.Source.MethodologyLink)

	if string(r.Dataset.Meta.Source.NationalStatisticCertified) == "Y" {
		resp.Dataset.NationalStatistic = true
	}

	for _, edge := range r.Dataset.Variables.Edges {

		resp.Version.Dimensions = append(resp.Version.Dimensions, Dimension{Name: string(edge.Node.Name), Description: string(edge.Node.Meta.ONSVariable.VariableDescription)})

		resp.Dataset.UnitOfMeasure = string(edge.Node.Meta.ONSVariable.StatisticalUnit.StatisticalUnit)

		for _, kw := range edge.Node.Meta.ONSVariable.Keywords {
			resp.Dataset.Keywords = append(resp.Dataset.Keywords, string(kw))
		}
	}

	return resp
}
