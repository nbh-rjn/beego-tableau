package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"context"
	"encoding/json"
	"net/http"
)

func (c *TableauController) DownloadDataSource() {

	c.EnableRender = false
	var request models.DownloadRequest

	// dont use c.bindJSON
	if err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// make api request
	fileName, err := lib.TableauDownloadDataSource(request.DatasourceID)

	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data sources from Tableau")
	}

	storage := models.GetStorage(context.TODO())
	if storage == nil {
		HandleError(c, http.StatusInternalServerError, "No storage handler found")
		return
	}
	// Read the file from storage
	fileData, err := storage.Read(context.TODO(), fileName)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to read file from storage")
		return
	}

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Body(fileData)

}
