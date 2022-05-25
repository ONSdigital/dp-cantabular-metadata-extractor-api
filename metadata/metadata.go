package metadata

import (
	"context"
	"errors"
	"log"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

type Resp struct {
	Dataset struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		Description   string `json:"description"` // Summary on page?
		UnitOfMeasure string `json:"unit_of_measure"`

		Contacts []Contact `json:"contacts"` // slice/Original js was array (!?)

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

type Contact struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

type Dimension struct {
	Name        string `json:"name"` // Title?
	Description string `json:"description"`
	Label       string `json:"label"` // Title?
}

// Client is the client for interacting with the Cantabular API
type Metadata struct {
	Client *cantabular.Client
}

// TODO add lang: cy
func (m *Metadata) GetMetadataTable(datasetID string) (mt *cantabular.MetadataTableQuery, dims []string, err error) {
	req := cantabular.MetadataTableQueryRequest{Variables: []string{datasetID}}

	mt, err = m.Client.MetadataTableQuery(context.Background(), req)
	if err != nil {
		return mt, dims, err
	}

	if len(mt.Service.Tables) == 0 {
		return mt, dims, errors.New("no dims/vars")
	}

	for _, v := range mt.Service.Tables[0].Vars {
		dims = append(dims, string(v))
	}

	return mt, dims, err
}

// TODO add lang: cy
// XXXXXXXXXXXXXXXXXXXXXXXX rename
func (m *Metadata) GetMetadataDataset(cantDataset string, dimensions []string) *cantabular.MetadataDatasetQuery {
	req := cantabular.MetadataDatasetQueryRequest{}
	req.Dataset = cantDataset
	req.Variables = dimensions

	r, err := m.Client.MetadataDatasetQuery(context.Background(), req)

	if err != nil {
		log.Print(err)
	}

	return r
}
