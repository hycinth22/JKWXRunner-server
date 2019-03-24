package config

const (
	Release = false
	Debug   = !Release

	ListenAddr = "localhost:8066"

	DBLogMode     = true
	DSN           = `.\app_new3.db?mode=rw`
	DBMaxIdleConn = 1
	DBMaxOpenConn = 10
)
