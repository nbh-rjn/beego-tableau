package utils

import (
	"encoding/xml"
	"fmt"
	"os"
)

type DatasourceGeneration struct {
	Datasource
	FormattedName      string             `xml:"formatted-name,attr"`
	Inline             bool               `xml:"inline,attr"`
	SourcePlatform     string             `xml:"source-platform,attr"`
	Version            string             `xml:"version,attr"`
	XMLBase            string             `xml:"xml:base,attr"`
	XMLNSUser          string             `xml:"xmlns:user,attr"`
	RepositoryLocation RepositoryLocation `xml:"repository-location"`
	Connection         Connection         `xml:"connection"`
	Aliases            Aliases            `xml:"aliases"`
	Columns            []Column           `xml:"column"`
}

type RepositoryLocation struct {
	XMLName  xml.Name `xml:"repository-location"`
	ID       string   `xml:"id,attr"`
	Path     string   `xml:"path,attr"`
	Revision string   `xml:"revision,attr"`
	Site     string   `xml:"site,attr"`
}

type Connection struct {
	XMLName         xml.Name        `xml:"connection"`
	Class           string          `xml:"class,attr"`
	DBName          string          `xml:"dbname,attr"`
	Server          string          `xml:"server,attr"`
	Relations       []Relation      `xml:"relation"`
	MetadataRecords MetadataRecords `xml:"metadata-records"`
}

type MetadataRecords struct {
	Records []MetadataRecord `xml:"metadata-record"`
}

type Relation struct {
	XMLName xml.Name `xml:"relation"`
	Name    string   `xml:"name,attr"`
	Table   string   `xml:"table,attr"`
	Type    string   `xml:"type,attr"`
}

type MetadataRecord struct {
	XMLName      xml.Name `xml:"metadata-record"`
	Class        string   `xml:"class,attr"`
	RemoteName   string   `xml:"remote-name"`
	RemoteType   int      `xml:"remote-type"`
	LocalName    string   `xml:"local-name"`
	ParentName   string   `xml:"parent-name"`
	LocalType    string   `xml:"local-type"`
	ContainsNull bool     `xml:"contains-null"`
	RemoteAlias  string   `xml:"remote-alias"`
	Ordinal      int      `xml:"ordinal"`
}

type Aliases struct {
	Enabled string `xml:"enabled,attr"`
}

type Column struct {
	XMLName  xml.Name `xml:"column"`
	Caption  string   `xml:"caption,attr"`
	Datatype string   `xml:"datatype,attr"`
	Name     string   `xml:"name,attr"`
	Role     string   `xml:"role,attr"`
	Type     string   `xml:"type,attr"`
}

func generateTDS(formattedName string, inline bool, sourcePlatform string, version string, xmlBase string, xmlnsUser string, datasource DatasourceStruct) ([]byte, error) {
	ds := DatasourceGeneration{
		FormattedName:  formattedName,
		Inline:         inline,
		SourcePlatform: sourcePlatform,
		Version:        version,
		XMLBase:        xmlBase,
		XMLNSUser:      xmlnsUser,
		RepositoryLocation: RepositoryLocation{
			ID:       datasource.Datasource,
			Path:     "/t/testsiteintern/datasources", // replace
			Revision: "1.0",
			Site:     "testsiteintern", // replace
		},
		Connection: Connection{
			Class:     "sqlserver",
			DBName:    datasource.Database,
			Server:    datasource.DBType,
			Relations: getRelations(datasource),
			MetadataRecords: MetadataRecords{
				Records: getMetadataRecords(datasource),
			},
		},
		Aliases: Aliases{
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

func GenerateTDSFile(datasources []DatasourceStruct) error {
	// create xml for file content
	tdsData, err := generateTDS("test", true, "win", "18.1", "https://10ax.online.tableau.com", "http://www.tableausoftware.com/xml/user", datasources[0])
	if err != nil {
		return err
	}

	// create tds file
	file, err := os.Create("sync.tds")
	if err != nil {
		return err
	}
	defer file.Close()

	// write to file
	_, err = file.Write(tdsData)
	if err != nil {
		return err
	}

	return nil
}

func getRelations(datasource DatasourceStruct) []Relation {
	var relations []Relation
	for _, table := range datasource.Tables {
		relation := Relation{
			Name:  fmt.Sprintf("[%s]", table.TableName),
			Table: fmt.Sprintf("[%s].[%s]", datasource.Schema, table.TableName),
			Type:  "table",
		}
		relations = append(relations, relation)
	}
	return relations
}

func getMetadataRecords(datasource DatasourceStruct) []MetadataRecord {

	var metadatarecords []MetadataRecord

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			metadatarecord := MetadataRecord{
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

func getColumns(datasource DatasourceStruct) []Column {
	var columns []Column

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			col := Column{
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
