package lib

import "net/http"

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
