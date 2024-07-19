package controllers

import (
	"beego-project/utils"
	"encoding/json"
	"net/http"
)

type AttributeMap struct {
	DataElements   string `json:"data_elements"`
	ContentProfile string `json:"content_profile"`
}

type InstanceMap map[string]string

type SyncRequest struct {
	Filename        string       `json:"filename"`
	SiteID          string       `json:"siteID"`
	CreateNewAssets bool         `json:"create_new_assets"`
	EntityType      string       `json:"entity_type"`
	AttributeMap    AttributeMap `json:"attribute_map"`
	InstanceMap     InstanceMap  `json:"instance_map"`
}

func (c *TableauController) PostSync() {
	c.EnableRender = false

	// parse request to struct
	var request SyncRequest
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request)

	// check for errors in JSON format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function takes file URL and returns parsed struct array
	records := utils.ParseCSV(request.Filename)

	if records == nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Could not parse raw CSV file"}
		c.ServeJSON()
		return
	}

	// synchronize records
	err = utils.TableauSyncRecords(records, request.SiteID)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Could not sync data sources to Tableau"}
		c.ServeJSON()
		return
	}

}
