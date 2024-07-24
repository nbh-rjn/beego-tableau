package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"encoding/json"
	"net/http"
)

func (c *TableauController) DownloadDataSource() {

	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	var request models.DownloadRequest

	// dont use c.bindJSON
	if err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	filePath := "storage/download.tds"
	// utility function to communicate with Tableau API
	if err := lib.TableauDownloadDataSource(models.Get_token(), request.SiteID, request.DatasourceID, filePath); err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data sources from Tableau")
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Download(filePath)

}
