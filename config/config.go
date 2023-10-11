package config

type Config struct {
	OfflineUUID    bool    `yaml:"offlineUUID"`
	Port           string  `yaml:"port"`
	Log            Log     `yaml:"log"`
	Sql            Sql     `yaml:"sql"`
	Debug          bool    `yaml:"debug"`
	Cache          Cache   `yaml:"cache"`
	RaelIP         bool    `yaml:"raelIP"`
	MaxIpUser      int     `yaml:"maxIpUser"`
	RsaPriKey      string  `yaml:"rsaPriKey"`
	TexturePath    string  `yaml:"texturePath"`
	TextureBaseUrl string  `yaml:"textureBaseUrl"`
	WebBaseUrl     string  `yaml:"webBaseUrl"`
	ServerName     string  `yaml:"serverName"`
	Captcha        Captcha `yaml:"captcha"`
}

type Log struct {
	Level string `yaml:"level"`
	Json  bool   `yaml:"json"`
}

type Sql struct {
	DriverName string `yaml:"driverName"`
	Dsn        string `yaml:"dsn"`
}

type Cache struct {
	Type     string `yaml:"type"`
	Ram      int    `yaml:"ram"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type Captcha struct {
	Type    string `yaml:"type"`
	SiteKey string `yaml:"siteKey"`
	Secret  string `yaml:"ecret"`
}
