package lib

import (
	"archive/zip"
	"beego-project/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

// downloads .tds file of the given id and returns path where it has been stored
func TableauDownloadDataSource(datasourceID string) (string, error) {

	fileName := fmt.Sprintf("download-%s.tds", datasourceID)

	// get storage handler
	storage := models.GetStorage(context.TODO())
	if storage == nil {
		return "", fmt.Errorf("no storage handler found")
	}

	// make request to tableau
	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/datasources/" + datasourceID + "/content"
	response, err := TableauRequest(url, "", "GET", "")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// save response body .tdsx file to storage
	filenameTDSX := fmt.Sprintf("download-%s.tdsx", datasourceID)
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, response.Body); err != nil {
		return "", err
	}
	if err := storage.Write(context.TODO(), filenameTDSX, buffer.Bytes()); err != nil {
		return "", err
	}

	// read .tdsx file from storage
	dataTDSX, err := storage.Read(context.TODO(), filenameTDSX)
	if err != nil {
		return "", err
	}

	// Extract the .tds file from the .tdsx archive
	zipReader, err := zip.NewReader(bytes.NewReader(dataTDSX), int64(len(dataTDSX)))
	if err != nil {
		return "", err
	}
	var tdsFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, ".tds") {
			tdsFile = f
			break
		}
	}
	if tdsFile == nil {
		return "", fmt.Errorf(".tds file not found in .tdsx archive")
	}

	// Save the extracted .tds file to storage
	tdsFileData := new(bytes.Buffer)
	rc, err := tdsFile.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	if _, err := io.Copy(tdsFileData, rc); err != nil {
		return "", err
	}

	if err := storage.Write(context.TODO(), fileName, tdsFileData.Bytes()); err != nil {
		return "", err
	}

	return fileName, nil
}
