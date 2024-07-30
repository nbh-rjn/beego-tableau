package controllers

import (
	"beego-project/lib"
	"beego-project/models"
	"beego-project/utils"
	"fmt"
	"net/http"
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
			if err := utils.GenerateTDSFile(fileNameTDS, datasourceRecord); err != nil {
				errorMsg = err
			}

			// publish it
			if _, err := lib.PublishDatasource(fileNameTDS, datasourceRecord.Datasource, requestBody.ProjectID); err != nil {
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
	lib.CreateCategory(requestBody.AttributeMap.ContentProfile)
	lib.CreateCategory(requestBody.AttributeMap.DataElements)

	for _, datasourceRecord := range datasourceRecords {

		// loop through all tables
		for _, table := range datasourceRecord.Tables {

			// get ids of table and its columns
			tableID, columnIDs, err := lib.GetColumnIDs(datasourceRecord.Database, table.TableName)
			if err != nil {
				errorMsg = err
			}

			// label table
			if err := lib.TableauLabelAsset(table.ContentProfiles, requestBody.AttributeMap.ContentProfile, "table", tableID); err != nil {
				errorMsg = err
			}

			// loop through all columns of table
			for _, column := range table.Columns {

				// label column
				columnID := columnIDs[column.ColumnName]
				if err := lib.TableauLabelAsset(column.DataElements, requestBody.AttributeMap.DataElements, "column", columnID); err != nil {
					errorMsg = err
				}

			}

		}

	}

	if errorMsg != nil {
		HandleError(c, http.StatusInternalServerError, fmt.Sprintf("error labelling assets: %s", errorMsg.Error()))
	}

	// success message
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]string{"success": "Records sync-ed successfully"}
	c.ServeJSON()

}
