package controllers

import (
	"net/http"
	"strings"
)

func (c *TableauController) GetAttribute() {
	// converted to lowercase in case something like "dataLabels" is entered
	param := strings.ToLower(c.Ctx.Input.Param(":param"))

	switch param {
	case "datasources":
		c.GetDataSources()

	case "datalabels":
		c.GetDataLabels()

	case "projects":
		c.GetProjects()

	default:
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = map[string]string{"Error": "Invalid JSON format"}
		c.ServeJSON()
	}
}
