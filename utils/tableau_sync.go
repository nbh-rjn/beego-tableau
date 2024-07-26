package utils

import (
	"beego-project/lib"
	"beego-project/models"
	"fmt"
)

// ** return slice of datasource ids
func CreateDatasources(datasourceRecords []models.DatasourceStruct, siteID string, projectID string) ([]string, error) {

	var datasourceIDs []string
	var errorMsg error = nil

	// can handle more than one datasource per CSV file
	for _, datasourceRecord := range datasourceRecords {

		// tds filename
		fileNameTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, siteID)

		// generate tds file for each datasource struct
		if err := GenerateTDSFile(fileNameTDS, datasourceRecords); err != nil {
			errorMsg = err
			continue
		}

		// publish it
		datasourceID, err := lib.PublishDatasource(fileNameTDS, siteID, datasourceRecord.Datasource, projectID)
		if err != nil {
			errorMsg = err
			continue
		}
		datasourceIDs = append(datasourceIDs, datasourceID)

	}
	return datasourceIDs, errorMsg
}

// ** return [] label ids
func LabelAssets(datasourceRecords []models.DatasourceStruct, siteID string, tableCategory string, columnCategory string) error {

	var errorMsg error = nil

	for _, datasourceRecord := range datasourceRecords {

		// assuming a datasource == one db
		databaseName := datasourceRecord.Database

		// for table in tables
		for _, table := range datasourceRecord.Tables {

			tableID, columnIDs, err := lib.GetColumnIDs(databaseName, table.TableName)
			if err != nil {
				errorMsg = err
				continue
			}

			if err := lib.TableauLabelAsset(table.ContentProfiles, tableCategory, "table", tableID); err != nil {
				errorMsg = err
				continue
			}

			for _, column := range table.Columns {

				columnID := columnIDs[column.ColumnName]

				if err := lib.TableauLabelAsset(column.DataElements, columnCategory, "column", columnID); err != nil {
					errorMsg = err
					continue
				}

			}

		}

	}

	return errorMsg

}
