package lib

import (
	"bytes"
	"net/http"
)

func TableauAuthRequest(xmlData string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.4/auth/signin"

	// new post request
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(xmlData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/xml")

	// send request using client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}
