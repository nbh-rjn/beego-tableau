package lib

import (
	"beego-project/models"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// publishes file
func PublishDatasource(filePath string, filenameTDS string, datasourceName string, projectID string) (string, error) {
	url := models.TableauURL() + "sites/" + models.Get_siteID() + "/datasources?datasourceType=tds&overwrite=true"

	// construct request payload
	requestPayload := fmt.Sprintf(
		`<tsRequest><datasource name="%s">
			<project id="%s" />
		</datasource></tsRequest>`, datasourceName, projectID)

	// buffer to store request
	requestBody := &bytes.Buffer{}

	// multipart writer
	writer := multipart.NewWriter(requestBody)

	// add request_payload field
	if err := writer.WriteField("request_payload", requestPayload); err != nil {
		return "", err
	}

	// add tableau_datasource file field
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// create form file part
	part, err := writer.CreateFormFile("tableau_datasource", filenameTDS)
	if err != nil {
		return "", err
	}

	// copy file content into the part
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	// close multipart writer, write the boundary
	if err := writer.Close(); err != nil {
		return "", err
	}

	// create HTTP request
	request, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return "", err
	}

	// set headers
	// do not replace the sprintf with string concatenation
	request.Header.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()))
	request.Header.Set("X-Tableau-Auth", models.Get_token())

	// send HTTP request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check if the response status code is not OK (200)
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}

	var extractresponse models.PublishDSResponse

	// Unmarshal the XML data into the struct
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if err := xml.Unmarshal(responseBody, &extractresponse); err != nil {
		return "", err
	}

	return extractresponse.Datasource.Id, nil
}
