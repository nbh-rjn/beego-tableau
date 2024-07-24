package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"fmt"

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

	//  construct XML request to tableau
	xmlData := utils.CredentialsXML(
		requestBody.PersonalAccessTokenName,
		requestBody.PersonalAccessTokenSecret,
		requestBody.ContentUrl,
	)

	// send request to tableau api
	response, err := lib.TableauAuthRequest(xmlData)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
	}
	defer response.Body.Close()

	// extract token from response
	credentialsToken, err := utils.ExtractToken(response)
	if err != nil {
		HandleError(c, http.StatusServiceUnavailable, "Failed to extract credentials from Tableau response")
	}

	// save session token
	models.SaveToken(credentialsToken)
	models.SaveCredentialsDB(requestBody)

	fmt.Println(lib.GetAssetID("table", "AdventureWorks", "SalesOrderDetail", ""))
	//fmt.Println(lib.GetAssetID("column", "AdventureWorks", "SalesOrderDetail", "ProductID"))
	fmt.Println(lib.GetColumns("AdventureWorks", "SalesOrderDetail"))

	// Return response data
	c.Data["json"] = map[string]interface{}{"credentialsToken": credentialsToken}
	c.ServeJSON()

}
