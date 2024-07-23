package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"net/http"
)

func (c *TableauController) PostSync() {
	c.EnableRender = false

	// parse request to struct
	var requestBody models.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format in request"}
		c.ServeJSON()
		return
	}

	// synchronize records
	if requestBody.CreateNewAssets {
		if err := utils.TableauCreateDatasources(requestBody.Filename, requestBody.SiteID, requestBody.CreateNewAssets); err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": err.Error()}
			c.ServeJSON()
			return
		}
	} else {
		if err := utils.UpdateDataLabels(requestBody.Filename, requestBody.SiteID, requestBody.CreateNewAssets); err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": err.Error()}
			c.ServeJSON()
			return
		}
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"Success": "Records sync-ed successfully"}
	c.ServeJSON()

}
