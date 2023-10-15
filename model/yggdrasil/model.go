package yggdrasil

type Pass struct {
	// 目前只能是 email
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Authenticate struct {
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
	ID         string `json:"id"`
	Properties []any  `json:"properties,omitempty"`
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
	Properties []UserProperties `json:"properties,omitempty"`
}

type UserProperties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature,omitempty"`
}

type Session struct {
	AccessToken     string `json:"accessToken" validate:"required,jwt"`
	SelectedProfile string `json:"selectedProfile" validate:"required"`
	ServerID        string `json:"serverId"`
}

type Yggdrasil struct {
	Meta               YggdrasilMeta `json:"meta"`
	SignaturePublickey string        `json:"signaturePublickey"`
	SkinDomains        []string      `json:"skinDomains"`
}

type YggdrasilMeta struct {
	ImplementationName    string             `json:"implementationName"`
	ImplementationVersion string             `json:"implementationVersion"`
	Links                 YggdrasilMetaLinks `json:"links"`
	ServerName            string             `json:"serverName"`
	EnableProfileKey      bool               `json:"feature.enable_profile_key"`
}

type YggdrasilMetaLinks struct {
	Homepage string `json:"homepage"`
	Register string `json:"register"`
}

type Certificates struct {
	ExpiresAt            string              `json:"expiresAt"`
	KeyPair              CertificatesKeyPair `json:"keyPair"`
	PublicKeySignature   string              `json:"publicKeySignature"`
	PublicKeySignatureV2 string              `json:"publicKeySignatureV2"`
	RefreshedAfter       string              `json:"refreshedAfter"`
}

type CertificatesKeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type PublicKeys struct {
	PlayerCertificateKeys []PublicKeyList `json:"playerCertificateKeys"`
	ProfilePropertyKeys   []PublicKeyList `json:"profilePropertyKeys"`
}

type PublicKeyList struct {
	PublicKey string `json:"publicKey"`
}
