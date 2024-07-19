package utils

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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

func TableauGetDataSources(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.4/sites/" + site_id + "/datasources"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Tableau-Auth", token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func TableauDownloadDataSource(token string, siteID string, datasourceID string) error {
	url := "https://10ax.online.tableau.com/api/3.4/sites/" + siteID + "/datasources/" + datasourceID + "/content"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Tableau-Auth", token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	out, err := os.Create("downloaded_file.tdsx")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer out.Close()

	// Copy the downloaded file content from response body to the local file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	fmt.Println("tdsx file downloaded successfully.")

	r, err := zip.OpenReader(out.Name())
	if err != nil {
		fmt.Println("Error opening .tdsx file as zip:", err)
		return err
	}
	defer r.Close()

	// Find the .tds file within the zip archive
	var tdsFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".tds") {
			tdsFile = f
			break
		}
	}

	if tdsFile == nil {
		fmt.Println("No .tds file found in .tdsx archive")
		return err
	}

	// Create a new file to save the extracted .tds content
	outFile, err := os.Create("extracted_file.tds")
	if err != nil {
		fmt.Println("Error creating .tds file:", err)
		return err
	}
	defer outFile.Close()

	// Open the .tds file from the zip archive
	rc, err := tdsFile.Open()
	if err != nil {
		fmt.Println("Error opening .tds file from zip:", err)
		return err
	}
	defer rc.Close()

	// Copy the contents of the .tds file from zip to the new file
	_, err = io.Copy(outFile, rc)
	if err != nil {
		fmt.Println("Error saving .tds file:", err)
		return err
	}

	fmt.Println("File extracted successfully:", outFile.Name())

	return nil
}
