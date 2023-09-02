package model

type API[T any] struct {
	Code APIStatus `json:"code"`
	Data T         `json:"data"`
	Msg  string    `json:"msg"`
}

type User struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=50"`
	Name     string `validate:"required,min=3,max=16"`
}
