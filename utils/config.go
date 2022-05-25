package utils

type Config struct {
	Server  Server  `yaml:"server"`
	Storage Storage `yaml:"storage"`
}

type Server struct {
	WebHostport  string `yaml:"webHostport"`
	GrpcHostport string `yaml:"grpcHostport"`
	Credentials  string `yaml:"credentials"`
}

// Exactly one of the storage fields may be supplied.
type Storage struct {
	Sqllite sqlLite `yaml:"sqllite"`
}

type sqlLite struct {
	Fpath string `yaml:"fpath"`
}
