package lib

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func TableauGetAttributes(param string) ([]map[string]interface{}, error) {
	attributeMap := map[string]string{
		"datalabels":  "/labelValues",
		"datasources": "/datasources",
		"projects":    "/projects",
	}
	attribute, found := attributeMap[param]
	if !found {
		return nil, fmt.Errorf("invalid attribute")
	}
	url := models.TableauURL() + "sites/" + models.Get_siteID() + attribute

	// make new get request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Tableau-Auth", models.Get_token())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	attributes, err := extractAttributes(response, param)
	if err != nil {
		return nil, err
	}

	// save attributes to database
	if err := models.SaveAttributesDB(param, models.Get_siteID(), attributes); err != nil {
		return nil, err
	}

	return attributes, nil

}

func extractAttributes(response *http.Response, attributeType string) ([]map[string]interface{}, error) {

	//var tag1, tag2 []string
	//tagName1, tagName2 := "name", "id"

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
			attributes = append(attributes, map[string]interface{}{
				"name": labelValue.Name, "category": labelValue.Category,
			})
		}

	case "datasources":

		var tsResponse models.DatasourceResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, datasource := range tsResponse.Datasources.Datasource {
			attributes = append(attributes, map[string]interface{}{
				"name": datasource.Name, "id": datasource.Id,
			})
		}

	case "projects":

		var tsResponse models.ProjectResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, project := range tsResponse.Projects {
			attributes = append(attributes, map[string]interface{}{
				"name": project.Name, "id": project.ID,
			})
		}

	default:
		return nil, fmt.Errorf("invalid attribute type")
	}

	return attributes, nil

}
