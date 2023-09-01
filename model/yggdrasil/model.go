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
