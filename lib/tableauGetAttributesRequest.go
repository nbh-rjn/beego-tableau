package lib

import (
	"beego-project/models"
	"fmt"
	"net/http"
)

func TableauGetAttribute(param string, site_id string) (*http.Response, error) {
	attributeMap := map[string]string{
		"datalabels":  "/labelValues",
		"datasources": "/datasources",
		"projects":    "/projects",
	}

	attribute, found := attributeMap[param]
	if !found {
		return nil, fmt.Errorf("invalid attribute")
	}

	url := models.TableauURL() + "sites/" + site_id + attribute

	// make new get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", models.Get_token())

	// send using client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
