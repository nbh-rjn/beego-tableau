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
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// using tableau REST API
	response, err := lib.TableauGetAttribute(param, requestBody.SiteID)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data sources from Tableau")
	}

	// extrect details out of response from Tableau
	attributes, err = utils.ExtractAttributes(response, param)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to extract attribute from Tableau response")
	}

	// save attributes to database
	models.SaveAttributesDB(param, requestBody.SiteID, attributes)

	c.Data["json"] = map[string]interface{}{
		param: attributes,
	}
	c.ServeJSON()
}
