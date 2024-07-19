package utils

import (
	"encoding/xml"
	"fmt"
)

type Datasource struct {
	XMLName xml.Name `xml:"datasource"`
}

type DatasourceElement struct {
	Datasource
	ContentUrl string `xml:"contentUrl,attr"`
	Name       string `xml:"name,attr"`
	Id         string `xml:"id,attr"`
}

type Datasources struct {
	XMLName    xml.Name            `xml:"datasources"`
	Datasource []DatasourceElement `xml:"datasource"`
}

type DatasourceResponse struct {
	TSResponse
	Datasources Datasources `xml:"datasources"`
}

func ExtractDataSources(xmlData string) ([]string, []string) {
	var tsResponse DatasourceResponse
	var names []string
	var ids []string

	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling XML: %v\n", err)
		return names, ids
	}

	for _, ds := range tsResponse.Datasources.Datasource {
		names = append(names, ds.Name)
		ids = append(ids, ds.Id)
	}

	return names, ids
}
