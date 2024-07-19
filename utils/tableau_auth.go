package utils

import (
	"encoding/xml"
	"fmt"
)

// struct to hold the XML response
type TSResponse struct {
	XMLName xml.Name `xml:"tsResponse"`
}

type AuthResponse struct {
	TSResponse
	Credentials Credentials `xml:"credentials"`
}

type Credentials struct {
	Token                     string `xml:"token,attr"`
	EstimatedTimeToExpiration string `xml:"estimatedTimeToExpiration,attr"`
	Site                      Site   `xml:"site"`
	User                      User   `xml:"user"`
}

type Site struct {
	ID         string `xml:"id,attr"`
	ContentURL string `xml:"contentUrl,attr"`
}

type User struct {
	ID string `xml:"id,attr"`
}

// to extract token from xml response
func ExtractToken(responseBody string) (string, error) {

	// unmarshal XML using our structs
	var response AuthResponse
	err := xml.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return "", err
	}

	// Extract token from the struct
	token := response.Credentials.Token
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
