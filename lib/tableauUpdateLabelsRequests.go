package lib

import (
	"beego-project/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetTableID(databaseName string, tableName string) string {
	return ""
}

func CreateCategory(category string) error {

	// need to use api v3.22; others dont work
	url := "https://10ax.online.tableau.com/api/3.21/sites/2ff64d57-b7c1-4e99-803c-13bb81ae0371/labelCategories"

	// request body
	xmlBody := fmt.Sprintf(
		`<tsRequest>
		<labelCategory name="%s"
			description="%s" />
	</tsRequest>`, category, category)

	// new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(xmlBody))
	if err != nil {
		return err
	}

	// headers
	req.Header.Set("X-Tableau-Auth", models.Get_token())
	req.Header.Set("Content-Type", "application/xml")

	// send req w client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func GetColumnID(databaseName string, tableName string, columnName string) (string, error) {

	//  graphQL query
	// do not format for readability, it wont work
	query := fmt.Sprintf(`
	{
		"query": "query useMetadataApiToQueryOrdersDatabases{\n    databases (filter: { name: \"%s\"}) {\n        tables (filter: { name: \"%s\"}){\n            columns (filter: { name: \"%s\"}){\n                luid\n            }\n        }\n    }\n}",
		"variables": {}
	}`, databaseName, tableName, columnName)

	// cant use global url cuz no api version used here
	url := "https://10ax.online.tableau.com/api/metadata/graphql"

	// JSON payload for the POST request
	payload := []byte(query)

	// new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// headers
	req.Header.Set("X-Tableau-Auth", models.Get_token())
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// Extract the luid from the response
	databases, ok := data["data"].(map[string]interface{})["databases"].([]interface{})
	if !ok || len(databases) == 0 {
		return "", fmt.Errorf("no databases found")
	}

	tables := databases[0].(map[string]interface{})["tables"].([]interface{})
	if len(tables) == 0 {
		return "", fmt.Errorf("no tables found in the database")
	}

	columns := tables[0].(map[string]interface{})["columns"].([]interface{})
	if len(columns) == 0 {
		return "", fmt.Errorf("no columns found in the table")
	}

	firstColumn := columns[0].(map[string]interface{})
	luid, ok := firstColumn["luid"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract luid")
	}

	return luid, nil
}

func CreateLabelValue(label string, category string) error {
	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
		   <labelValue name="%s"
		     category="%s"
		     description="Created via API" />
		</tsRequest>`, label, category)

	// need v3.21
	url := "https://10ax.online.tableau.com/api/3.21/sites/2ff64d57-b7c1-4e99-803c-13bb81ae0371/labelValues"

	// new put request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	// headers
	req.Header.Set("X-Tableau-Auth", models.Get_token())
	req.Header.Set("Content-Type", "application/xml")

	// send req
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func ApplyLabelValue(asset string, columnID string, label string) error {
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
	url := "https://10ax.online.tableau.com/api/3.21/sites/2ff64d57-b7c1-4e99-803c-13bb81ae0371/labels"

	// Create PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	// Headers
	req.Header.Set("X-Tableau-Auth", models.Get_token())
	req.Header.Set("Content-Type", "application/xml")

	// HTTP client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
