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
				fmt.Println("error fetching data sources . . . ")
				continue
			}

			// check if data source already exists
			// if it exists we need to recreate to prevent loss when overwriting

			if id := utils.DatasourceExists(currentDatasources, datasourceRecord.Datasource); id != "" {
				// download existing data source
				if err := lib.TableauDownloadDataSource(id); err != nil {
					errorMsg = err
					continue
				}

				// extract metadata of existing data source
				datasourceRecord, err = utils.ExtractAssets(datasourceRecord)
				if err != nil {
					errorMsg = err
					continue
				}

			}

			// generate tds file for each datasource struct
			if err := utils.GenerateTDSFile(fileNameTDS, datasourceRecord); err != nil {
				errorMsg = err
				continue
			}

			// publish it
			if _, err := lib.PublishDatasource(fileNameTDS, datasourceRecord.Datasource, requestBody.ProjectID); err != nil {
				errorMsg = err
				continue
			}

		}

		if errorMsg != nil {
			HandleError(c, http.StatusInternalServerError, fmt.Sprintf("error creating assets: %s", errorMsg.Error()))
		}
	}

	// ** assets not being recognized in dvdrentals
	/*
		if err := utils.LabelAssets(datasourceRecords, requestBody.SiteID, requestBody.AttributeMap.ContentProfile, requestBody.AttributeMap.DataElements); err != nil {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
	*/

	errorMsg = nil
	lib.CreateCategory(requestBody.AttributeMap.ContentProfile)
	lib.CreateCategory(requestBody.AttributeMap.DataElements)

	for _, datasourceRecord := range datasourceRecords {

		// for table in tables
		for _, table := range datasourceRecord.Tables {
			tableID, columnIDs, err := lib.GetColumnIDs(datasourceRecord.Database, table.TableName)
			fmt.Println(table.TableName)
			if err != nil {
				fmt.Println("error fetching ids . . .")
				errorMsg = err
				continue
			}

			fmt.Println(tableID, columnIDs)

			if err := lib.TableauLabelAsset(table.ContentProfiles, requestBody.AttributeMap.ContentProfile, "table", tableID); err != nil {
				fmt.Println("error labelling table . . . ")
				errorMsg = err
				continue
			}

			for _, column := range table.Columns {

				columnID := columnIDs[column.ColumnName]

				if err := lib.TableauLabelAsset(column.DataElements, requestBody.AttributeMap.DataElements, "column", columnID); err != nil {
					fmt.Println("error labeling column ", column.ColumnName)
					errorMsg = err
					continue
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
