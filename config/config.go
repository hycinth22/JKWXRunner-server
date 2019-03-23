package config

const (
	ListenAddr = "localhost:8066"

	DBLogMode     = true
	DSN           = `.\app_new2.db?mode=rw`
	DBMaxIdleConn = 1
	DBMaxOpenConn = 10
)
