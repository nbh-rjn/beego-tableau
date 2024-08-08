package lib

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
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
	response, err := TableauRequest(url, "", "GET", "")
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	xmlData := string(responseBody)
	var attributes []map[string]interface{}

	switch param {

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

	// save attributes to database
	if err := models.SaveAttributesDB(param, models.GetSiteID(), attributes); err != nil {
		return nil, err
	}

	return attributes, nil

}
