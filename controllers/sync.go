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
	var errorMsg error = nil

	// parse request to struct
	var requestBody models.SyncRequest
	if err := c.BindJSON(&requestBody); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid JSON format in request")
	}

	// get metadata from csv
	datasourceRecords, err := utils.ParseCSV(requestBody.Filename)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// add more tables, cols (assets)
	if requestBody.CreateNewAssets {
		for _, datasourceRecord := range datasourceRecords {

			// tds filename
			fileNameTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, models.Get_siteID())
			filePath := "storage/" + fileNameTDS

			// get all current data sources
			currentDatasources, err := lib.TableauGetAttributes("datasources")
			if err != nil {
				errorMsg = err
				continue
			}

			// check if data source already exists
			// if it exists we need to recreate to prevent loss when overwriting

			if id := utils.DatasourceExists(currentDatasources, datasourceRecord.Datasource); id != "" {

				// download existing data source
				filePath, err := lib.TableauDownloadDataSource(id)
				if err != nil {
					errorMsg = err

					// we dont continue with the sync if the data source exists but we cant download it
					// otherwise it will get overwritten
					continue
				}

				// extract metadata of existing data source
				// dvd rentals columns issue
				datasourceRecord, err = utils.ExtractAssets(filePath, datasourceRecord)
				if err != nil {
					errorMsg = err
				}

			}

			// generate tds file for each datasource struct
			if err := utils.GenerateTDSFile(filePath, datasourceRecord); err != nil {
				errorMsg = err
			}

			// publish it
			if _, err := lib.TableauPublishDatasource(filePath, fileNameTDS, datasourceRecord.Datasource, requestBody.ProjectID); err != nil {
				errorMsg = err
			}

		}

	}

	if errorMsg != nil {
		HandleError(c, http.StatusInternalServerError, fmt.Sprintf("error creating assets: %s", errorMsg.Error()))
	}

	// cant label columns in dvdrentals
	// cols not being recognized in dvdrentals

	errorMsg = nil

	// create label categories
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

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "Records sync-ed successfully"}
	c.ServeJSON()

}

func worker(id int, labelInfo <-chan models.WorkerLabelInfo, results chan<- error) {
	//var info models.WorkerLabelInfo = <-labelInfo

	for info := range labelInfo {
		tableID, columnIDs, err := lib.TableauGetAssetIDs(info.DatabaseName, info.TableInfo.TableName)
		if err != nil {
			//results <- err
			//return
		}

		// label table
		if info.TableInfo.ContentProfiles != "" && info.TableCategory != "" {
			if err := lib.TableauLabelAsset(info.TableInfo.ContentProfiles, info.TableCategory, "table", tableID); err != nil {
				results <- err
				return
			}
		}

		// loop through all columns of table
		for _, column := range info.TableInfo.Columns {

			if column.DataElements != "" && info.ColumnCategory != "" { // label column
				columnID := columnIDs[column.ColumnName]
				if err := lib.TableauLabelAsset(column.DataElements, info.ColumnCategory, "column", columnID); err != nil {
					//results <- err
					//return
				}
			}

		}

		log.Printf("worker %d completed %s\n", id, info.TableInfo.TableName)
		results <- nil
	}
	//return

}
