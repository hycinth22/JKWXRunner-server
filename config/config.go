package config

//noinspection GoBoolExpressions,GoUnusedConst
const (
	Release = false
	Debug   = !Release

	ListenAddr = "localhost:8066"

	DBLogMode     = true
	DBLogDir      = `.\data\logs\database`
	DBLogFileName = `20060102_150405.000000000.log` // Time Package Format Layout

	DSN           = `.\data\v3_2018-2019-2_s.db?mode=rw`
	DBMaxIdleConn = 1
	DBMaxOpenConn = 10
)
