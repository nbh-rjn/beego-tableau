package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	//"fmt"
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

	// IMPORTANT!!!!!!
	// THIS Ctx.Input.CopyBody FUNCTION IS IMPORTANT
	// DON'T REPLACE IT WITH Ctx.Input.RequestBody
	// OTHERWISE IT WILL KEEP SHOWING EMPTY
	// ie "unexpected end of json input"

	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &req)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	response, err := utils.Tableau_get_data_label_values(models.Get_token(), req.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	bodyread, _ := ioutil.ReadAll(response.Body)
	labels, _ := utils.ExtractLabelNames(string(bodyread))
	c.Data["json"] = map[string]interface{}{"datalabels": labels}
	c.ServeJSON()
	/*
		contentType := response.Header.Get("Content-Type")
		c.Ctx.Output.ContentType(contentType)
		c.Ctx.Output.Body(bodyread)
	*/

}
