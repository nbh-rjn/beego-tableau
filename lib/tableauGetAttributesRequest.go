package lib

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// accepts attribute "datalabels", "projects", or "datasources" as param and returns info about the attribute
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
	url := models.TableauURL() + "sites/" + models.GetSiteID() + attribute
	response, err := MakeRequest(url, "", "GET", "")
	if err != nil {
		return nil, err
	}

	attributes, err := extractAttributes(response, param)
	if err != nil {
		return nil, err
	}

	// save attributes to database
	if err := models.SaveAttributesDB(param, models.GetSiteID(), attributes); err != nil {
		return nil, err
	}

	return attributes, nil

}

// extracts attributes from an http response and returns them
func extractAttributes(response *http.Response, attributeType string) ([]map[string]interface{}, error) {

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
