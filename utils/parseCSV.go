package utils

import (
	"beego-project/models"
	"encoding/csv"
	"fmt"
	"os"
)

func ParseCSV(filename string) (models.DatasourceStruct, error) {

	var ds models.DatasourceStruct
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return ds, err
	}
	defer file.Close()

	// read file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return ds, err
	}

	// create empty array
	var fcrecords []models.CSVRecords

	for i, record := range records {
		if i == 0 {
			continue
		}
		// store record values in elements of struct
		fcrecord := models.CSVRecords{
			Id:                record[0],
			Datasource:        record[1],
			Host:              record[2],
			Port:              record[3],
			DatabaseType:      record[4],
			DBUsername:        record[5],
			Database:          record[6],
			Schema:            record[7],
			Table:             record[8],
			TableType:         record[9],
			ContentProfiles:   record[10],
			Column:            record[11],
			ColumnType:        record[12],
			ColumnDescription: record[13],
			DataElements:      record[14],
		}

		// append to array
		fcrecords = append(fcrecords, fcrecord)
	}

	if len(fcrecords) < 1 {
		return ds, fmt.Errorf("no records to parse in CSV")
	}

	// return parsed array
	return organizeRecords(fcrecords)
}

// takes fcrecords which is just lines of csv file
// returns hierarchically arranged slice of structs
// each struct representing one datasource
func organizeRecords(records []models.CSVRecords) (models.DatasourceStruct, error) {

	var datasource models.DatasourceStruct

	if len(records) < 1 {
		return datasource, fmt.Errorf("no records obtained from CSV to organize")
	}

	dsDetail := records[0]
	ds := models.DatasourceStruct{
		Datasource: dsDetail.Datasource,
		Host:       dsDetail.Host,
		Port:       dsDetail.Port,
		Database:   dsDetail.Database,
		Schema:     dsDetail.Schema,
		DBUsername: dsDetail.DBUsername,
		DBType:     dsDetail.DatabaseType,
	}

	lastTb := ""

	for _, record := range records {
		if record.Table != lastTb {
			ds.Tables = append(ds.Tables, models.TableStruct{
				Id:              record.Id,
				TableName:       record.Table,
				TableType:       record.TableType,
				ContentProfiles: record.ContentProfiles,
			})
			lastTb = record.Table
		}

		currentTableIdx := len(ds.Tables) - 1

		ds.Tables[currentTableIdx].Columns = append(ds.Tables[currentTableIdx].Columns, models.ColumnStruct{
			ColumnName:        record.Column,
			ColumnType:        record.ColumnType,
			ColumnDescription: record.ColumnDescription,
			DataElements:      record.DataElements,
		})

	}
	return ds, nil

	/*

		ds_idx, tb_idx := "", ""
		dsi, tbi := -1, -1

		for _, record := range records {

			if ds_idx != record.Datasource {
				ds := models.DatasourceStruct{
					Datasource: record.Datasource,
					Host:       record.Host,
					Port:       record.Port,
					Database:   record.Database,
					Schema:     record.Schema,
					DBUsername: record.DBUsername,
					DBType:     record.DatabaseType,
				}

				datasources = append(datasources, ds)
				ds_idx = record.Datasource
				dsi = dsi + 1

			}
			if tb_idx != record.Table {
				tb := models.TableStruct{
					Id:              record.Id,
					TableName:       record.Table,
					TableType:       record.TableType,
					ContentProfiles: record.ContentProfiles,
				}

				datasources[dsi].Tables = append(datasources[dsi].Tables, tb)
				tb_idx = record.Table
				tbi = tbi + 1
			}
			col := models.ColumnStruct{
				ColumnName:        record.Column,
				ColumnType:        record.ColumnType,
				ColumnDescription: record.ColumnDescription,
				DataElements:      record.DataElements,
			}
			datasources[dsi].Tables[tbi].Columns = append(datasources[dsi].Tables[tbi].Columns, col)

		}
	*/

	//return datasource, nil

}
