package utils

import (
	"beego-project/models"
	"fmt"
	"io"
	"sort"
)

func Tableau_sync(records []FCRecords, siteID string) {

	hierarchy_records := organize_records(records)

	// LETS GOOOOO HOGAYA
	for _, datasource := range hierarchy_records {
		fmt.Printf("%s \n", datasource.Datasource)
		for _, table := range datasource.Tables {
			fmt.Printf("\t %s \n", table.TableName)
			for _, col := range table.Columns {
				fmt.Printf("\t \t %s \n", col.ColumnName)
			}
		}
	}

	Gen_xml(hierarchy_records)

}

func datasourceExists(siteID string, datasource string) bool {

	response, err := Tableau_get_data_sources(models.Get_token(), siteID)
	if err != nil {
		fmt.Println("error communicating with tableau API")
		return false
	}

	// read body of response
	bodyread, _ := io.ReadAll(response.Body)

	// utility function to extract relevant info
	datasources, _ := Extract_data_sources_xml(string(bodyread))

	sort.Strings(datasources)
	i := sort.SearchStrings(datasources, datasource)
	found := i < len(datasources) && datasources[i] == datasource

	return found
}
