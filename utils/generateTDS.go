package utils

import (
	"beego-project/models"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

func GenerateTDSFile(filePath string, datasource models.DatasourceStruct) error {
	// create xml for file content
	tdsBody, err := generateTDSBody("test", true, "win", "18.1", "https://10ax.online.tableau.com", "http://www.tableausoftware.com/xml/user", datasource)
	if err != nil {
		return err
	}

	// create tds file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write to file
	if _, err = file.Write(tdsBody); err != nil {
		return err
	}

	return nil
}
func generateTDSBody(formattedName string, inline bool, sourcePlatform string, version string, xmlBase string, xmlnsUser string, datasource models.DatasourceStruct) ([]byte, error) {
	TableauMetadataRecords, err := getMetadataRecords(datasource)
	if err != nil {
		return nil, err
	}

	columns, err := getColumns(datasource)
	if err != nil {
		return nil, err
	}
	ds := models.DatasourceGeneration{
		FormattedName:  formattedName,
		Inline:         inline,
		SourcePlatform: sourcePlatform,
		Version:        version,
		XMLBase:        xmlBase,
		XMLNSUser:      xmlnsUser,
		RepositoryLocation: models.RepositoryLocation{
			ID:       datasource.Datasource,
			Path:     "/t/testsiteintern/datasources", // replace
			Revision: "1.0",
			Site:     "testsiteintern", // replace
		},
		Connection: models.Connection{
			Class:     "sqlserver",
			DBName:    datasource.Database,
			Server:    datasource.DBType,
			Relations: getRelations(datasource),
			MetadataRecords: models.MetadataRecords{
				Records: TableauMetadataRecords,
			},
		},
		Aliases: models.Aliases{
			Enabled: "yes",
		},
		Columns: columns,
	}

	xmlData, err := xml.MarshalIndent(ds, "", "\t")
	if err != nil {
		return nil, err
	}

	xmlHeader := []byte(xml.Header)
	xmlOutput := append(xmlHeader, xmlData...)
	return xmlOutput, nil
}

func getRelations(datasource models.DatasourceStruct) []models.Relation {
	var relations []models.Relation
	for _, table := range datasource.Tables {
		relation := models.Relation{
			Name:  table.TableName, //fmt.Sprintf("[%s]", table.TableName),
			Table: fmt.Sprintf("[%s].[%s]", datasource.Schema, table.TableName),
			Type:  "table",
		}
		relations = append(relations, relation)
	}
	return relations
}

func getMetadataRecords(datasource models.DatasourceStruct) ([]models.MetadataRecord, error) {

	var metadatarecords []models.MetadataRecord
	ordinalCount := 0

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			TableauLocalType, err := mapLocalType(column.ColumnType)
			if err != nil {
				return nil, err
			}
			TableauRemoteType, err := mapRemoteType(TableauLocalType)
			if err != nil {
				return nil, err
			}
			metadatarecord := models.MetadataRecord{
				Class:        "column",
				RemoteName:   fmt.Sprintf("%s.%s", table.TableName, column.ColumnName),   //
				RemoteType:   TableauRemoteType,                                          //3,                                                          // assuming
				LocalName:    fmt.Sprintf("[%s.%s]", table.TableName, column.ColumnName), //
				ParentName:   fmt.Sprintf("[%s]", table.TableName),
				LocalType:    TableauLocalType,
				ContainsNull: true,                                                     // assuming, since no nullflag in csv
				RemoteAlias:  fmt.Sprintf("%s.%s", table.TableName, column.ColumnName), //
				Ordinal:      ordinalCount,                                             // assuming
			}
			metadatarecords = append(metadatarecords, metadatarecord)
			ordinalCount = ordinalCount + 1
		}
	}

	return metadatarecords, nil

}

func getColumns(datasource models.DatasourceStruct) ([]models.Column, error) {
	var columns []models.Column

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			TableauDataType, err := mapLocalType(column.ColumnType)
			if err != nil {
				return nil, err
			}
			col := models.Column{
				Caption:  fmt.Sprintf("%s.%s", table.TableName, column.ColumnName),
				Datatype: TableauDataType,
				Name:     fmt.Sprintf("[%s.%s]", table.TableName, column.ColumnName),
				Role:     "dimension",
				Type:     "nominal",
			}

			columns = append(columns, col)
		}
	}
	return columns, nil
}

func mapLocalType(columnType string) (string, error) {
	typeMaps := map[string][]string{
		"string":   {"char", "varchar", "text", "character varying", "character", "uuid", "nvarchar2", "nchar", "string", "national character varying", "national character", "character large object", "clob", "long", "long text", "mediumtext", "tinytext", "long varchar", "longnvarchar", "uniqueidentifier", "nstring", "nvarchar", "longnvarchar", "nchar varying", "nclob", "ntext", "json", "jsonb", "xml", "varchar(max)", "nvarchar(max)", "mpaa_rating"},
		"real":     {"decimal", "numeric", "float", "double", "real", "money", "smallmoney"},
		"integer":  {"tinyint", "smallint", "mediumint", "int", "integer", "bigint", "number", "smallserial", "serial", "bigserial", "int4", "int2"},
		"date":     {"date", "year", "year to month", "year to second", "month", "day", "yearmonth", "year to fraction"},
		"datetime": {"datetime", "smalldatetime", "datetime2", "datetimeoffset", "timestamp", "timestamp without time zone", "timestamp with time zone", "time", "time without time zone", "time with time zone", "interval", "hour", "minute", "second", "fraction", "timetz"},
		"boolean":  {"boolean", "bool", "bit"},
	}
	for TableauType, datatypes := range typeMaps {
		if Contains(datatypes, strings.ToLower(columnType)) {
			return TableauType, nil
		}
	}
	return "", fmt.Errorf("could not convert to Tableau-compatible data type")
}

func Contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func mapRemoteType(columnType string) (int, error) {
	typeMaps := map[string]int{
		"string":   130, // DBTYPE_BSTR
		"real":     5,   // DBTYPE_R8
		"integer":  3,   // DBTYPE_I4
		"date":     133, // DBTYPE_DBDATE
		"datetime": 135, // DBTYPE_DBTimeStamp
		"boolean":  11,  // DBTYPE_BOOL
	}

	if oleDBType, exists := typeMaps[columnType]; exists {
		return oleDBType, nil
	}
	return 12, errors.New("unknown column type")
}
