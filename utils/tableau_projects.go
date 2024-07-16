package utils

import (
	"encoding/xml"
	"net/http"
)

func Tableau_get_projects(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + site_id + "/projects"

	// make new get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", token)

	// send using client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ExtractProjectNames(xmlData string) ([]string, []string, error) {

	type Owner struct {
		ID string `xml:"id,attr"`
	}
	type Project struct {
		ID                 string `xml:"id,attr"`
		Name               string `xml:"name,attr"`
		Description        string `xml:"description,attr"`
		CreatedAt          string `xml:"createdAt,attr"`
		UpdatedAt          string `xml:"updatedAt,attr"`
		ContentPermissions string `xml:"contentPermissions,attr"`
		Owner              Owner  `xml:"owner"`
	}
	type Pagination struct {
		PageNumber     int `xml:"pageNumber,attr"`
		PageSize       int `xml:"pageSize,attr"`
		TotalAvailable int `xml:"totalAvailable,attr"`
	}
	type TSResponse struct {
		XMLName    xml.Name   `xml:"tsResponse"`
		Pagination Pagination `xml:"pagination"`
		Projects   []Project  `xml:"projects>project"`
	}

	// unmarshaling
	var tsResponse TSResponse
	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		return nil, nil, err
	}

	// Extract label names from struct
	var projectNames []string
	var projectIDs []string
	for _, project := range tsResponse.Projects {
		projectNames = append(projectNames, project.Name)
		projectIDs = append(projectIDs, project.ID)

	}
	return projectNames, projectIDs, nil
}
