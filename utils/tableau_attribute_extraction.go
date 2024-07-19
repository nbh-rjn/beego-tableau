package utils

import (
	"encoding/xml"
	"fmt"
)

type SiteL struct {
	ID string `xml:"id,attr"`
}

type LabelValue struct {
	XMLName     xml.Name `xml:"labelValue"`
	Name        string   `xml:"name,attr"`
	Category    string   `xml:"category,attr"`
	Description string   `xml:"description,attr"`
	Internal    bool     `xml:"internal,attr"`
	Elevated    bool     `xml:"elevatedDefault,attr"`
	BuiltIn     bool     `xml:"builtIn,attr"`
	Site        SiteL    `xml:"site"`
}

type LabelValueList struct {
	XMLName     xml.Name     `xml:"labelValueList"`
	LabelValues []LabelValue `xml:"labelValue"`
}

type TsResponseL struct {
	XMLName        xml.Name       `xml:"tsResponse"`
	LabelValueList LabelValueList `xml:"labelValueList"`
}

type Datasource struct {
	XMLName xml.Name `xml:"datasource"`
}

type DatasourceElement struct {
	Datasource
	ContentUrl string `xml:"contentUrl,attr"`
	Name       string `xml:"name,attr"`
	Id         string `xml:"id,attr"`
}

type Datasources struct {
	XMLName    xml.Name            `xml:"datasources"`
	Datasource []DatasourceElement `xml:"datasource"`
}

type DatasourceResponse struct {
	TSResponse
	Datasources Datasources `xml:"datasources"`
}

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

type ProjectResponse struct {
	TSResponse
	Pagination Pagination `xml:"pagination"`
	Projects   []Project  `xml:"projects>project"`
}

func ExtractLabelNames(xmlData string) ([]string, error) {

	// unmarshaling
	var tsResponse TsResponseL
	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		return nil, err
	}

	// Extract label names from struct
	var labelNames []string
	for _, LabelValue := range tsResponse.LabelValueList.LabelValues {
		labelNames = append(labelNames, LabelValue.Name)
	}
	return labelNames, nil
}

func ExtractDataSources(xmlData string) ([]string, []string) {
	var tsResponse DatasourceResponse
	var names []string
	var ids []string

	err := xml.Unmarshal([]byte(xmlData), &tsResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling XML: %v\n", err)
		return names, ids
	}

	for _, ds := range tsResponse.Datasources.Datasource {
		names = append(names, ds.Name)
		ids = append(ids, ds.Id)
	}

	return names, ids
}

func ExtractProjectNames(xmlData string) ([]string, []string, error) {

	// unmarshaling
	var tsResponse ProjectResponse
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
