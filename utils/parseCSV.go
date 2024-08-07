package utils

import (
	"beego-project/models"
	"encoding/csv"
	"fmt"
	"os"
)

func ParseCSV(filename string) (models.DatasourceStruct, error) {

	var datasource models.DatasourceStruct

	// read file
	file, err := os.Open(filename)
	if err != nil {
		return datasource, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return datasource, err
	}

	// return if empty
	if len(records) < 2 {
		return datasource, fmt.Errorf("no records to parse in CSV")
	}
	/*
		[0]  Id,
		[1]  Datasource
		[2]  Host
		[3]  Port
		[4]  DatabaseType
		[5]  DBUsername
		[6]  Database
		[7]  Schema
		[8]  Table
		[9]  TableType
		[10] ContentProfiles
		[11] Column
		[12] ColumnType
		[13] ColumnDescription
		[14] DataElements
	*/
	// data source info
	datasource.Datasource = records[1][1]
	datasource.Host = records[1][2]
	datasource.Port = records[1][3]
	datasource.DBType = records[1][4]
	datasource.DBUsername = records[1][5]
	datasource.Database = records[1][6]
	datasource.Schema = records[1][7]

	// organize cols and tables within struct

	for i, record := range records {
		if i == 0 {
			continue
		}

		if record[8] != records[i-1][8] {
			datasource.Tables = append(datasource.Tables, models.TableStruct{
				Id:              record[0],
				TableName:       record[8],
				TableType:       record[9],
				ContentProfiles: record[10],
			})
		}
		currentTableIdx := len(datasource.Tables) - 1

		datasource.Tables[currentTableIdx].Columns = append(datasource.Tables[currentTableIdx].Columns, models.ColumnStruct{
			ColumnName:        record[11],
			ColumnType:        record[12],
			ColumnDescription: record[13],
			DataElements:      record[14],
		})

	}

	return datasource, nil
}
