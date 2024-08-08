package lib

import (
	"beego-project/models"
	"fmt"
)

func TableauCreateCategory(category string) error {

	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/labelCategories"

	// request body
	payload := fmt.Sprintf(
		`<tsRequest>
		<labelCategory name="%s"
			description="%s" />
	</tsRequest>`, category, category)

	response, err := TableauRequest(url, payload, "POST", "xml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return err
}

func TableauLabelAsset(label string, category string, assetType string, assetID string) error {

	if label == "" || category == "" {
		return nil
	}

	// in case category doesnt exist
	TableauCreateCategory(category)

	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
		   <labelValue name="%s"
		     category="%s"
		     description="Created via API" />
		</tsRequest>`, label, category)

	url := models.TableauURL() + "sites/" + models.GetSiteID() + "/labelValues"

	// new put request
	response, err := TableauRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	response.Body.Close()

	// XML payload for applying label value
	payload = fmt.Sprintf(`
		<tsRequest>
		  <contentList>
		    <content contentType="%s" id="%s" />
		  </contentList>
		  <label
		      value="%s"/>
		</tsRequest>`, assetType, assetID, label)

	url = models.TableauURL() + "sites/" + models.GetSiteID() + "/labels"

	// PUT request
	response, err = TableauRequest(url, payload, "PUT", "xml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
