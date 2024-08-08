package utils

import (
	"beego-project/lib"
	"beego-project/models"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

func ExtractAssets(fileName string, dsCSV models.DatasourceStruct) (models.DatasourceStruct, error) {

	// read existing assets from datasource
	storage := models.GetStorage(context.TODO())

	data, err := storage.Read(context.TODO(), fileName)
	if err != nil {
		return dsCSV, err
	}

	var dsTDS models.DatasourceGeneration
	err = xml.Unmarshal(data, &dsTDS)
	if err != nil {
		return dsCSV, err
	}

	// create map of all current assets
	outerMap := make(map[string]map[string]string)

	for _, tableTDS := range dsTDS.Connection.Relations {
		outerMap[tableTDS.Name] = make(map[string]string)
	}

	// the inner maps have column information
	for _, columnTDS := range dsTDS.Connection.MetadataRecords.Records {
		columnParent := strings.Trim(columnTDS.ParentName, "[]")

		if innerMap, ok := outerMap[columnParent]; ok {
			innerMap[columnTDS.RemoteName] = columnTDS.LocalType
		}
	}

	for _, table := range dsCSV.Tables {

		// if csv table aready exists in published data source
		if columnsMap, exists := outerMap[table.TableName]; exists {

			// delete the csv columns from map of existing columns
			for _, column := range table.Columns {
				delete(columnsMap, column.ColumnName)
			}

			// fetch ids of columns within it
			// we need col ids to retrieve existing labels
			tableID, columnIDs, err := lib.TableauGetAssetIDs(dsCSV.Database, table.TableName)
			if err != nil {
				log.Printf("error in extracting assets: could not retrieve asset IDs...")
			}

			// remaining columns in map will be those that are not already published
			for colName, dataType := range columnsMap {

				colNameTableau := fmt.Sprintf("%s.%s", table.TableName, colName)
				columnLabel, err := lib.TableauGetAssetLabel("column", columnIDs[colNameTableau])
				if err != nil {
					columnLabel = ""
				}

				table.Columns = append(table.Columns, models.ColumnStruct{
					ColumnName:        colName,
					ColumnType:        dataType,
					ColumnDescription: "",
					DataElements:      columnLabel,
				})
			}

			tableLabel, err := lib.TableauGetAssetLabel("table", tableID)
			if err != nil {
				tableLabel = ""
			}
			if table.ContentProfiles == "" && tableLabel != "" {
				table.ContentProfiles = tableLabel
			}

			delete(outerMap, table.TableName)
		}
	}

	// Append remaining tables from the map to the slice
	for tableName, columnsMap := range outerMap {
		tableID, columnIDs, err := lib.TableauGetAssetIDs(dsCSV.Database, tableName)
		if err != nil {
			log.Printf("error in extracting assets: could not retrieve asset IDs...")
		}

		var columns []models.ColumnStruct
		for columnName, dataType := range columnsMap {
			columnLabel, err := lib.TableauGetAssetLabel("column", columnIDs[columnName])
			if err != nil {
				columnLabel = ""
			}

			columns = append(columns, models.ColumnStruct{
				ColumnName:        columnName,
				ColumnType:        dataType,
				ColumnDescription: "",
				DataElements:      columnLabel,
			})
		}

		tableLabel, err := lib.TableauGetAssetLabel("table", tableID)
		if err != nil {
			tableLabel = ""
		}
		dsCSV.Tables = append(dsCSV.Tables, models.TableStruct{
			Id:              "",
			TableName:       tableName,
			TableType:       "",
			ContentProfiles: tableLabel,
			Columns:         columns,
		})
	}

	return dsCSV, nil

}
