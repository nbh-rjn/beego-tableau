package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/beego/beego/orm"
)

func (c *TableauController) GetAttribute() {
	c.EnableRender = false
	param := strings.ToLower(c.Ctx.Input.Param(":param"))

	var requestBody models.SiteRequest
	var attributes []map[string]interface{}

	// dont use c.bindjson
	if err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format in request"}
		c.ServeJSON()
		return
	}

	// using tableau REST API
	response, err := lib.TableauGetAttribute(param, requestBody.SiteID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch data sources from Tableau"}
		c.ServeJSON()
		return
	}

	// extrect details out of response from Tableau
	attributes, err = utils.ExtractAttributes(response, param)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to extract attribute from Tableau response"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	for _, attribute := range attributes {
		name := string(attribute["name"].(string))

		switch param {
		case "datalabels":
			label := models.LabelsTable{
				LabelName: name,
				SiteID:    requestBody.SiteID,
			}
			o.Insert(&label)
		case "datasources":
			datasource := models.DatasourcesTable{
				DatasourceName: name,
				SiteID:         requestBody.SiteID,
			}
			o.Insert(&datasource)
		case "projects":
			project := models.ProjectsTable{
				ProjectName: name,
				SiteID:      requestBody.SiteID,
			}
			o.Insert(&project)
		default:
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]string{"error": "Invalid attribute type"}
			c.ServeJSON()
			return
		}
	}

	// return JSON
	c.Data["json"] = map[string]interface{}{
		param: attributes,
	}
	c.ServeJSON()
}
