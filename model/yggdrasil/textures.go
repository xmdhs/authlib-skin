package yggdrasil

import (
	"encoding/base64"
	"encoding/json"
)

type UserTextures struct {
	// UUID
	ProfileID   string              `json:"profileId"`
	ProfileName string              `json:"profileName"`
	Textures    map[string]Textures `json:"textures"`
	// 时间戳 毫秒
	Timestamp string `json:"timestamp"`
}

type Textures struct {
	Url      string            `json:"url"`
	Metadata map[string]string `json:"metadata"`
}

func (u UserTextures) Base64() string {
	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
