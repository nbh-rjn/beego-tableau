package lib

import (
	"beego-project/models"
	"bytes"
	"fmt"
	"net/http"
)

func MakeRequest(url string, payload string, method string, contenttype string) (*http.Response, error) {

	// Create PUT request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Set("X-Tableau-Auth", models.Get_token())
	if contenttype != "" {
		req.Header.Set("Content-Type", "application/"+contenttype)
	}

	// HTTP client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil

}
