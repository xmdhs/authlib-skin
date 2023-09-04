package config

type Config struct {
	OfflineUUID bool
	Port        string
	Log         struct {
		Level string
		Json  bool
	}
	Sql struct {
		MysqlDsn string
	}
	Debug  bool
	JwtKey string
	Cache  struct {
		Type string
		Ram  int
	}
	RaelIP bool
}
