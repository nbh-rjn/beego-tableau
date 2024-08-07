package lib

import (
	"beego-project/models"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// returns first label on an asset
func TableauGetAssetLabel(assetType string, assetID string) (string, error) {
	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
			<contentList>
   				<content contentType="%s"
      			id="%s" />
 			</contentList>
		</tsRequest>`, assetType, assetID)

	// need v3.21
	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/labels"

	// new put request
	response, err := TableauRequest(url, payload, "POST", "xml")
	if err != nil {
		return "", err
	}

	xmlData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseBody models.LabelResponse
	if err := xml.Unmarshal([]byte(xmlData), &responseBody); err != nil {
		return "", err
	}

	response.Body.Close()
	if len(responseBody.LabelList) == 0 {
		return "", nil
	}
	return responseBody.LabelList[0].Value, nil
}

func TableauGetAssetIDs(databaseName string, tableName string) (string, map[string]string, error) {
	payload := fmt.Sprintf(`
		{
			"query": "query columnQuery{\n    databases (filter: { name: \"%s\"}) {\n        tables (filter: { name: \"%s\"}){\n            luid\n            columns{\n                luid\n                 name\n           }\n        }\n    }\n}",
			"variables": {}
		}`, databaseName, tableName)

	// API endpoint
	url := "https://10ax.online.tableau.com/api/metadata/graphql"
	resp, err := TableauRequest(url, payload, "POST", "json")
	if err != nil {
		return "", nil, err

	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	// extract json
	var responseData models.TableResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", nil, err
	}

	// error handling
	if len(responseData.Data.Databases) == 0 {
		return "", nil, fmt.Errorf("could not retrieve asset id, no database found")
	}
	if len(responseData.Data.Databases[0].Tables) == 0 {
		return "", nil, fmt.Errorf("could not retrieve asset id, no table found")
	}
	if len(responseData.Data.Databases[0].Tables[0].Columns) == 0 {
		return responseData.Data.Databases[0].Tables[0].LUID, nil, fmt.Errorf("could not retrieve asset id, no columns found")
	}

	// return map of col names and ids
	columns := make(map[string]string)
	for _, c := range responseData.Data.Databases[0].Tables[0].Columns {
		columns[c.Name] = c.LUID
	}
	return responseData.Data.Databases[0].Tables[0].LUID, columns, nil

}
