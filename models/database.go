package models

import (
	"github.com/beego/beego/orm"
)

type CredentialsTable struct {
	Id        int    `orm:"column(id);pk;auto"`
	PATName   string `orm:"column(pat_name);size(255)"`
	PATSecret string `orm:"column(pat_secret);size(255)"`
	SiteID    string `orm:"column(site_id);size(255)"`
}

type LabelsTable struct {
	Id        int    `orm:"column(id);pk;auto"`
	LabelName string `orm:"column(label_name);size(255)"`
	SiteID    string `orm:"column(site_id);size(255)"`
}

type ProjectsTable struct {
	Id          int    `orm:"column(id);pk;auto"`
	ProjectName string `orm:"column(project_name);size(255)"`
	SiteID      string `orm:"column(site_id);size(255)"`
}

type DatasourcesTable struct {
	Id             int    `orm:"column(id);pk;auto"`
	DatasourceName string `orm:"column(datasource_name);size(255)"`
	SiteID         string `orm:"column(site_id);size(255)"`
}

func init() {
	orm.RegisterModel(new(LabelsTable))
	orm.RegisterModel(new(CredentialsTable))
	orm.RegisterModel(new(ProjectsTable))
	orm.RegisterModel(new(DatasourcesTable))
}
