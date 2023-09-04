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

type Token struct {
	AccessToken       string         `json:"accessToken"`
	AvailableProfiles []TokenProfile `json:"availableProfiles,omitempty"`
	ClientToken       string         `json:"clientToken"`
	SelectedProfile   TokenProfile   `json:"selectedProfile"`
	User              TokenUser      `json:"user,omitempty"`
}

type TokenProfile struct {
	// 就是 uuid
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TokenUser struct {
	ID         string `json:"id"`
	Properties []any  `json:"properties"`
}

type ValidateToken struct {
	AccessToken string `json:"accessToken" validate:"required,jwt"`
	ClientToken string `json:"clientToken"`
}

type RefreshToken struct {
	ValidateToken
	RequestUser     bool         `json:"requestUser"`
	SelectedProfile TokenProfile `json:"selectedProfile"`
}
