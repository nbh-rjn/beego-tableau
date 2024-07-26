package controllers

import (
	"beego-project/lib"
	"beego-project/models"

	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauController struct {
	beego.Controller
}

func (c *TableauController) PostAuth() {
	// no .tpl to render
	c.EnableRender = false

	// read request body
	var requestBody models.CredentialStruct
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// get token from tableau api
	credentialsToken, siteID, err := lib.TableauAuthRequest(requestBody.PersonalAccessTokenName, requestBody.PersonalAccessTokenSecret, requestBody.ContentUrl)

	if err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
	}

	models.SaveCredentials(credentialsToken, siteID)
	models.SaveCredentialsDB(requestBody.PersonalAccessTokenName, requestBody.PersonalAccessTokenSecret, siteID)

	// Return response data
	c.Data["json"] = map[string]interface{}{
		"credentialsToken": credentialsToken, "siteID": siteID,
	}
	c.ServeJSON()

}
