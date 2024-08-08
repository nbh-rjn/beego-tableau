package controllers

import (
	"beego-project/lib"
	"beego-project/logger"
	"beego-project/models"

	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauController struct {
	beego.Controller
	Logger *logger.ZapLogger
}

func (c *TableauController) PostAuth() {

	c.EnableRender = false // no .tpl to render

	// read request body
	var requestBody models.CredentialStruct
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	c.Logger.Info("auth request body parsed successfully")

	// get token from tableau api
	credentialsToken := ""
	siteID := ""

	call := func() error {
		ct, sid, err := lib.TableauAuthentication(requestBody.PersonalAccessTokenName, requestBody.PersonalAccessTokenSecret, requestBody.ContentUrl)
		credentialsToken = ct
		siteID = sid
		if err != nil {
			return err
		}
		return nil
	}

	if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	// save for future use
	models.SaveCredentials(credentialsToken, siteID)
	models.SaveCredentialsDB(requestBody.PersonalAccessTokenName, requestBody.PersonalAccessTokenSecret, siteID)

	// Return response
	c.Data["json"] = map[string]interface{}{
		"credentialsToken": credentialsToken, "siteID": siteID,
	}
	c.ServeJSON()

	c.Logger.Info("authentication successful")

}
