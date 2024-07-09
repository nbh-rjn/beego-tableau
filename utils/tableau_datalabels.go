package utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func Tableau_get_data_label_values(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + site_id + "/labelValues"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tableau-Auth", token)

	//fmt.Printf("\t token id: %s \n", token)
	//fmt.Printf("\t url: %s \n", url)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ExtractLabelNames(xmlData string) ([]string, error) {
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

	var tsResponse TsResponseL
	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		return nil, err
	}

	// Extract label names
	var labelNames []string
	for _, lv := range tsResponse.LabelValueList.LabelValues {
		labelNames = append(labelNames, lv.Name)
	}
	fmt.Println(labelNames)
	return labelNames, nil
}
