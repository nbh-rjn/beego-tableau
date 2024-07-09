package utils

import (
	"encoding/xml"
	"fmt"
)

type TsResponse struct {
	XMLName     xml.Name    `xml:"tsResponse"`
	Datasources Datasources `xml:"datasources"`
}

type Datasources struct {
	XMLName    xml.Name     `xml:"datasources"`
	Datasource []Datasource `xml:"datasource"`
}

type Datasource struct {
	XMLName    xml.Name `xml:"datasource"`
	ContentUrl string   `xml:"contentUrl,attr"`
}

func Extract_data_sources_xml(xmlData string) []string {

	var tsResponse TsResponse
	var contentUrls []string

	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling XML: %v\n", err)
		return contentUrls
	}

	for _, ds := range tsResponse.Datasources.Datasource {
		contentUrls = append(contentUrls, ds.ContentUrl)
	}

	fmt.Println("Content URLs extracted:")
	for _, url := range contentUrls {
		fmt.Println(url)
	}
	return contentUrls
}
