package lib

import (
	"beego-project/models"
	"encoding/json"
	"fmt"
	"io"
)

func GetAssetID(assetType string, databaseName string, tableName string, columnName string) (string, error) {
	payload := fmt.Sprintf(`
				{
					"query": "query tableQuery{\n    databases (filter: { name: \"%s\"}) {\n        tables (filter: { name: \"%s\"}){\n            luid\n        }\n    }\n}",
					"variables": {}
				}`, databaseName, tableName)

	// API endpoint
	url := "https://10ax.online.tableau.com/api/metadata/graphql"
	resp, err := MakeRequest(url, payload, "POST", "json")

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseData models.TableResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", err
	}

	if len(responseData.Data.Databases) == 0 || len(responseData.Data.Databases[0].Tables) == 0 {
		return "", fmt.Errorf("no tables found")
	}

	return responseData.Data.Databases[0].Tables[0].LUID, nil

}

func CreateCategory(siteID string, category string) error {

	// need to use api v3.21; others dont work
	url := models.TableauURL() + "sites/" + siteID + "/labelCategories"

	// request body
	payload := fmt.Sprintf(
		`<tsRequest>
		<labelCategory name="%s"
			description="%s" />
	</tsRequest>`, category, category)

	response, err := MakeRequest(url, payload, "POST", "xml")

	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func CreateLabelValue(siteID string, label string, category string) error {
	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
		   <labelValue name="%s"
		     category="%s"
		     description="Created via API" />
		</tsRequest>`, label, category)

	// need v3.21
	url := models.TableauURL() + "sites/" + siteID + "/labelValues"

	// new put request
	response, err := MakeRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func ApplyLabelValue(siteID string, asset string, columnID string, label string) error {
	// XML payload for applying label value
	payload := fmt.Sprintf(`
		<tsRequest>
		  <contentList>
		    <content contentType="column" id="%s" />
		  </contentList>
		  <label
		      value="%s"/>
		</tsRequest>`, columnID, label)

	// API endpoint
	url := models.TableauURL() + "sites/" + siteID + "/labels"

	// Create PUT request
	response, err := MakeRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func GetColumns(databaseName string, tableName string) (map[string]string, error) {
	payload := fmt.Sprintf(`
		{
			"query": "query columnQuery{\n    databases (filter: { name: \"%s\"}) {\n        tables (filter: { name: \"%s\"}){\n            columns{\n                luid\n                 name\n           }\n        }\n    }\n}",
			"variables": {}
		}`, databaseName, tableName)

	// API endpoint
	url := "https://10ax.online.tableau.com/api/metadata/graphql"
	resp, err := MakeRequest(url, payload, "POST", "json")
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// extract json
	var responseData models.TableResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	if len(responseData.Data.Databases) == 0 || len(responseData.Data.Databases[0].Tables) == 0 || len(responseData.Data.Databases[0].Tables[0].Columns) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	// return map of col names and ids
	columns := make(map[string]string)
	for _, c := range responseData.Data.Databases[0].Tables[0].Columns {
		columns[c.Name] = c.LUID
	}

	return columns, nil

}
