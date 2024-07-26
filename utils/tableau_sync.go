package utils

import (
	"beego-project/lib"
	"beego-project/models"
	"fmt"
)

// ** return slice of datasource ids
// rename to upload data sources??
func TableauCreateDatasources(datasourceRecords []models.DatasourceStruct, siteID string, projectID string) error {

	// one datasource at a time
	// can handle more than one datasource per CSV file
	for _, datasourceRecord := range datasourceRecords {

		// tds filename
		fileNameTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, siteID)

		// generate tds file for each datasource struct
		if err := GenerateTDSFile(fileNameTDS, datasourceRecords); err != nil {
			return err
		}

		// publish it
		if _, err := lib.PublishDatasource(fileNameTDS, siteID, datasourceRecord.Datasource, projectID); err != nil {
			return err
		}

	}
	return nil
}

// ** return [] label ids
func LabelAssets(datasourceRecords []models.DatasourceStruct, siteID string, tableCategory string, columnCategory string) error {

	for _, datasourceRecord := range datasourceRecords {

		// assuming a datasource == one db
		databaseName := datasourceRecord.Database

		// for table in tables
		for _, table := range datasourceRecord.Tables {

			// ** category as per sync req body
			// ** remove from here
			/*
				if table.ContentProfiles != "" {
					lib.CreateCategory(tableCategory)
				}
			*/

			// ** make one func to get table ID and all its columns together

			// fetch table id from graphql
			tableID, err := lib.GetTableID(databaseName, table.TableName)

			if err != nil {

				// ** flow must continue
				// ** continue
				// ** dont return err
				return err
			}

			// ** one func

			if err := lib.TableauLabelAsset(table.ContentProfiles, tableCategory, "table", tableID); err != nil {

				// ** flow must continue
				// ** continue
				// ** dont return err
				return err
			}
			//lib.ApplyLabelValue(siteID, "table", tableID, table.ContentProfiles)

			// for column in column
			columnIDs, err := lib.GetColumns(databaseName, table.TableName)
			if err != nil {
				return err
				// continue
				// dont return
			}
			for _, column := range table.Columns {

				// fetch column id from graphql
				columnID := columnIDs[column.ColumnName]

				// ** make one func for creating and applying labels

				// create labelvalue acc to datasourceRecord.table.column.DataElements
				if err := lib.TableauLabelAsset(column.DataElements, columnCategory, "column", columnID); err != nil {
					return err
				}

				//apply on columns
				/*
					if err := lib.ApplyLabelValue(siteID, "column", columnID, column.DataElements); err != nil {
						return err
					}
				*/
			}

		}

	}
	return nil

}
