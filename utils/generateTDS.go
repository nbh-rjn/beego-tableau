package utils

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func GenerateTDSFile(filenameTDS string, datasources []models.DatasourceStruct) error {
	// create xml for file content
	tdsBody, err := generateTDSBody("test", true, "win", "18.1", "https://10ax.online.tableau.com", "http://www.tableausoftware.com/xml/user", datasources[0])
	if err != nil {
		return err
	}

	// create tds file
	file, err := os.Create("storage/" + filenameTDS)
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
				Records: getMetadataRecords(datasource),
			},
		},
		Aliases: models.Aliases{
			Enabled: "yes",
		},
		Columns: getColumns(datasource),
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
			Name:  fmt.Sprintf("[%s]", table.TableName),
			Table: fmt.Sprintf("[%s].[%s]", datasource.Schema, table.TableName),
			Type:  "table",
		}
		relations = append(relations, relation)
	}
	return relations
}

func getMetadataRecords(datasource models.DatasourceStruct) []models.MetadataRecord {

	var metadatarecords []models.MetadataRecord

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			metadatarecord := models.MetadataRecord{
				Class:        "column",
				RemoteName:   fmt.Sprintf("%s.%s", table.TableName, column.ColumnName),   //
				RemoteType:   3,                                                          // idk what this value is yet so i will leave it hardcoded like the example for now
				LocalName:    fmt.Sprintf("[%s.%s]", table.TableName, column.ColumnName), //
				ParentName:   fmt.Sprintf("[%s]", table.TableName),
				LocalType:    standardiseDatatypes(column.ColumnType),
				ContainsNull: true,                                                     // assuming, since no nullflag in csv
				RemoteAlias:  fmt.Sprintf("%s.%s", table.TableName, column.ColumnName), //
				Ordinal:      6,                                                        // idk what this value is yet so i will leave it hardcoded like the example for now
			}
			metadatarecords = append(metadatarecords, metadatarecord)
		}
	}

	return metadatarecords

}

func getColumns(datasource models.DatasourceStruct) []models.Column {
	var columns []models.Column

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			col := models.Column{
				Caption:  fmt.Sprintf("%s.%s", table.TableName, column.ColumnName),   //
				Datatype: standardiseDatatypes(column.ColumnType),                    //"integer",
				Name:     fmt.Sprintf("[%s.%s]", table.TableName, column.ColumnName), //
				Role:     "dimension",
				Type:     "nominal",
			}

			columns = append(columns, col)
		}
	}
	return columns
}

func standardiseDatatypes(columnType string) string {
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
			return TableauType
		}
	}
	return "N/A"
}

func Contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func standardiseDatatypesOLD(old_data_type string) string {
	typeMaps := map[string]string{
		// Mapping old data types to standardized data types

		"char":                        "string",
		"varchar":                     "string",
		"text":                        "string",
		"character varying":           "string",
		"character":                   "string",
		"uuid":                        "string",
		"nvarchar2":                   "string",
		"nchar":                       "string",
		"string":                      "string",
		"national character varying":  "string",
		"national character":          "string",
		"character large object":      "string",
		"clob":                        "string",
		"long":                        "string",
		"long text":                   "string",
		"mediumtext":                  "string",
		"tinytext":                    "string",
		"long varchar":                "string",
		"longnvarchar":                "string",
		"uniqueidentifier":            "string",
		"nstring":                     "string",
		"nvarchar":                    "string",
		"nchar varying":               "string",
		"nclob":                       "string",
		"ntext":                       "string",
		"json":                        "string",
		"jsonb":                       "string",
		"xml":                         "string",
		"varchar(max)":                "string",
		"nvarchar(max)":               "string",
		"mpaa_rating":                 "string", //
		"decimal":                     "real",
		"numeric":                     "real",
		"float":                       "real",
		"double":                      "real",
		"real":                        "real",
		"money":                       "real",
		"smallmoney":                  "real",
		"tinyint":                     "integer",
		"smallint":                    "integer",
		"mediumint":                   "integer",
		"int":                         "integer",
		"integer":                     "integer",
		"bigint":                      "integer",
		"number":                      "integer",
		"smallserial":                 "integer",
		"serial":                      "integer",
		"bigserial":                   "integer",
		"int4":                        "integer",
		"int2":                        "integer",
		"date":                        "date",
		"year":                        "date",
		"year to month":               "date",
		"year to second":              "date",
		"month":                       "date",
		"day":                         "date",
		"yearmonth":                   "date",
		"year to fraction":            "date",
		"datetime":                    "datetime",
		"smalldatetime":               "datetime",
		"datetime2":                   "datetime",
		"datetimeoffset":              "datetime",
		"timestamp":                   "datetime",
		"timestamp without time zone": "datetime",
		"timestamp with time zone":    "datetime",
		"time":                        "datetime",
		"time without time zone":      "datetime",
		"time with time zone":         "datetime",
		"interval":                    "datetime",
		"hour":                        "datetime",
		"minute":                      "datetime",
		"second":                      "datetime",
		"fraction":                    "datetime",
		"timetz":                      "datetime",
		"boolean":                     "boolean",
		"bool":                        "boolean",
		"bit":                         "boolean",
	}

	// Lookup the old_data_type in the map
	new_data_type, ok := typeMaps[old_data_type]
	if !ok {
		// Handle case where old_data_type is not found (optional)
		return "TEST" // or some default value
	}

	return new_data_type
}
