package utils

import (
	"beego-project/lib"
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

	var found bool = false
	var errorMsg error

	// for every table in the already existing data source
	for _, table := range dsTDS.Connection.Relations {

		// fetch its ID and those of its columns
		//tableID, columnIDs, err := lib.TableauGetAssetIDs(dsCSV.Database, strings.Trim(table.Name, "[]"))
		// / fmt.Println("dsCSV.Database, table.Name", dsCSV.Database, table.Name)
		tableID, columnIDs, err := lib.TableauGetAssetIDs(dsCSV.Database, table.Name)

		if err != nil {
			errorMsg = err
			return dsCSV, err
		}

		// initially assume it isnt in our csv
		found = false

		// check if that table has already been parsed from csv
		for _, t := range dsCSV.Tables {
			//if strings.Trim(table.Name, "[]") == t.TableName {

			// / fmt.Println("table.Name == t.TableName, ", table.Name, t.TableName)
			if table.Name == t.TableName {

				found = true

				// for every column of that table in already existing data source
				for _, col := range dsTDS.Connection.MetadataRecords.Records {

					//if strings.Trim(col.ParentName, "[]") == strings.Trim(table.Name, "[]") {

					// / fmt.Println("strings.Trim(col.ParentName, xx ) == table.Name", strings.Trim(col.ParentName, "[]"), table.Name)
					if strings.Trim(col.ParentName, "[]") == table.Name {

						// if the column is not already parsed from csv
						if !containsColumn(strings.Trim(col.RemoteName, "[]"), table.Name, dsCSV.Tables) {

							// / fmt.Println("columnIDs[col.RemoteName]", columnIDs[col.RemoteName])

							label, err := lib.TableauGetAssetLabel("column", columnIDs[col.RemoteName])
							if err != nil {
								label = ""
							}
							t.Columns = append(t.Columns, models.ColumnStruct{
								ColumnName:        strings.Split(col.RemoteName, ".")[1],
								ColumnType:        col.LocalType,
								ColumnDescription: "",
								DataElements:      label,
							})
						}

					}
				}

			}
		}
		if !found {
			// if not found just recreate the entire table as is

			var columns []models.ColumnStruct

			for _, col := range dsTDS.Connection.MetadataRecords.Records {
				//if strings.Trim(col.ParentName, "[]") == strings.Trim(table.Name, "[]") {

				fmt.Println("strings.Trim(col.ParentName,) == table.Name", strings.Trim(col.ParentName, "[]"), table.Name)
				if strings.Trim(col.ParentName, "[]") == table.Name {

					fmt.Println("!found ; columnIDs[col.RemoteName]", columnIDs[col.RemoteName])
					label, err := lib.TableauGetAssetLabel("column", columnIDs[col.RemoteName])
					if err != nil {
						label = ""

					}

					columns = append(columns, models.ColumnStruct{
						ColumnName:        strings.Split(col.RemoteName, ".")[1],
						ColumnType:        col.LocalType,
						ColumnDescription: "",
						DataElements:      label,
					})
				}
			}
			fmt.Println("tableID", tableID)
			label, err := lib.TableauGetAssetLabel("table", tableID)
			if err != nil {
				label = ""
			}
			dsCSV.Tables = append(dsCSV.Tables, models.TableStruct{
				Id:              "",
				TableName:       table.Name, //strings.Trim(table.Name, "[]"),
				TableType:       "",
				ContentProfiles: label,
				Columns:         columns,
			})
		}

	}

	return dsCSV, errorMsg
}

func containsColumn(columnname, tableName string, tables []models.TableStruct) bool {
	for _, table := range tables {
		if table.TableName == tableName {
			for _, column := range table.Columns {
				if column.ColumnName == columnname {
					return true
				}
			}
		}
	}
	return false
}
