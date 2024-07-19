package utils

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// to extract token from xml response
func ExtractToken(response *http.Response) (string, error) {

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	xmlData := string(responseBody)

	// unmarshal XML using our structs
	var tsResponse models.AuthResponse

	if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
		return "", err
	}

	// Extract token from the struct
	token := tsResponse.Credentials.Token
	return token, nil
}

// construct xml request to tableau api
// using the info from the request to beego
func CredentialsXML(PersonalAccessTokenName string, PersonalAccessTokenSecret string, ContentUrl string) string {
	xml := fmt.Sprintf(
		`<tsRequest>
		            <credentials personalAccessTokenName="%s"
		                personalAccessTokenSecret="%s">
		                	<site contentUrl="%s" />
		            </credentials>
		        </tsRequest>`, PersonalAccessTokenName, PersonalAccessTokenSecret, ContentUrl)
	return xml
}
