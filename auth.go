package main

type User struct {
	Login string
}

func userFromToken(token string) *User {
	// FIXME: JWT, Oauth2 ...
	if token == "aditya123Pass" {
		return &User{"Aditya"}
	}
	return nil
}
