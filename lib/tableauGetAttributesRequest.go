package lib

import "net/http"

func TableauGetDataLabelValues(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + site_id + "/labelValues"

	// make new get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", token)

	// send using client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TableauGetDataSources(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.4/sites/" + site_id + "/datasources"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Tableau-Auth", token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func TableauGetProjects(token string, site_id string) (*http.Response, error) {
	url := "https://10ax.online.tableau.com/api/3.20/sites/" + site_id + "/projects"

	// make new get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tableau-Auth", token)

	// send using client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
