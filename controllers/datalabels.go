package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauControllerDL struct {
	beego.Controller
}

func (c *TableauControllerDL) GetDataLabels() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	type SiteRequest struct {
		SiteID string `json:"siteID"`
	}
	var req SiteRequest

	// DON'T REPLACE Ctx.Input.CopyBody WITH Ctx.Input.RequestBody
	// OTHERWISE ERROR "unexpected end of json input"
	// ie json will be empty

	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &req)

	// check JSON format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau api
	response, err := utils.Tableau_get_data_label_values(models.Get_token(), req.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// read response body
	bodyread, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	labels, _ := utils.ExtractLabelNames(string(bodyread))

	// return response
	c.Data["json"] = map[string]interface{}{"datalabels": labels}
	c.ServeJSON()

}
