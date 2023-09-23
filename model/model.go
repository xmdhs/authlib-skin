package model

import "github.com/golang-jwt/jwt/v5"

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

type TokenClaims struct {
	// token id 验证 token 是否过期
	Tid string `json:"tid"`
	// ClientToken Yggdrasil 协议中使用
	CID string `json:"cid"`
	// 用户 id
	UID int `json:"uid"`
	jwt.RegisteredClaims
}

type Captcha struct {
	Type    string `json:"type"`
	SiteKey string `json:"siteKey"`
}
