package models

import "encoding/xml"

// auth req to beego
type CredentialStruct struct {
	PersonalAccessTokenName   string `json:"personalAccessTokenName"`
	PersonalAccessTokenSecret string `json:"personalAccessTokenSecret"`
	ContentUrl                string `json:"contentUrl"`
}

// auth resp from tableau
type AuthResponse struct {
	TSResponse
	Credentials Credentials `xml:"credentials"`
}

type Credentials struct {
	Token                     string `xml:"token,attr"`
	EstimatedTimeToExpiration string `xml:"estimatedTimeToExpiration,attr"`
	Site                      Site   `xml:"site"`
	User                      User   `xml:"user"`
}

type Site struct {
	ID         string `xml:"id,attr"`
	ContentURL string `xml:"contentUrl,attr"`
}

type User struct {
	ID string `xml:"id,attr"`
}

// sync
type AttributeMap struct {
	DataElements   string `json:"data_elements"`
	ContentProfile string `json:"content_profile"`
}

type InstanceMap map[string]string

type SyncRequest struct {
	Filename        string       `json:"filename"`
	SiteID          string       `json:"siteID"`
	CreateNewAssets bool         `json:"create_new_assets"`
	EntityType      string       `json:"entity_type"`
	AttributeMap    AttributeMap `json:"attribute_map"`
	InstanceMap     InstanceMap  `json:"instance_map"`
}

// get attributes body
type SiteRequest struct {
	SiteID string `json:"siteID"`
}

// download datasource body
type DownloadRequest struct {
	SiteRequest
	DatasourceID string `json:"datasourceID"`
}

// tableau response body
type TSResponse struct {
	XMLName xml.Name `xml:"tsResponse"`
}

// label attributes
type SiteL struct {
	ID string `xml:"id,attr"`
}

type LabelValue struct {
	XMLName     xml.Name `xml:"labelValue"`
	Name        string   `xml:"name,attr"`
	Category    string   `xml:"category,attr"`
	Description string   `xml:"description,attr"`
	Internal    bool     `xml:"internal,attr"`
	Elevated    bool     `xml:"elevatedDefault,attr"`
	BuiltIn     bool     `xml:"builtIn,attr"`
	Site        SiteL    `xml:"site"`
}

type LabelValueList struct {
	XMLName     xml.Name     `xml:"labelValueList"`
	LabelValues []LabelValue `xml:"labelValue"`
}

type LabelValueResponse struct {
	TSResponse
	LabelValueList LabelValueList `xml:"labelValueList"`
}

// datasource attributes
type Datasource struct {
	XMLName xml.Name `xml:"datasource"`
}

type DatasourceElement struct {
	Datasource
	ContentUrl string `xml:"contentUrl,attr"`
	Name       string `xml:"name,attr"`
	Id         string `xml:"id,attr"`
}

type Datasources struct {
	XMLName    xml.Name            `xml:"datasources"`
	Datasource []DatasourceElement `xml:"datasource"`
}

type DatasourceResponse struct {
	TSResponse
	Datasources Datasources `xml:"datasources"`
}

// project attributes
type Owner struct {
	ID string `xml:"id,attr"`
}
type Project struct {
	ID                 string `xml:"id,attr"`
	Name               string `xml:"name,attr"`
	Description        string `xml:"description,attr"`
	CreatedAt          string `xml:"createdAt,attr"`
	UpdatedAt          string `xml:"updatedAt,attr"`
	ContentPermissions string `xml:"contentPermissions,attr"`
	Owner              Owner  `xml:"owner"`
}
type Pagination struct {
	PageNumber     int `xml:"pageNumber,attr"`
	PageSize       int `xml:"pageSize,attr"`
	TotalAvailable int `xml:"totalAvailable,attr"`
}

type ProjectResponse struct {
	TSResponse
	Pagination Pagination `xml:"pagination"`
	Projects   []Project  `xml:"projects>project"`
}

// generate TDS

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

// parse CSV

type CSVRecords struct {
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
