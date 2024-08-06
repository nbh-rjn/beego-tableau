package lib

import (
	"beego-project/models"
	"encoding/json"
	"fmt"
	"io"
)

func TableauCreateCategory(category string) error {

	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/labelCategories"

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

	return err
}

func TableauLabelAsset(label string, category string, assetType string, assetID string) error {

	if label == "" || category == "" {
		return nil
	}

	// in case category doesnt exist
	TableauCreateCategory(category)

	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
		   <labelValue name="%s"
		     category="%s"
		     description="Created via API" />
		</tsRequest>`, label, category)

	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/labelValues"

	// new put request
	response, err := MakeRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	response.Body.Close()

	// XML payload for applying label value
	payload = fmt.Sprintf(`
		<tsRequest>
		  <contentList>
		    <content contentType="%s" id="%s" />
		  </contentList>
		  <label
		      value="%s"/>
		</tsRequest>`, assetType, assetID, label)

	url = models.TableauURL() + "sites/" + models.GetSiteID() + "/labels"

	// PUT request
	response, err = MakeRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func TableauGetAssetIDs(databaseName string, tableName string) (string, map[string]string, error) {
	payload := fmt.Sprintf(`
		{
			"query": "query columnQuery{\n    databases (filter: { name: \"%s\"}) {\n        tables (filter: { name: \"%s\"}){\n            luid\n            columns{\n                luid\n                 name\n           }\n        }\n    }\n}",
			"variables": {}
		}`, databaseName, tableName)

	// API endpoint
	url := "https://10ax.online.tableau.com/api/metadata/graphql"
	resp, err := MakeRequest(url, payload, "POST", "json")
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
