package utils

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ExtractAttributes(response *http.Response, attributeType string) ([]map[string]interface{}, error) {

	var attributes []map[string]interface{}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	xmlData := string(responseBody)

	switch attributeType {

	case "datalabels":
		var tsResponse models.LabelValueResponse

		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, labelValue := range tsResponse.LabelValueList.LabelValues {
			attributes = append(attributes, map[string]interface{}{
				"name": labelValue.Name, "category": labelValue.Category,
			})
		}

	case "datasources":

		var tsResponse models.DatasourceResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, datasource := range tsResponse.Datasources.Datasource {
			attributes = append(attributes, map[string]interface{}{
				"name": datasource.Name, "id": datasource.Id,
			})
		}

	case "projects":

		var tsResponse models.ProjectResponse
		if err := xml.Unmarshal([]byte(xmlData), &tsResponse); err != nil {
			return nil, err
		}

		for _, project := range tsResponse.Projects {
			attributes = append(attributes, map[string]interface{}{
				"name": project.Name, "id": project.ID,
			})
		}

	default:
		return nil, fmt.Errorf("invalid attribute type")
	}

	return attributes, nil

}

func ExtractAssets(filePath string, dsCSV models.DatasourceStruct) (models.DatasourceStruct, error) {

	// read existing assets from datasource
	file, err := os.Open(filePath)
	if err != nil {
		return dsCSV, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
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

	// Populate the inner maps with column information
	for _, columnTDS := range dsTDS.Connection.MetadataRecords.Records {
		columnParent := strings.Trim(columnTDS.ParentName, "[]")

		if innerMap, ok := outerMap[columnParent]; ok {
			innerMap[columnTDS.RemoteName] = columnTDS.LocalType
		}
	}

	for _, table := range dsCSV.Tables {

		if columnsMap, exists := outerMap[table.TableName]; exists {

			for _, column := range table.Columns {

				delete(columnsMap, column.ColumnName)
			}

			for colName, dataType := range columnsMap {
				table.Columns = append(table.Columns, models.ColumnStruct{
					ColumnName:        colName,
					ColumnType:        dataType,
					ColumnDescription: "",
					DataElements:      "",
				})
			}

			delete(outerMap, table.TableName)
		}
	}

	// Append remaining tables from the map to the slice
	for tableName, columnsMap := range outerMap {
		var columns []models.ColumnStruct
		for columnName, dataType := range columnsMap {
			columns = append(columns, models.ColumnStruct{
				ColumnName:        columnName,
				ColumnType:        dataType,
				ColumnDescription: "",
				DataElements:      "",
			})
		}
		dsCSV.Tables = append(dsCSV.Tables, models.TableStruct{
			Id:              "",
			TableName:       tableName,
			TableType:       "",
			ContentProfiles: "",
			Columns:         columns,
		})
	}

	return dsCSV, nil

}
