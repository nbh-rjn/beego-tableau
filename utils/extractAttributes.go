package utils

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func ExtractAttributes(response *http.Response, attributeType string) ([]map[string]interface{}, error) {

	var tag1, tag2 []string
	tagName1, tagName2 := "name", "id"

	var attributes []map[string]interface{}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	xmlData := string(responseBody)

	switch attributeType {

	case "datalabels":

		var tsResponse models.LabelValueResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, labelValue := range tsResponse.LabelValueList.LabelValues {
			tag1 = append(tag1, labelValue.Name)
			tag2 = append(tag2, labelValue.Category)
		}

		tagName1, tagName2 = "name", "category"

	case "datasources":

		var tsResponse models.DatasourceResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, datasource := range tsResponse.Datasources.Datasource {
			tag1 = append(tag1, datasource.Name)
			tag2 = append(tag2, datasource.Id)
		}

	case "projects":

		var tsResponse models.ProjectResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, project := range tsResponse.Projects {
			tag1 = append(tag1, project.Name)
			tag2 = append(tag2, project.ID)
		}

	default:
		return nil, fmt.Errorf("invalid attribute type")
	}

	for i := 0; i < len(tag1); i++ {
		attributes = append(attributes, map[string]interface{}{
			tagName1: tag1[i], tagName2: tag2[i],
		})
	}
	return attributes, nil

}
