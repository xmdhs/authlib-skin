package model

type API[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type User struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,sha256"`
	Name     string `validate:"required,min=3,max=16"`
}
