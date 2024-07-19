package lib

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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

	r, err := zip.OpenReader(out.Name())
	if err != nil {
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
		return err
	}

	// Create a new file to save the extracted .tds content
	outFile, err := os.Create("extracted_file.tds")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Open the .tds file from the zip archive
	rc, err := tdsFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Copy the contents of the .tds file from zip to the new file
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	return nil
}
