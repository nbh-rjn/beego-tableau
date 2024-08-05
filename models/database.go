package models

import (
	"fmt"

	"github.com/beego/beego/orm"
)

type CredentialsTable struct {
	Id        int    `orm:"column(id);pk;auto"`
	PATName   string `orm:"column(pat_name);size(255)"`
	PATSecret string `orm:"column(pat_secret);size(255)"`
	SiteID    string `orm:"column(site_id);size(255)"`
}

type LabelsTable struct {
	Id            int    `orm:"column(id);pk;auto"`
	LabelName     string `orm:"column(label_name);size(255)"`
	LabelCategory string `orm:"column(label_category);size(255)"`
	SiteID        string `orm:"column(site_id);size(255)"`
}

type ProjectsTable struct {
	Id          int    `orm:"column(id);pk;auto"`
	ProjectName string `orm:"column(project_name);size(255)"`
	ProjectID   string `orm:"column(project_id);size(255)"`
	SiteID      string `orm:"column(site_id);size(255)"`
}

type DatasourcesTable struct {
	Id             int    `orm:"column(id);pk;auto"`
	DatasourceName string `orm:"column(datasource_name);size(255)"`
	DatasourceID   string `orm:"column(datasource_id);size(255)"`
	SiteID         string `orm:"column(site_id);size(255)"`
}

func init() {
	orm.RegisterModel(new(LabelsTable))
	orm.RegisterModel(new(CredentialsTable))
	orm.RegisterModel(new(ProjectsTable))
	orm.RegisterModel(new(DatasourcesTable))
}

func SaveCredentialsDB(patName string, patSecret string, siteID string) {
	o := orm.NewOrm()

	c := CredentialsTable{
		PATName:   patName,
		PATSecret: patSecret,
		SiteID:    siteID,
	}

	o.Insert(&c)

}

func SaveAttributesDB(param string, siteID string, attributes []map[string]interface{}) error {

	o := orm.NewOrm()

	for _, attribute := range attributes {

		switch param {
		case "datalabels":
			label := LabelsTable{
				LabelName:     string(attribute["name"].(string)),
				LabelCategory: string(attribute["category"].(string)),
				SiteID:        siteID,
			}
			_, err := o.Insert(&label)
			if err != nil {
				return err
			}

		case "datasources":
			datasource := DatasourcesTable{
				DatasourceName: string(attribute["name"].(string)),
				DatasourceID:   string(attribute["id"].(string)),
				SiteID:         siteID,
			}
			_, err := o.Insert(&datasource)
			if err != nil {
				return err
			}

		case "projects":
			project := ProjectsTable{
				ProjectName: string(attribute["name"].(string)),
				ProjectID:   string(attribute["id"].(string)),
				SiteID:      siteID,
			}
			_, err := o.Insert(&project)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid attribute type")
		}
	}
	return nil
}
