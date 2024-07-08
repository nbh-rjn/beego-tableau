package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

// TSResponse struct to hold the XML response
type TSResponse struct {
	XMLName     xml.Name    `xml:"tsResponse"`
	Credentials Credentials `xml:"credentials"`
}

// Credentials struct to represent the credentials node
type Credentials struct {
	Token                     string `xml:"token,attr"`
	EstimatedTimeToExpiration string `xml:"estimatedTimeToExpiration,attr"`
	Site                      Site   `xml:"site"`
	User                      User   `xml:"user"`
}

// Site struct to represent the site node
type Site struct {
	ID         string `xml:"id,attr"`
	ContentURL string `xml:"contentUrl,attr"`
}

// User struct to represent the user node
type User struct {
	ID string `xml:"id,attr"`
}

func Get_token(resp_body string) string {
	var response TSResponse
	err := xml.Unmarshal([]byte(resp_body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return ""
	}

	// Extract token from the struct
	token := response.Credentials.Token
	fmt.Println("Token:", token)
	return token
}

func Make_xml(pat string, pats string, siteid string) string {
	xml := fmt.Sprintf(
		`<tsRequest>
		            <credentials personalAccessTokenName="%s"
		                personalAccessTokenSecret="%s">
		                	<site contentUrl="%s" />
		            </credentials>
		        </tsRequest>`, pat, pats, siteid)
	return xml
}

func Tableau_req(xmlData string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.4/auth/signin"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(xmlData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil

}
