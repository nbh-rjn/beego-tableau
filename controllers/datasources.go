package controllers

import (
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"io"
	"net/http"
)

func (c *TableauController) GetDataSources() {
	// so it doesnt go looking in views for a tpl to render
	c.EnableRender = false

	var request SiteRequest

	// dont use Ctx.Input.RequestBody
	err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request)

	// check for correct request format
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// utility function to communicate with Tableau API
	response, err := utils.TableauGetDataSources(models.Get_token(), request.SiteID)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// read body of response
	responseBody, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	datasourceNames, datasourceIDs := utils.ExtractDataSources(string(responseBody))

	var datasources []map[string]interface{}
	for i := 0; i < len(datasourceNames); i++ {
		datasources = append(datasources, map[string]interface{}{
			"Name": datasourceNames[i],
			"ID":   datasourceIDs[i],
		})
	}

	// Set response data
	c.Data["json"] = map[string]interface{}{
		"Data sources": datasources,
	}

	// Serve JSON response
	c.ServeJSON()

	// return info in response
	c.Data["json"] = map[string]interface{}{"Data sources": datasources}
	c.ServeJSON()

}
