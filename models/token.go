package models

var credential_token = ""

// save credentials token for the session
func SaveToken(t string) {
	credential_token = t
}

// so it can be accessed when servicing other endpoints
func Get_token() string {
	return credential_token
}
