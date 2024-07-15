package utils

import (
	"encoding/xml"
	"net/http"
)

func Tableau_get_data_label_values(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + site_id + "/labelValues"

	// make new get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", token)

	// send using client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ExtractLabelNames(xmlData string) ([]string, error) {

	// structs for unmarshalling
	type SiteL struct {
		ID string `xml:"id,attr"`
	}

	type LabelValue struct {
		XMLName     xml.Name `xml:"labelValue"`
		Name        string   `xml:"name,attr"`
		Category    string   `xml:"category,attr"`
		Description string   `xml:"description,attr"`
		Internal    bool     `xml:"internal,attr"`
		Elevated    bool     `xml:"elevatedDefault,attr"`
		BuiltIn     bool     `xml:"builtIn,attr"`
		Site        SiteL    `xml:"site"`
	}

	type LabelValueList struct {
		XMLName     xml.Name     `xml:"labelValueList"`
		LabelValues []LabelValue `xml:"labelValue"`
	}

	type TsResponseL struct {
		XMLName        xml.Name       `xml:"tsResponse"`
		LabelValueList LabelValueList `xml:"labelValueList"`
	}

	// unmarshaling
	var tsResponse TsResponseL
	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		return nil, err
	}

	// Extract label names from struct
	var labelNames []string
	for _, LabelValue := range tsResponse.LabelValueList.LabelValues {
		labelNames = append(labelNames, LabelValue.Name)
	}
	return labelNames, nil
}
