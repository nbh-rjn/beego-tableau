package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

/*
sample xml response

<?xml version='1.0' encoding='UTF-8'?>
<tsResponse xmlns="http://tableau.com/api" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://tableau.com/api https://help.tableau.com/samples/en-us/rest_api/ts-api_3_4.xsd">
    <credentials token="abced" estimatedTimeToExpiration="173:21:12">
        <site id="abcd" contentUrl="abcd"/>
        <user id="abcd"/>
    </credentials>
</tsResponse>

*/

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

// communicating with tableau api
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
