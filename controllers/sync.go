package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"net/http"
)

func (c *TableauController) PostSync() {
	// make catgeories here from payload
	c.EnableRender = false

	// parse request to struct
	var requestBody models.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	datasourceRecords, err := utils.ParseCSV(requestBody.Filename)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if requestBody.CreateNewAssets {
		// ** handle partial data
		// ** download existing data
		// ** merge both to handle overwrite, then publish
		// **** check for existing API
		if _, err := utils.CreateDatasources(datasourceRecords, requestBody.SiteID, requestBody.ProjectID); err != nil {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	// assets not being recognized in dvdrentals
	if err := utils.LabelAssets(datasourceRecords, requestBody.SiteID, requestBody.AttributeMap.ContentProfile, requestBody.AttributeMap.DataElements); err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "Records sync-ed successfully"}
	c.ServeJSON()

}
