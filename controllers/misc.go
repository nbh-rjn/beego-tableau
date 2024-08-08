package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (c *TableauController) DownloadDataSource() {

	c.EnableRender = false
	var request models.DownloadRequest

	// dont use c.bindJSON
	if err := json.Unmarshal((c.Ctx.Input.CopyBody(1000)), &request); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
		return
	}
	c.Logger.Info("download request body parsed successfully")

	// make api request
	fileName := ""
	call := func() error {
		f, err := lib.TableauDownloadDataSource(request.DatasourceID)
		fileName = f
		if err != nil {
			return err
		}
		return nil
	}

	if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to fetch data sources from Tableau")
		return
	}

	c.Logger.Info("data source fetched from Tableau successfully")

	storage := models.GetStorage(context.TODO())
	if storage == nil {
		HandleError(c, http.StatusInternalServerError, "No storage handler found")
		return
	}

	// Read the file from storage
	fileData, err := storage.Read(context.TODO(), fileName)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to read file from storage")
		return
	}

	c.Logger.Info("downloaded file read from storage successfully")

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Body(fileData)

}

func (c *TableauController) GetAttribute() {
	c.EnableRender = false
	param := strings.ToLower(c.Ctx.Input.Param(":param"))

	// using tableau REST API
	var attributes []map[string]interface{}
	call := func() error {
		a, err := lib.TableauGetAttributes(param)
		attributes = a
		if err != nil {
			return err
		}
		return nil
	}

	if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Logger.Info("attributes fetched successfully", "param", param)

	c.Data["json"] = map[string]interface{}{
		param: attributes,
	}
	c.ServeJSON()
}

func HandleError(c *TableauController, status uint, errormsg string) {
	c.Logger.Error("errormsg")
	c.Ctx.Output.SetStatus(int(status))
	c.Data["json"] = map[string]string{"error": errormsg}
	c.ServeJSON()
}
