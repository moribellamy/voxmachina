package utils

type sqlLite struct {
	Fpath string
}

type cache struct {
	WebHostport  string
	GrpcHostport string
}

type Config struct {
	Cache       cache
	Credentials string
	Sqllite     sqlLite
}
