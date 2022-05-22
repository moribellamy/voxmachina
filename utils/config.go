package utils

type sqlLite struct {
	Fpath string
}

type Config struct {
	Hostport    string
	Credentials string
	Sqllite     sqlLite
}
