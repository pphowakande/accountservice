package main

type Account struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}
