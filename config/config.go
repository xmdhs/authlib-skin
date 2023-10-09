package config

type Config struct {
	OfflineUUID bool   `yaml:"offlineUUID"`
	Port        string `yaml:"port"`
	Log         struct {
		Level string `yaml:"level"`
		Json  bool   `yaml:"json"`
	} `yaml:"log"`
	Sql struct {
		MysqlDsn string `yaml:"mysqlDsn"`
	} `yaml:"sql"`
	Debug bool `yaml:"debug"`
	Cache struct {
		Type string `yaml:"type"`
		Ram  int    `yaml:"ram"`
	} `yaml:"cache"`
	RaelIP         bool   `yaml:"raelIP"`
	MaxIpUser      int    `yaml:"maxIpUser"`
	RsaPriKey      string `yaml:"rsaPriKey"`
	TexturePath    string `yaml:"texturePath"`
	TextureBaseUrl string `yaml:"textureBaseUrl"`
	WebBaseUrl     string `yaml:"webBaseUrl"`
	ServerName     string `yaml:"serverName"`

	Captcha struct {
		Type    string `yaml:"type"`
		SiteKey string `yaml:"siteKey"`
		Secret  string `yaml:"ecret"`
	} `yaml:"captcha"`
}
