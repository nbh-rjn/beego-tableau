package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

func (c *TableauController) PostSync() {
	c.EnableRender = false

	// parse request to struct
	var requestBody models.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// get metadata from csv
	datasourceRecords, err := utils.ParseCSV(requestBody.Filename)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// add more tables, cols (assets)
	if requestBody.CreateNewAssets {
		var errorMsg error = nil

		for _, datasourceRecord := range datasourceRecords {

			// get all current data sources
			currentDatasources, err := lib.TableauGetAttributes("datasources")
			if err != nil {
				log.Println("error fetching current data sources", err)
				errorMsg = err
				continue
			}

			// if ds alr exists
			// we recreate, merge, publish

			if datasourceID, exists := utils.DatasourceExists(currentDatasources, datasourceRecord.Datasource); exists {
				log.Printf("data source %s exists  . . .\n", datasourceRecord.Datasource)

				// download existing data source
				filePathDownload, err := lib.TableauDownloadDataSource(datasourceID)
				if err != nil {
					log.Println("error downloading tds file", err)
					errorMsg = err

					continue // to prevent overwriting
				}

				// extract metadata of existing data source
				// dvd rentals columns issue
				datasourceRecord, err = utils.ExtractAssets(filePathDownload, datasourceRecord)
				if err != nil {
					log.Println("error extracting existing assets: ", err)
					errorMsg = err
				}

			}

			fileNameTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, models.Get_siteID())
			filePathPublish := "storage/" + fileNameTDS

			// generate tds file for each datasource struct
			if err := utils.GenerateTDSFile(filePathPublish, datasourceRecord); err != nil {
				log.Println("error generating tds file from struct:", err)
				errorMsg = err
			}

			// publish it
			if _, err := lib.TableauPublishDatasource(filePathPublish, fileNameTDS, datasourceRecord.Datasource, requestBody.ProjectID); err != nil {
				log.Println("error publishing data source:", err)
				errorMsg = err
			}

		}
		if errorMsg != nil {
			HandleError(c, http.StatusInternalServerError, fmt.Sprintf("error creating assets: %s", errorMsg.Error()))
		}
	}

	if requestBody.LabelAssets {
		var errorMsg error = nil

		lib.TableauCreateCategory(requestBody.AttributeMap.ContentProfile)
		lib.TableauCreateCategory(requestBody.AttributeMap.DataElements)

		for _, datasourceRecord := range datasourceRecords {

			// make channels
			jobs := make(chan models.WorkerLabelInfo, 20)
			results := make(chan error, 20)
			var wg sync.WaitGroup

			// make workers
			for id := 1; id <= runtime.NumCPU(); id++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					worker(id, jobs, results)
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

			for i := 1; i <= len(datasourceRecord.Tables); i++ {
				errorMsg = <-results
				if errorMsg != nil {
					log.Println(errorMsg)
				}
			}

			wg.Wait()

		}

		if errorMsg != nil {
			HandleError(c, http.StatusInternalServerError, fmt.Sprintf("error labelling assets: %s", errorMsg.Error()))
		}
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "Records sync-ed successfully"}
	c.ServeJSON()

}

func worker(id int, labelInfo <-chan models.WorkerLabelInfo, results chan<- error) {

	for info := range labelInfo {
		tableID, columnIDs, err := lib.TableauGetAssetIDs(info.DatabaseName, info.TableInfo.TableName)
		if err != nil {
			log.Println("worker", id, "error:", err)
			results <- err
			return
		}

		// label table
		if info.TableInfo.ContentProfiles != "" && info.TableCategory != "" {
			log.Println(info.TableInfo.TableName, tableID)
			if err := lib.TableauLabelAsset(info.TableInfo.ContentProfiles, info.TableCategory, "table", tableID); err != nil {
				results <- err
				return
			}
		}

		// loop through all columns of table
		for _, column := range info.TableInfo.Columns {

			if column.DataElements != "" && info.ColumnCategory != "" {
				// label column

				columnID := columnIDs[fmt.Sprintf("%s.%s", info.TableInfo.TableName, column.ColumnName)]

				if err := lib.TableauLabelAsset(column.DataElements, info.ColumnCategory, "column", columnID); err != nil {
					results <- err
					return
				}
			}

		}
		results <- nil
	}

}
