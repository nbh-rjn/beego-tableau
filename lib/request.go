package lib

import (
	"beego-project/models"
	"bytes"
	"fmt"
	"net/http"
)

func TableauRequest(url string, payload string, method string, contentType string) (*http.Response, error) {

	// Create PUT request
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Headers
	request.Header.Set("X-Tableau-Auth", models.GetToken())
	if contentType != "" {
		request.Header.Set("Content-Type", "application/"+contentType)
	}

	// HTTP client
	client := &http.Client{}

	// Send request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Check response status
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return response, nil

}
