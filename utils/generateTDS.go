package utils

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"os"
)

func GenerateTDSFile(filenameTDS string, datasources []models.DatasourceStruct) error {
	// create xml for file content
	tdsBody, err := generateTDSBody("test", true, "win", "18.1", "https://10ax.online.tableau.com", "http://www.tableausoftware.com/xml/user", datasources[0])
	if err != nil {
		return err
	}

	// create tds file
	file, err := os.Create(filenameTDS)
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

func standardiseDatatypes(old_data_type string) string {
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
		"mpaa_rating":                 "string",
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
