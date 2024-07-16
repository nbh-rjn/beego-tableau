package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauControllerDS struct {
	beego.Controller
}

func (c *TableauControllerDS) GetDataSources() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	type SiteRequest struct {
		SiteID string `json:"siteID"`
	}
	var req SiteRequest

	// DON'T REPLACE Ctx.Input.CopyBody WITH Ctx.Input.RequestBody
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &req)

	// check for correct request format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau API
	response, err := utils.Tableau_get_data_sources(models.Get_token(), req.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// read body of response
	bodyread, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	datasources, _ := utils.Extract_data_sources_xml(string(bodyread))

	// return info in response
	c.Data["json"] = map[string]interface{}{"datasources": datasources}
	c.ServeJSON()

}
