package utils

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Datasource struct {
	XMLName            xml.Name           `xml:"datasource"`
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
	XMLName         xml.Name         `xml:"connection"`
	Class           string           `xml:"class,attr"`
	DBName          string           `xml:"dbname,attr"`
	Server          string           `xml:"server,attr"`
	Relations       []Relation       `xml:"relation"`
	MetadataRecords []MetadataRecord `xml:"metadata-records"`
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

func GenerateXML(formattedName string, inline bool, sourcePlatform string, version string, xmlBase string, xmlnsUser string, datasource DatasourceStruct) ([]byte, error) {
	ds := Datasource{
		FormattedName:  formattedName,
		Inline:         inline,
		SourcePlatform: sourcePlatform,
		Version:        version,
		XMLBase:        xmlBase,
		XMLNSUser:      xmlnsUser,
		RepositoryLocation: RepositoryLocation{
			//path and site hardcoded for now
			ID:       datasource.Datasource, //"newdatasource5",
			Path:     "/t/testsiteintern/datasources",
			Revision: "1.0",
			Site:     "testsiteintern",
		},
		Connection: Connection{
			Class:           "sqlserver",
			DBName:          datasource.Database, //"SampleDatabase2",
			Server:          datasource.DBType,   //"127.0.0.1",
			Relations:       getrelations(datasource),
			MetadataRecords: getmetadatarecords(datasource),
		},
		Aliases: Aliases{
			Enabled: "yes",
		},
		Columns: getcolumns(datasource),
	}

	xmlData, err := xml.MarshalIndent(ds, "", "\t")
	if err != nil {
		return nil, err
	}

	xmlHeader := []byte(xml.Header)
	xmlOutput := append(xmlHeader, xmlData...)
	return xmlOutput, nil
}

func Gen_xml(datasources []DatasourceStruct) {
	xmlData, err := GenerateXML("test", true, "win", "18.1", "https://10ax.online.tableau.com", "http://www.tableausoftware.com/xml/user", datasources[0])
	if err != nil {
		fmt.Println("Error generating XML:", err)
		return
	}

	file, err := os.Create("xml.tds")
	if err != nil {
		fmt.Println("error creating file . . . ")
		return

	}
	defer file.Close()

	_, err = file.Write(xmlData)
	if err != nil {
		fmt.Println("error writing data . . .")
		return
	}
	fmt.Println(string(xmlData))
}

func getrelations(datasource DatasourceStruct) []Relation {
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

func getmetadatarecords(datasource DatasourceStruct) []MetadataRecord {

	var metadatarecords []MetadataRecord

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			metadatarecord := MetadataRecord{
				Class:        "column",
				RemoteName:   column.ColumnName,
				RemoteType:   3, // idk what this value is yet so i will leave it hardcoded like the example for now
				LocalName:    fmt.Sprintf("[%s]", column.ColumnName),
				ParentName:   fmt.Sprintf("[%s]", table.TableName),
				LocalType:    column.ColumnType,
				ContainsNull: true, // assuming, since no nullflag in csv
				RemoteAlias:  column.ColumnName,
				Ordinal:      6, // idk what this value is yet so i will leave it hardcoded like the example for now
			}
			metadatarecords = append(metadatarecords, metadatarecord)
		}
	}

	return metadatarecords

}

func getcolumns(datasource DatasourceStruct) []Column {
	var columns []Column

	for _, table := range datasource.Tables {
		for _, column := range table.Columns {
			col := Column{
				Caption:  column.ColumnName,                      //"OrderID",
				Datatype: column.ColumnType,                      //"integer",
				Name:     fmt.Sprintf("[%s]", column.ColumnName), //"[OrderID]",
				Role:     "dimension",
				Type:     "nominal",
			}

			columns = append(columns, col)
		}
	}
	return columns
}
