package models

var credential_token = ""

func Set_token(t string) {
	credential_token = t
}

func Get_token() string {
	return credential_token
}
