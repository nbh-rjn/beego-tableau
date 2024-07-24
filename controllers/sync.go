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
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// synchronize records
	if requestBody.CreateNewAssets {
		if err := utils.TableauCreateDatasources(requestBody.Filename, requestBody.SiteID, requestBody.CreateNewAssets); err != nil {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := utils.UpdateDataLabels(requestBody.Filename, requestBody.SiteID, requestBody.CreateNewAssets); err != nil {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"Success": "Records sync-ed successfully"}
	c.ServeJSON()

}
