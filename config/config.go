package config

type Config struct {
	OfflineUUID    bool        `toml:"offlineUUID" comment:"为 true 则 uuid 生成方式于离线模式相同，若从离线模式切换不会丢失数据。\n已有用户数据的情况下勿更改此项"`
	Port           string      `toml:"port"`
	Log            Log         `toml:"log"`
	Sql            Sql         `toml:"sql"`
	Debug          bool        `toml:"debug" comment:"输出每条执行的 sql 语句"`
	Cache          Cache       `toml:"cache"`
	RaelIP         bool        `toml:"raelIP" comment:"位于反向代理后启用，用于记录真实 ip\n若直接提供服务，请勿打开，否则会被伪造 ip"`
	MaxIpUser      int         `toml:"maxIpUser" comment:"ip 段最大注册用户，ipv4 为 /24 ipv6 为 /48"`
	RsaPriKey      string      `toml:"rsaPriKey,multiline" comment:"运行后勿修改，若为集群需设置为一致"`
	TexturePath    string      `toml:"texturePath" comment:"材质文件保存路径，如果需要对象存储可以把对象储存挂载到本地目录上"`
	TextureBaseUrl string      `toml:"textureBaseUrl" comment:"材质静态文件提供基础地址\n如果静态文件位于 oss 上，比如 https://s3.amazonaws.com/example/1.png\n则填写 https://s3.amazonaws.com/example \n若通过反向代理提供服务并启用了 https，请在在此处填写带有 https 的基础路径，否则游戏内无法加载皮肤"`
	WebBaseUrl     string      `toml:"webBaseUrl" comment:"用于在支持的启动器中展示本站的注册地址\n填写类似 https://example.com"`
	ServerName     string      `toml:"serverName" comment:"皮肤站名字，用于在多个地方展示"`
	Captcha        Captcha     `toml:"captcha"`
	Email          EmailConfig `toml:"email"`
}

type Log struct {
	Level string `toml:"level"`
	Json  bool   `toml:"json" comment:"json 格式输出"`
}

type Sql struct {
	DriverName string `toml:"driverName" comment:"可填 mysql 或 sqlite3"`
	Dsn        string `toml:"dsn" comment:"填写见 mysql https://github.com/go-sql-driver/mysql#dsn-data-source-name\nsqlite https://github.com/mattn/go-sqlite3\n例如 mysql 用户名:密码@tcp(mysqlIP:端口)/数据库名\nsqlite data.db?_txlock=IMMEDIATE&_journal_mode=WAL&_fk=true"`
}

type Cache struct {
	Type     string `toml:"type" comment:"默认使用内存缓存，若需要集群部署，填写 redis"`
	Ram      int    `toml:"ram" comment:"内存缓存使用大小，单位 b"`
	Addr     string `toml:"addr" comment:"redis 服务端地址，如 127.0.0.1:6379"`
	Password string `toml:"password" comment:"redis 密码"`
}

type Captcha struct {
	Type    string `toml:"type" comment:"验证码类型，目前只支持 cloudflare turnstile，若需要填写 turnstile"`
	SiteKey string `toml:"siteKey"`
	Secret  string `toml:"secret"`
}

type EmailConfig struct {
	Enable      bool       `toml:"enable" comment:"注册验证邮件，且允许使用邮箱找回账号"`
	Smtp        []SmtpUser `toml:"smtp"`
	AllowDomain []string   `toml:"allow_domain" comment:"允许用于注册的邮箱域名，留空则允许全部"`
	EmailReg    string     `toml:"email_reg" comment:"邮箱正则，留空则不处理，如 ^[0-9]+@qq.com$|^[^+\\.A-Z]+@gmail.com$"`
	EmailRegMsg string     `toml:"email_reg_msg" comment:"不满足要求时的提示信息"`
}

type SmtpUser struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	SSL  bool   `toml:"SSL" comment:"启用 ssl"`
	Name string `toml:"name"`
	Pass string `toml:"password"`
}

func Default() Config {
	return Config{
		OfflineUUID: true,
		Port:        ":8080",
		Log: Log{
			Level: "debug",
			Json:  false,
		},
		Sql: Sql{
			DriverName: "sqlite3",
			Dsn:        "data.db?_txlock=IMMEDIATE&_journal_mode=WAL&_fk=true",
		},
		Debug: false,
		Cache: Cache{
			Type:     "",
			Ram:      10000000,
			Addr:     "",
			Password: "",
		},
		RaelIP:         false,
		MaxIpUser:      0,
		RsaPriKey:      "",
		TexturePath:    "",
		TextureBaseUrl: "",
		WebBaseUrl:     "",
		ServerName:     "没有设置名字",
		Captcha:        Captcha{},
		Email: EmailConfig{
			Smtp: []SmtpUser{
				{
					Host: "",
					Port: 0,
					SSL:  false,
					Name: "",
					Pass: "",
				},
			},
			AllowDomain: []string{},
		},
	}
}
