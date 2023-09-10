package yggdrasil

type Pass struct {
	// 目前只能是 email
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Authenticate struct {
	Agent struct {
		Name    string `json:"name" validate:"required,eq=Minecraft"`
		Version int    `json:"version" validate:"required,eq=1"`
	} `json:"agent"`
	ClientToken string `json:"clientToken"`
	RequestUser bool   `json:"requestUser"`
	Pass
}

type Error struct {
	Cause        string `json:"cause,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type TokenUserID struct {
	ID string `json:"id"`
}

type Token struct {
	AccessToken       string      `json:"accessToken"`
	AvailableProfiles []UserInfo  `json:"availableProfiles,omitempty"`
	ClientToken       string      `json:"clientToken"`
	SelectedProfile   UserInfo    `json:"selectedProfile"`
	User              TokenUserID `json:"user,omitempty"`
}

type ValidateToken struct {
	// jwt
	AccessToken string `json:"accessToken" validate:"required,jwt"`
	ClientToken string `json:"clientToken"`
}

type RefreshToken struct {
	ValidateToken
	RequestUser     bool     `json:"requestUser"`
	SelectedProfile UserInfo `json:"selectedProfile"`
}

type UserInfo struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Properties []UserProperties `json:"properties"`
}

type UserProperties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}
