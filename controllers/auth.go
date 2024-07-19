package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"

	//"fmt"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type TableauController struct {
	beego.Controller
}

type credentialStruct struct {
	PersonalAccessTokenName   string `json:"personalAccessTokenName"`
	PersonalAccessTokenSecret string `json:"personalAccessTokenSecret"`
	ContentUrl                string `json:"contentUrl"`
}

func (c *TableauController) PostAuth() {

	// no .tpl to render
	c.EnableRender = false

	var requestBody credentialStruct

	// read auth request and handle error
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

	// send request to tableau api, recieve response
	response, err := lib.TableauAuthRequest(xmlData)

	// error in creating our request
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{
			"error":  "Failed to create authentication request to Tableau",
			"detail": err.Error(),
		}
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

	// Read response body
	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read response body"}
		c.ServeJSON()
		return
	}

	// extract token from response
	credentialsToken, err := utils.ExtractToken(string(responseBody))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to extract credentials from Tableau response"}
		c.ServeJSON()
		return
	}

	// save session token
	// so we dont have to keep including it in request body from other endpoints
	models.SaveToken(credentialsToken)

	// Return response data
	c.Data["json"] = map[string]interface{}{"Credentials Token": credentialsToken}
	c.ServeJSON()
}
