package api_test

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/shurcooL/graphql"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMetadataTable(t *testing.T) {
	ctx := context.Background()

	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorApi.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorApi.CantExtAPI = &mock.CantExtAPIMock{

			MetadataTableQueryFunc: func(ctx context.Context, req cantabular.MetadataTableQueryRequest) (*cantabular.MetadataTableQuery, error) {
				mt := getMT()
				return &mt, nil
			},
		}
		Convey("GetMetadataTable method should return correct dimensions", func() { // XXX
			expected := []string{"Region", "Ethnic Group", "Sex", "Age"}
			_, dims, err := cantMetadataExtractorApi.GetMetadataTable(ctx, "Teaching-Dataset", "en")
			if err != nil {
				t.Error(err)
			}
			So(dims, ShouldResemble, expected)
		})
	})

	// ... moar
}

func TestGetMetadataDataset(t *testing.T) {

	ctx := context.Background()
	cantMetadataExtractorApi := &api.CantabularMetadataExtractorAPI{}
	cantMetadataExtractorApi.Cfg, _ = config.Get()

	Convey("Given a mock CantExtAPI client and datasetID/cantabular table", t, func() {
		cantMetadataExtractorApi.CantExtAPI = &mock.CantExtAPIMock{

			MetadataDatasetQueryFunc: func(ctx context.Context, req cantabular.MetadataDatasetQueryRequest) (*cantabular.MetadataDatasetQuery, error) {
				md := &cantabular.MetadataDatasetQuery{}
				md.Dataset.Description = graphql.String("This is some summary test...")
				return md, nil
			},
		}

		Convey("getDimensions method should return correct dimensions", func() { // XXX
			md, err := cantMetadataExtractorApi.GetMetadataDataset(ctx, "Teaching-Dataset", []string{"Age", "Sex"}, "en")
			if err != nil {
				t.Fail()
			}
			So(md.Dataset.Description, ShouldResemble, graphql.String("This is some summary test..."))

		})
	})
}

func getMT() cantabular.MetadataTableQuery {

	return cantabular.MetadataTableQuery{Service: struct {
		Tables []struct {
			Name        graphql.String
			Label       graphql.String
			Description graphql.String
			Vars        []graphql.String
			Meta        struct {
				Contact struct {
					ContactName    graphql.String "graphql:\"Contact_Name\""
					ContactEmail   graphql.String "graphql:\"Contact_Email\""
					ContactPhone   graphql.String "graphql:\"Contact_Phone\""
					ContactWebsite graphql.String "graphql:\"Contact_Website\""
				} "graphql:\"Contact\""
				CensusReleases []struct {
					CensusReleaseDescription graphql.String "graphql:\"Census_Release_Description\""
					CensusReleaseNumber      graphql.String "graphql:\"Census_Release_Number\""
					ReleaseDate              graphql.String "graphql:\"Release_Date\""
				} "graphql:\"Census_Releases\""
				DatasetMnemonic2011        graphql.String   "graphql:\"Dataset_Mnemonic_2011\""
				DatasetPopulation          graphql.String   "graphql:\"Dataset_Population\""
				DisseminationSource        graphql.String   "graphql:\"Dissemination_Source\""
				GeographicCoverage         graphql.String   "graphql:\"Geographic_Coverage\""
				GeographicVariableMnemonic graphql.String   "graphql:\"Geographic_Variable_Mnemonic\""
				LastUpdated                graphql.String   "graphql:\"Last_Updated\""
				Keywords                   []graphql.String "graphql:\"Keywords\""
				Publications               []struct {
					PublisherName    graphql.String "graphql:\"Publisher_Name\""
					PublicationTitle graphql.String "graphql:\"Publication_Title\""
					PublisherWebsite graphql.String "graphql:\"Publisher_Website\""
				} "graphql:\"Publications\""
				RelatedDatasets  []graphql.String "graphql:\"Related_Datasets\""
				ReleaseFrequency graphql.String   "graphql:\"Release_Frequency\""
				StatisticalUnit  struct {
					StatisticalUnit            graphql.String "graphql:\"Statistical_Unit\""
					StatisticalUnitDescription graphql.String "graphql:\"Statistical_Unit_Description\""
				} "graphql:\"Statistical_Unit\""
				UniqueUrl graphql.String "graphql:\"Unique_Url\""
				Version   graphql.String "graphql:\"Version\""
			}
		} "graphql:\"tables(names: $vars)\""
	}{Tables: []struct {
		Name        graphql.String
		Label       graphql.String
		Description graphql.String
		Vars        []graphql.String
		Meta        struct {
			Contact struct {
				ContactName    graphql.String "graphql:\"Contact_Name\""
				ContactEmail   graphql.String "graphql:\"Contact_Email\""
				ContactPhone   graphql.String "graphql:\"Contact_Phone\""
				ContactWebsite graphql.String "graphql:\"Contact_Website\""
			} "graphql:\"Contact\""
			CensusReleases []struct {
				CensusReleaseDescription graphql.String "graphql:\"Census_Release_Description\""
				CensusReleaseNumber      graphql.String "graphql:\"Census_Release_Number\""
				ReleaseDate              graphql.String "graphql:\"Release_Date\""
			} "graphql:\"Census_Releases\""
			DatasetMnemonic2011        graphql.String   "graphql:\"Dataset_Mnemonic_2011\""
			DatasetPopulation          graphql.String   "graphql:\"Dataset_Population\""
			DisseminationSource        graphql.String   "graphql:\"Dissemination_Source\""
			GeographicCoverage         graphql.String   "graphql:\"Geographic_Coverage\""
			GeographicVariableMnemonic graphql.String   "graphql:\"Geographic_Variable_Mnemonic\""
			LastUpdated                graphql.String   "graphql:\"Last_Updated\""
			Keywords                   []graphql.String "graphql:\"Keywords\""
			Publications               []struct {
				PublisherName    graphql.String "graphql:\"Publisher_Name\""
				PublicationTitle graphql.String "graphql:\"Publication_Title\""
				PublisherWebsite graphql.String "graphql:\"Publisher_Website\""
			} "graphql:\"Publications\""
			RelatedDatasets  []graphql.String "graphql:\"Related_Datasets\""
			ReleaseFrequency graphql.String   "graphql:\"Release_Frequency\""
			StatisticalUnit  struct {
				StatisticalUnit            graphql.String "graphql:\"Statistical_Unit\""
				StatisticalUnitDescription graphql.String "graphql:\"Statistical_Unit_Description\""
			} "graphql:\"Statistical_Unit\""
			UniqueUrl graphql.String "graphql:\"Unique_Url\""
			Version   graphql.String "graphql:\"Version\""
		}
	}{{Name: "LC2101EW", Label: "Ethnic group by sex by age", Description: "This dataset provides 2011 Census estimates that classify usual residents in England and Wales by ethnic group, by sex and by age. The estimates are as at census day, 27 March 2011.\n\nThis information helps public bodies meet statutory obligations relating to race equality. It is also used for resource allocation and to develop and monitor policy on improving the life-chances for disadvantaged groups, including many ethnic minority groups.\n\nThe statistics also provide a better understanding of communities and are used for the government-wide race equality and community cohesion strategy, which seeks to improve race equality outcomes in areas such as housing, education, health and criminal justice for all groups across society.", Vars: []graphql.String{"Region", "Ethnic Group", "Sex", "Age"}, Meta: struct {
		Contact struct {
			ContactName    graphql.String "graphql:\"Contact_Name\""
			ContactEmail   graphql.String "graphql:\"Contact_Email\""
			ContactPhone   graphql.String "graphql:\"Contact_Phone\""
			ContactWebsite graphql.String "graphql:\"Contact_Website\""
		} "graphql:\"Contact\""
		CensusReleases []struct {
			CensusReleaseDescription graphql.String "graphql:\"Census_Release_Description\""
			CensusReleaseNumber      graphql.String "graphql:\"Census_Release_Number\""
			ReleaseDate              graphql.String "graphql:\"Release_Date\""
		} "graphql:\"Census_Releases\""
		DatasetMnemonic2011        graphql.String   "graphql:\"Dataset_Mnemonic_2011\""
		DatasetPopulation          graphql.String   "graphql:\"Dataset_Population\""
		DisseminationSource        graphql.String   "graphql:\"Dissemination_Source\""
		GeographicCoverage         graphql.String   "graphql:\"Geographic_Coverage\""
		GeographicVariableMnemonic graphql.String   "graphql:\"Geographic_Variable_Mnemonic\""
		LastUpdated                graphql.String   "graphql:\"Last_Updated\""
		Keywords                   []graphql.String "graphql:\"Keywords\""
		Publications               []struct {
			PublisherName    graphql.String "graphql:\"Publisher_Name\""
			PublicationTitle graphql.String "graphql:\"Publication_Title\""
			PublisherWebsite graphql.String "graphql:\"Publisher_Website\""
		} "graphql:\"Publications\""
		RelatedDatasets  []graphql.String "graphql:\"Related_Datasets\""
		ReleaseFrequency graphql.String   "graphql:\"Release_Frequency\""
		StatisticalUnit  struct {
			StatisticalUnit            graphql.String "graphql:\"Statistical_Unit\""
			StatisticalUnitDescription graphql.String "graphql:\"Statistical_Unit_Description\""
		} "graphql:\"Statistical_Unit\""
		UniqueUrl graphql.String "graphql:\"Unique_Url\""
		Version   graphql.String "graphql:\"Version\""
	}{Contact: struct {
		ContactName    graphql.String "graphql:\"Contact_Name\""
		ContactEmail   graphql.String "graphql:\"Contact_Email\""
		ContactPhone   graphql.String "graphql:\"Contact_Phone\""
		ContactWebsite graphql.String "graphql:\"Contact_Website\""
	}{ContactName: "Census Customer Services", ContactEmail: "census.customerservices@ons.gov.uk", ContactPhone: "01329 444 972", ContactWebsite: "https://www.ons.gov.uk/census/censuscustomerservices"}, CensusReleases: []struct {
		CensusReleaseDescription graphql.String "graphql:\"Census_Release_Description\""
		CensusReleaseNumber      graphql.String "graphql:\"Census_Release_Number\""
		ReleaseDate              graphql.String "graphql:\"Release_Date\""
	}{{CensusReleaseDescription: "Example release: ethnicity, national identity, language and religion", CensusReleaseNumber: "2", ReleaseDate: "30/07/2013"}}, DatasetMnemonic2011: "LC2101EW", DatasetPopulation: "All usual residents", DisseminationSource: "Census 2011", GeographicCoverage: "England and Wales", GeographicVariableMnemonic: "Region", LastUpdated: "30/07/2013", Keywords: []graphql.String{"Ethnic group", "Sex", "Age"}, Publications: []struct {
		PublisherName    graphql.String "graphql:\"Publisher_Name\""
		PublicationTitle graphql.String "graphql:\"Publication_Title\""
		PublisherWebsite graphql.String "graphql:\"Publisher_Website\""
	}{}, RelatedDatasets: []graphql.String{"LC2107EW"}, ReleaseFrequency: "", StatisticalUnit: struct {
		StatisticalUnit            graphql.String "graphql:\"Statistical_Unit\""
		StatisticalUnitDescription graphql.String "graphql:\"Statistical_Unit_Description\""
	}{StatisticalUnit: "People", StatisticalUnitDescription: "People living in England and Wales"}, UniqueUrl: "", Version: "1"}}}}}
}
