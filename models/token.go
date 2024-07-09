package models

var ct = ""

func Set_token(t string) {
	ct = t
}

func Get_token() string {
	return ct
}
