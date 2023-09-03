package yggdrasil

type Authenticate struct {
	Agent struct {
		Name    string `json:"name" validate:"required,eq=Minecraft"`
		Version int    `json:"version" validate:"required,eq=1"`
	} `json:"agent"`
	ClientToken string `json:"clientToken"`
	Password    string `json:"password" validate:"required"`
	RequestUser bool   `json:"requestUser"`
	Username    string `json:"username" validate:"required"`
}

type Error struct {
	Cause        string `json:"cause,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type Token struct {
	AccessToken       string         `json:"accessToken"`
	AvailableProfiles []TokenProfile `json:"availableProfiles"`
	ClientToken       string         `json:"clientToken"`
	SelectedProfile   TokenProfile   `json:"selectedProfile"`
	User              TokenUser      `json:"user,omitempty"`
}

type TokenProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TokenUser struct {
	ID         string `json:"id"`
	Properties []any  `json:"properties"`
}
