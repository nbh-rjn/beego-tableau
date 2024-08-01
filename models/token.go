package models

var credential_token = ""
var site_id = ""

// save credentials token for the session
func SaveCredentials(t string, s string) {
	credential_token = t
	site_id = s
}

// so it can be accessed when servicing other endpoints
func Get_token() string {
	return credential_token
}

func Get_siteID() string {
	return site_id
}

func TableauURL() string {
	return "https://10ax.online.tableau.com/api/3.21/"
}
