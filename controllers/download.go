package controllers

import (
	"beego-project/lib"
	"beego-project/models"
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

	filePath := "storage/download.tds"

	// download this tds from this site and save to this filepath
	if err := lib.TableauDownloadDataSource(request.DatasourceID); err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data sources from Tableau")
	}

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Download(filePath)

}
