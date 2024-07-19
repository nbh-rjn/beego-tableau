package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io"
	"net/http"
)

func (c *TableauController) GetProjects() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	var request SiteRequest

	// dont use Ctx.Input.RequestBody

	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request)

	// check JSON format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau api
	response, err := utils.TableauGetProjects(models.Get_token(), request.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// read response body
	responseBody, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	projectNames, projectIDs, _ := utils.ExtractProjectNames(string(responseBody))

	var projects []map[string]interface{}
	for i := 0; i < len(projectNames); i++ {
		projects = append(projects, map[string]interface{}{
			"Name": projectNames[i],
			"ID":   projectIDs[i],
		})
	}

	// Set response data
	c.Data["json"] = map[string]interface{}{
		"Projects": projects,
	}

	// Serve JSON response
	c.ServeJSON()

}
