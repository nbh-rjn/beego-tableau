package controllers

import (
	"beego-project/utils"
	_ "encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "github.com/nbh-rjn"
	c.Data["Email"] = "nabiha.rajani@foundri.net"
	c.TplName = "index.tpl"
}

type TableauController struct {
	beego.Controller
}

func (c *TableauController) PostAuth() {

	// this is imp otherwise it looks for tpl file to render and we get error
	c.EnableRender = false

	// json to struct
	type creds struct {
		PersonalAccessTokenName   string `json:"personalAccessTokenName"`
		PersonalAccessTokenSecret string `json:"personalAccessTokenSecret"`
		ContentUrl                string `json:"contentUrl"`
	}
	var reqBody creds

	// handle the errors
	if err := c.BindJSON(&reqBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// get the creds from req body
	personalAccessTokenName := reqBody.PersonalAccessTokenName
	personalAccessTokenSecret := reqBody.PersonalAccessTokenSecret
	contentUrl := reqBody.ContentUrl

	//  make xml request body
	xmlData := utils.Make_xml(personalAccessTokenName, personalAccessTokenSecret, contentUrl)

	// send req to tableau api, recieve response
	resp, u_err := utils.Tableau_req(xmlData)

	if u_err != nil {
		fmt.Println("Error creating request:", u_err)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	defer resp.Body.Close()

	// check for error in response
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed. Status code:", resp.StatusCode)
		c.Data["json"] = map[string]string{"error": "Request failed"}
		c.ServeJSON()
		return
	}

	// Read response body
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		c.Data["json"] = map[string]string{"error": "Failed to read response"}
		c.ServeJSON()
		return
	}

	// extract token from response
	cred_token := utils.Get_token(string(responseBody))

	// Return response data
	c.Data["json"] = map[string]interface{}{"credential_token": cred_token}
	c.ServeJSON()
}
