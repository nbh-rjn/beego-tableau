package controllers

import (
	"beego-project/lib"
	"net/http"
	"strings"
)

func (c *TableauController) GetAttribute() {
	c.EnableRender = false
	param := strings.ToLower(c.Ctx.Input.Param(":param"))

	// using tableau REST API
	attributes, err := lib.TableauGetAttributes(param)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	c.Data["json"] = map[string]interface{}{
		param: attributes,
	}
	c.ServeJSON()
}
