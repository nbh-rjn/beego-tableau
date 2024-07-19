package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

type FCRecords struct {
	Id                string `csv:"Id"`
	Datasource        string `csv:"Datasource"`
	Host              string `csv:"Host"`
	Port              string `csv:"Port"`
	DatabaseType      string `csv:"DatabaseType"`
	DBUsername        string `csv:"DBUsername"`
	Database          string `csv:"Database"`
	Schema            string `csv:"Schema"`
	Table             string `csv:"Table"`
	TableType         string `csv:"TableType"`
	ContentProfiles   string `csv:"ContentProfiles"`
	Column            string `csv:"Column"`
	ColumnType        string `csv:"ColumnType"`
	ColumnDescription string `csv:"ColumnDescription"`
	DataElements      string `csv:"DataElements"`
}
type ColumnStruct struct {
	ColumnName        string
	ColumnType        string
	ColumnDescription string
	DataElements      string
}

type TableStruct struct {
	Id              string
	TableName       string
	TableType       string
	ContentProfiles string
	Columns         []ColumnStruct
}

type DatasourceStruct struct {
	Datasource string
	Host       string
	Port       string
	Database   string
	Schema     string
	DBUsername string
	DBType     string
	Tables     []TableStruct
}

func ParseCSV(filename string) []DatasourceStruct {

	// open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	// read file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// create empty array
	var fcrecords []FCRecords

	for i, record := range records {
		if i == 0 {
			continue
		}
		// store record values in elements of struct
		fcrecord := FCRecords{
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

	// return parsed array
	return organizeRecords(fcrecords)
}

// takes fcrecords which is just a single line of csv file
// returns hierarchically arranged slice of structs
// each struct representing one datasource
func organizeRecords(records []FCRecords) []DatasourceStruct {

	var datasources []DatasourceStruct
	ds_idx, tb_idx := "", ""
	dsi, tbi := -1, -1

	for _, record := range records {

		if ds_idx != record.Datasource {
			ds := DatasourceStruct{
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

		} else {
			if tb_idx != record.Table {
				tb := TableStruct{
					Id:              record.Id,
					TableName:       record.Table,
					TableType:       record.TableType,
					ContentProfiles: record.ContentProfiles,
				}

				datasources[dsi].Tables = append(datasources[dsi].Tables, tb)
				tb_idx = record.Table
				tbi = tbi + 1
			}
			col := ColumnStruct{
				ColumnName:        record.Column,
				ColumnType:        record.ColumnType,
				ColumnDescription: record.ColumnDescription,
				DataElements:      record.DataElements,
			}
			datasources[dsi].Tables[tbi].Columns = append(datasources[dsi].Tables[tbi].Columns, col)

		}
	}

	return datasources

}
