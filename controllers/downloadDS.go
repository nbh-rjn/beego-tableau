package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type DownloadRequest struct {
	SiteRequest
	DatasourceID string `json:"datasourceID"`
}

func (c *TableauController) DownloadDataSource() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	var request DownloadRequest

	// dont use Ctx.Input.RequestBody
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request)

	// check for correct request format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau API
	err = utils.TableauDownloadDataSource(models.Get_token(), request.SiteID, request.DatasourceID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	filePath := "extracted_file.tds"

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Download(filePath)

	fmt.Println("downloaded . . .")

}
