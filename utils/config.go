package utils

type sqlLite struct {
	Fpath string `yaml:"fpath"`
}

type cache struct {
	WebHostport  string `yaml:"webHostport"`
	GrpcHostport string `yaml:"grpcHostport"`
}

type Config struct {
	Cache       cache   `yaml:"cache"`
	Credentials string  `yaml:"credentials"`
	Sqllite     sqlLite `yaml:"sqllite"`
}
