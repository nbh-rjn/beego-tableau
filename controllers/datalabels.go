package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io"
	"net/http"
)

type SiteRequest struct {
	SiteID string `json:"siteID"`
}

func (c *TableauController) GetDataLabels() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	var request SiteRequest

	// dont use Ctx.Input.RequestBody
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request)

	// check JSON format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// use utility function to communicate with Tableau api
	response, err := lib.TableauGetDataLabelValues(models.Get_token(), request.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// read response body
	responseBody, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	dataLabels, _ := utils.ExtractLabelNames(string(responseBody))

	// return response
	c.Data["json"] = map[string]interface{}{"Data labels": dataLabels}
	c.ServeJSON()

}
