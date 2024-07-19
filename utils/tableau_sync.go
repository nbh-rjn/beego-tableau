package utils

func TableauSyncRecords(datasourceRecords []DatasourceStruct, siteID string) error {
	for _, datasourceRecord := range datasourceRecords {

		// generate tds file for each datasource and publish it
		err := GenerateTDSFile(datasourceRecords)
		if err != nil {
			return err
		}

		err = PublishDatasource(siteID, datasourceRecord.Datasource)
		if err != nil {
			return err
		}

	}
	return nil
}

/*
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
*/
