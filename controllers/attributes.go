package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func (c *TableauController) GetAttribute() {
	c.EnableRender = false
	param := strings.ToLower(c.Ctx.Input.Param(":param"))

	var requestBody models.SiteRequest
	var attributes []map[string]interface{}

	// dont use c.bindjson
	if err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format in request"}
		c.ServeJSON()
		return
	}

	// using tableau REST API
	response, err := lib.TableauGetAttribute(param, requestBody.SiteID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// extrect details out of response from Tableau
	attributes, err = utils.ExtractAttributes(response, param)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to extract attribute from Tableau response"}
		c.ServeJSON()
		return
	}

	// return JSON
	c.Data["json"] = map[string]interface{}{
		param: attributes,
	}
	c.ServeJSON()
}
