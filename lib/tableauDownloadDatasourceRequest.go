package lib

import (
	"archive/zip"
	"beego-project/models"
	"io"
	"os"
	"strings"
)

// downloads .tds file of the given id and returns path where it has been stored
func TableauDownloadDataSource(datasourceID string) (string, error) {
	filePath := "storage/download.tds"

	url := models.TableauURL() + "sites/" + models.Get_siteID() + "/datasources/" + datasourceID + "/content"
	response, err := MakeRequest(url, "", "GET", "")
	if err != nil {
		return filePath, err
	}

	out, err := os.Create("storage/" + datasourceID + ".tdsx")
	if err != nil {
		return filePath, err
	}
	defer out.Close()

	// Copy the downloaded file content from response body to the local file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return filePath, err
	}

	r, err := zip.OpenReader(out.Name())
	if err != nil {
		return filePath, err
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
		return filePath, err
	}

	// Create a new file to save the extracted .tds content
	outFile, err := os.Create(filePath)
	if err != nil {
		return filePath, err
	}
	defer outFile.Close()

	// Open the .tds file from the zip archive
	rc, err := tdsFile.Open()
	if err != nil {
		return filePath, err
	}
	defer rc.Close()

	// Copy the contents of the .tds file from zip to the new file
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return filePath, err
	}

	return filePath, nil
}
