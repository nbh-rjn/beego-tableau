package controllers

import (
	"beego-project/utils"
	"encoding/json"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauControllerCSV struct {
	beego.Controller
}

func (c *TableauControllerCSV) PostSync() {
	c.EnableRender = false

	// structs for unmarshaling the request

	type AttributeMap struct {
		DataElements   string `json:"data_elements"`
		ContentProfile string `json:"content_profile"`
	}

	type InstanceMap map[string]string

	type JsonRequest struct {
		Filename        string       `json:"filename"`
		SiteID          string       `json:"siteID"`
		CreateNewAssets bool         `json:"create_new_assets"`
		EntityType      string       `json:"entity_type"`
		AttributeMap    AttributeMap `json:"attribute_map"`
		InstanceMap     InstanceMap  `json:"instance_map"`
	}

	// parse request to struct
	var req JsonRequest
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &req)

	// check for errors in JSON format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function takes file URL and returns parsed struct array
	records := utils.ParseCSV(req.Filename)
	utils.Tableau_sync(records, req.SiteID)

	//utils.Tableau_sync(records, req.SiteID)

}
