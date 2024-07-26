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

	//fmt.Println(lib.PublishDatasource("download.tds", "2ff64d57-b7c1-4e99-803c-13bb81ae0371", "testing", "f7eea7f7-2c14-4694-a6ac-27af2e0bc583"))

	//lib.TableauLabelAsset("testlabel3", "testcategory", "column", "bcb81941-6bb7-4054-a323-9aa82fc7d51e")

	//lib.ApplyLabelValue(models.Get_siteID(), "table", id, "label1")

	//lib.TableauLabelAsset("label2", "testcategory", "table", "2d418eed-ff45-4263-b146-c7010f69d938")
	lib.CreateCategory("testcategory")
	// Return response data
	c.Data["json"] = map[string]interface{}{
		"credentialsToken": credentialsToken, "siteID": siteID,
	}
	c.ServeJSON()

}
