package lib

import (
	"beego-project/models"
	"encoding/xml"
	"fmt"
	"io"
)

// returns first label on an asset
func TableauGetAssetLabel(assetType string, assetID string) (string, error) {
	// XML payload for creating label value
	payload := fmt.Sprintf(`
		<tsRequest>
			<contentList>
   				<content contentType="%s"
      			id="%s" />
 			</contentList>
		</tsRequest>`, assetType, assetID)

	// need v3.21
	url := models.TableauURL() + "sites/" + models.Get_siteID() + "/labels"

	// new put request
	response, err := MakeRequest(url, payload, "POST", "xml")
	if err != nil {
		return "", err
	}

	xmlData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseBody models.LabelResponse
	if err := xml.Unmarshal([]byte(xmlData), &responseBody); err != nil {
		return "", err
	}

	response.Body.Close()

	return responseBody.LabelList[0].Value, nil
}
