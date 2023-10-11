package model

import "github.com/golang-jwt/jwt/v5"

type API[T any] struct {
	Code APIStatus `json:"code"`
	Data T         `json:"data"`
	Msg  string    `json:"msg,omitempty"`
}

type List[T any] struct {
	Total int `json:"total"`
	List  []T `json:"list"`
}

type UserReg struct {
	Email        string `validate:"required,email"`
	Password     string `validate:"required,min=6,max=50"`
	Name         string `validate:"required,min=3,max=16"`
	CaptchaToken string
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

type UserInfo struct {
	UID     int    `json:"uid"`
	UUID    string `json:"uuid"`
	IsAdmin bool   `json:"is_admin"`
}

type ChangePasswd struct {
	Old string `json:"old"`
	New string `json:"new" validate:"required,min=6,max=50"`
}

type UserList struct {
	UserInfo
	Email     string `json:"email"`
	RegIp     string `json:"reg_ip"`
	Name      string `json:"name"`
	IsDisable bool   `json:"is_disable"`
}

type ChangeName struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
}

type Config struct {
	Captcha         Captcha `json:"captcha"`
	AllowChangeName bool
	ServerName      string `json:"serverName"`
}

type EditUser struct {
	Email       string `json:"email" validate:"omitempty,email"`
	Name        string `json:"name" validate:"omitempty,min=3,max=16"`
	Password    string `json:"password"`
	IsAdmin     *bool  `json:"is_admin"`
	IsDisable   *bool  `json:"is_disable"`
	DelTextures bool   `json:"del_textures"`
}

type Login struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password"`
	CaptchaToken string
}

type LoginRep struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
}
