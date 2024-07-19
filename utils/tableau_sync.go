package utils

import (
	"beego-project/lib"
	"fmt"

	"github.com/pkg/errors"
)

func TableauSyncRecords(filenameCSV string, siteID string) error {

	// parse CSV to slice of structs
	datasourceRecords := ParseCSV(filenameCSV)
	if datasourceRecords == nil {
		return errors.New("Could not parse raw CSV file")
	}

	// one datasource at a time
	// can handle more than one datasource per CSV file
	for _, datasourceRecord := range datasourceRecords {

		// tds filename
		filepathTDS := fmt.Sprintf("%s-%s.tds", datasourceRecord.Datasource, siteID)

		// generate tds file for each datasource struct
		if err := GenerateTDSFile(filepathTDS, datasourceRecords); err != nil {
			return err
		}

		// publish it
		if err := lib.PublishDatasource(filepathTDS, siteID, datasourceRecord.Datasource); err != nil {
			return err
		}

	}
	return nil
}
