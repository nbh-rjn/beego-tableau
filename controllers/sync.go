package controllers

import (
	"beego-project/utils"
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
	var requestBody SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format in request"}
		c.ServeJSON()
		return
	}

	// synchronize records
	if err := utils.TableauSyncRecords(requestBody.Filename, requestBody.SiteID); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"Success": "Records sync-ed successfully"}
	c.ServeJSON()

}
