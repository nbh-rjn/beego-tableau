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
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau API
	if err := lib.TableauDownloadDataSource(models.Get_token(), request.SiteID, request.DatasourceID); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	filePath := "extracted_file.tds"
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Download(filePath)

}
