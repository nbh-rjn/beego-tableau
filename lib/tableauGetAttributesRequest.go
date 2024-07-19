package lib

import (
	"beego-project/models"
	"fmt"
	"net/http"
)

func TableauGetAttribute(param string, site_id string) (*http.Response, error) {
	var attribute string

	switch param {
	case "datalabels":
		attribute = "/labelValues"
	case "datasources":
		attribute = "/datasources"
	case "projects":
		attribute = "/projects"
	default:
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
