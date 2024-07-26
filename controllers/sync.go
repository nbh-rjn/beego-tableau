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

	datasourceRecords := utils.ParseCSV(requestBody.Filename)
	if datasourceRecords == nil {
		HandleError(c, http.StatusInternalServerError, "Could not parse raw CSV file")
	}

	if requestBody.CreateNewAssets {
		// ** handle partial data
		// ** download existing data
		// ** merge both to handle overwrite, then publish
		// **** check for existing API
		if err := utils.TableauCreateDatasources(datasourceRecords, requestBody.SiteID, requestBody.ProjectID); err != nil {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
	}

	// assets not being recognized

	if err := utils.LabelAssets(datasourceRecords, requestBody.SiteID, requestBody.AttributeMap.ContentProfile, requestBody.AttributeMap.DataElements); err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "Records sync-ed successfully"}
	c.ServeJSON()

}
