package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

const maxRetries = 5

func (c *TableauController) PostSync() {
	c.EnableRender = false

	// parse request to struct
	var requestBody models.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// get metadata from csv
	datasourceRecord, err := utils.ParseCSV(requestBody.Filename)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "error fetching metadata from CSV: "+err.Error())
		return
	}

	// add more tables, cols (assets)
	if requestBody.CreateNewAssets {

		// get all current data sources
		var currentDatasources []map[string]interface{}

		call := func() error {
			cd, err := lib.TableauGetAttributes("datasources")
			currentDatasources = cd
			if err != nil {
				return err
			}
			return nil
		}

		if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
			HandleError(c, http.StatusInternalServerError, "error fetching current data sources: "+err.Error())
			return
		}

		// if ds already exists we merge existing assets in tds
		// to prevent overwriting
		// see documentation for detail

		if datasourceID, exists := utils.DatasourceExists(currentDatasources, datasourceRecord.Datasource); exists {

			// download existing data source
			fileName := ""
			call := func() error {
				f, err := lib.TableauDownloadDataSource(datasourceID)
				fileName = f
				if err != nil {
					return err
				}
				return nil
			}

			if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
				HandleError(c, http.StatusInternalServerError, "error downloading tds file: "+err.Error())
				return
			}

			// extract metadata of existing data source
			datasourceRecord, err = utils.ExtractAssets(fileName, datasourceRecord)
			if err != nil {
				HandleError(c, http.StatusInternalServerError, "error extracting current assets: "+err.Error())
				return
			}

		}

		fileNameTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, models.GetSiteID())

		// generate tds file w this name
		if err := utils.GenerateTDSFile(fileNameTDS, datasourceRecord); err != nil {
			HandleError(c, http.StatusInternalServerError, "error generating tds file from struct: "+err.Error())
			return
		}

		// publish this file to this project
		call = func() error {
			if _, err := lib.TableauPublishDatasource(fileNameTDS, datasourceRecord.Datasource, requestBody.ProjectId); err != nil {
				return err
			}
			return nil
		}

		if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
			HandleError(c, http.StatusInternalServerError, "error publishing data source: "+err.Error())
			return
		}

	}

	if requestBody.LabelAssets {

		call := func() error {
			if err := lib.TableauCreateCategory(requestBody.AttributeMap.ContentProfile); err != nil {
				return err
			}
			if err := lib.TableauCreateCategory(requestBody.AttributeMap.DataElements); err != nil {
				return err
			}
			return nil
		}

		if err := CallWithRetry(c.Ctx.Request.Context(), call); err != nil {
			log.Println("category already exists")
		}

		// make channels
		jobs := make(chan models.WorkerLabelInfo, 20)
		results := make(chan error, 20)
		var wg sync.WaitGroup

		// make workers
		for id := 1; id <= runtime.NumCPU(); id++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker(c.Ctx.Request.Context(), jobs, results)
			}()
		}

		// send jobs in channels
		for _, table := range datasourceRecord.Tables {
			info := models.WorkerLabelInfo{
				TableInfo:      table,
				DatabaseName:   datasourceRecord.Database,
				TableCategory:  requestBody.AttributeMap.ContentProfile,
				ColumnCategory: requestBody.AttributeMap.DataElements,
			}
			jobs <- info
		}
		// no more jobs
		close(jobs)

		// error handling
		var errorMsg error = nil
		for i := 1; i <= len(datasourceRecord.Tables); i++ {
			errorMsg = <-results
			if errorMsg != nil {
				log.Println(errorMsg.Error())
			}
		}

		wg.Wait()

		if errorMsg != nil {
			HandleError(c, http.StatusInternalServerError, "error labelling assets: "+errorMsg.Error())
		}
	}

	// success message

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "records sync-ed successfully"}
	c.ServeJSON()

}

func worker(c context.Context, labelInfo <-chan models.WorkerLabelInfo, results chan<- error) {

	for info := range labelInfo {

		var tableID string
		var columnIDs map[string]string

		call := func() error {
			t, c, err := lib.TableauGetAssetIDs(info.DatabaseName, info.TableInfo.TableName)
			tableID = t
			columnIDs = c
			if err != nil {
				return err
			}

			return nil
		}
		if err := CallWithRetry(c, call); err != nil {
			results <- err
		}

		// label table
		if info.TableInfo.ContentProfiles != "" && info.TableCategory != "" {
			call := func() error {
				if err := lib.TableauLabelAsset(info.TableInfo.ContentProfiles, info.TableCategory, "table", tableID); err != nil {
					results <- err
				}
				return nil
			}
			if err := CallWithRetry(c, call); err != nil {
				results <- err
			}
		}

		// loop through all columns of table
		for _, column := range info.TableInfo.Columns {

			if column.DataElements != "" && info.ColumnCategory != "" {

				// label column
				columnNameTableau := fmt.Sprintf("%s.%s", info.TableInfo.TableName, column.ColumnName)
				columnID := columnIDs[columnNameTableau]

				call := func() error {
					if err := lib.TableauLabelAsset(column.DataElements, info.ColumnCategory, "column", columnID); err != nil {
						results <- err
					}
					return nil
				}
				if err := CallWithRetry(c, call); err != nil {
					results <- err
				}

			}

		}
		results <- nil
	}

}
