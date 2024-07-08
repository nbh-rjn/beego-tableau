package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

/*
sample xml response to help in understanding

<?xml version='1.0' encoding='UTF-8'?>
<tsResponse xmlns="http://tableau.com/api" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://tableau.com/api https://help.tableau.com/samples/en-us/rest_api/ts-api_3_4.xsd">
    <credentials token="pi9TCS-4Ti-_xD0dN-H8ww|4qnj3HwbpVAdAdd7KkdnGGakiwxBnbcS|2ff64d57-b7c1-4e99-803c-13bb81ae0371" estimatedTimeToExpiration="173:21:12">
        <site id="2ff64d57-b7c1-4e99-803c-13bb81ae0371" contentUrl="testsiteintern"/>
        <user id="776f27ef-8997-422b-9348-a28d4f091fa4"/>
    </credentials>
</tsResponse>

*/

// struct to hold the XML response
type TSResponse struct {
	XMLName     xml.Name    `xml:"tsResponse"`
	Credentials Credentials `xml:"credentials"`
}

// struct to represent the credentials node
type Credentials struct {
	Token                     string `xml:"token,attr"`
	EstimatedTimeToExpiration string `xml:"estimatedTimeToExpiration,attr"`
	Site                      Site   `xml:"site"`
	User                      User   `xml:"user"`
}

// struct to represent the site node
type Site struct {
	ID         string `xml:"id,attr"`
	ContentURL string `xml:"contentUrl,attr"`
}

// struct to represent the user node
type User struct {
	ID string `xml:"id,attr"`
}

// to extract token from xml response
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

// construct xml request to tableau api
// using the info from the request to beego
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

// communicating with tableau api
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
