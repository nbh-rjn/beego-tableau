package utils

import (
	"beego-project/lib"

	"github.com/pkg/errors"
)

func TableauSyncRecords(filenameCSV string, siteID string) error {

	// tds filename
	filenameTDS := "sync.tds"

	// parse CSV to slice of structs
	datasourceRecords := ParseCSV(filenameCSV)
	if datasourceRecords == nil {
		return errors.New("Could not parse raw CSV file")
	}

	// one datasource at a time
	// can handle more than one datasource per CSV file
	for _, datasourceRecord := range datasourceRecords {

		// generate tds file for each datasource struct
		if err := GenerateTDSFile(filenameTDS, datasourceRecords); err != nil {
			return err
		}

		// publish it
		if err := lib.PublishDatasource(filenameTDS, siteID, datasourceRecord.Datasource); err != nil {
			return err
		}

	}
	return nil
}
