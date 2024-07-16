package utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func Extract_data_sources_xml(xmlData string) ([]string, []string) {
	type Datasource struct {
		XMLName    xml.Name `xml:"datasource"`
		ContentUrl string   `xml:"contentUrl,attr"`
		Name       string   `xml:"name,attr"`
		Id         string   `xml:"id,attr"`
	}

	type Datasources struct {
		XMLName    xml.Name     `xml:"datasources"`
		Datasource []Datasource `xml:"datasource"`
	}

	type TsResponse struct {
		XMLName     xml.Name    `xml:"tsResponse"`
		Datasources Datasources `xml:"datasources"`
	}

	var tsResponse TsResponse
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

func Tableau_get_data_sources(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.4/sites/" + site_id + "/datasources"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
