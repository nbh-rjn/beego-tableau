package utils

import (
	"beego-project/models"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func PublishDatasource(siteID string, datasourceName string) error {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + siteID + "/datasources?datasourceType=tds&overwrite=true"

	// Construct the request payload as XML
	requestPayload := `<tsRequest><datasource name="` + datasourceName + `"></datasource></tsRequest>`

	// Path to the XML file
	xmlFilePath := "xml.tds"

	// Create a new buffer to store the request body
	body := &bytes.Buffer{}

	// Create a multipart writer with a random boundary
	writer := multipart.NewWriter(body)

	// Add request_payload field
	_ = writer.WriteField("request_payload", requestPayload)

	// Add tableau_datasource file field
	file, err := os.Open(xmlFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create form file part with the correct field name and filename
	part, err := writer.CreateFormFile("tableau_datasource", xmlFilePath)
	if err != nil {
		return err
	}

	// Copy the file content into the part
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Close the multipart writer, this writes the boundary
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create HTTP request with the correct method, URL, and body
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Set headers, including Content-Type with the manually constructed boundary
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()))
	req.Header.Set("X-Tableau-Auth", models.Get_token())

	// Send HTTP request using a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is not OK (200)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	fmt.Println("Request successfully sent")
	return nil
}
