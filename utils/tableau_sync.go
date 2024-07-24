package utils

import (
	"beego-project/lib"
	"fmt"

	"github.com/pkg/errors"
)

func TableauCreateDatasources(filenameCSV string, siteID string, createAssets bool) error {

	// parse CSV to slice of struct
	datasourceRecords := ParseCSV(filenameCSV)
	if datasourceRecords == nil {
		return errors.New("Could not parse raw CSV file")
	}

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
		if err := lib.PublishDatasource(fileNameTDS, siteID, datasourceRecord.Datasource); err != nil {
			return err
		}

	}
	return nil
}

func UpdateDataLabels(filenameCSV string, siteID string, createAssets bool) error {

	// parse CSV to slice of structs
	datasourceRecords := ParseCSV(filenameCSV)

	if datasourceRecords == nil {
		return errors.New("Could not parse raw CSV file")
	}

	for _, datasourceRecord := range datasourceRecords {

		// assuming a datasource == one db
		databaseName := datasourceRecord.Database

		// for table in tables
		for _, table := range datasourceRecord.Tables {

			//create category acc to datasourceRecord.table.contentprofile
			if table.ContentProfiles != "" {
				lib.CreateCategory(siteID, table.ContentProfiles)
			}

			// use this if labels need to be applied on TABLES

			// fetch table id from graphql
			tableID, err := lib.GetTableID(databaseName, table.TableName)
			if err != nil {
				return err
			}

			// create label value to apply on table ? if needed
			lib.CreateLabelValue(siteID, table.ContentProfiles, table.ContentProfiles)
			lib.ApplyLabelValue(siteID, "table", tableID, table.ContentProfiles)

			// for column in column
			columnIDs, err := lib.GetColumns(databaseName, table.TableName)
			if err != nil {
				return err
			}
			for _, column := range table.Columns {

				// fetch column id from graphql
				columnID := columnIDs[column.ColumnName]

				// create labelvalue acc to datasourceRecord.table.column.DataElements
				if err := lib.CreateLabelValue(siteID, column.DataElements, table.ContentProfiles); err != nil {
					return err
				}

				//apply on columns
				if err := lib.ApplyLabelValue(siteID, "column", columnID, column.DataElements); err != nil {
					return err
				}
			}

		}

	}
	return nil

}
