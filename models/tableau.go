package models

var credentialsToken = ""
var siteID = ""

// save credentials token for the session
func SaveCredentials(t string, s string) {
	credentialsToken = t
	siteID = s
}

// so it can be accessed when servicing other endpoints
func GetToken() string {
	return credentialsToken
}

func GetSiteID() string {
	return siteID
}

func TableauURL() string {
	return "https://10ax.online.tableau.com/api/3.21/"
}
