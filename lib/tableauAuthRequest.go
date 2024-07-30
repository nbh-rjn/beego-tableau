package lib

import (
	"beego-project/models"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func TableauAuthRequest(patName string, patSecret string, contentURL string) (string, string, error) {

	url := models.TableauURL() + "auth/signin"
	payload := fmt.Sprintf(
		`<tsRequest>
		            <credentials personalAccessTokenName="%s"
		                personalAccessTokenSecret="%s">
		                	<site contentUrl="%s" />
		            </credentials>
		        </tsRequest>`, patName, patSecret, contentURL)

	// new post request
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", "application/xml")

	// send request using client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	xmlData := string(responseBody)

	// unmarshal XML using our structs
	var tsResponse models.AuthResponse

	if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
		return "", "", err
	}

	// Extract token from the struct
	token := tsResponse.Credentials.Token
	siteID := tsResponse.Credentials.Site.ID

	return token, siteID, nil

}
