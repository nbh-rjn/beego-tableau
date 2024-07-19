package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"

	//"fmt"

	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauController struct {
	beego.Controller
}

func (c *TableauController) PostAuth() {
	// no .tpl to render
	c.EnableRender = false

	var requestBody models.CredentialStruct

	// read auth request
	if err := c.BindJSON(&requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format in request"}
		c.ServeJSON()
		return
	}

	//  make XML for request to tableau
	xmlData := utils.CredentialsXML(
		requestBody.PersonalAccessTokenName,
		requestBody.PersonalAccessTokenSecret,
		requestBody.ContentUrl,
	)

	// send request to tableau api
	response, err := lib.TableauAuthRequest(xmlData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	defer response.Body.Close()

	//  error in the response from tableau
	if response.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(http.StatusServiceUnavailable)
		c.Data["json"] = map[string]string{"error": "Tableau response error"}
		c.ServeJSON()
		return
	}

	// extract token from response
	credentialsToken, err := utils.ExtractToken(response)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to extract credentials from Tableau response"}
		c.ServeJSON()
		return
	}

	// save session token
	models.SaveToken(credentialsToken)

	// Return response data
	c.Data["json"] = map[string]interface{}{"credentialsToken": credentialsToken}
	c.ServeJSON()
}
