package config

import (
	"flag"
)

var (
	Port = flag.Int("port", 7007, "Victord running on default port")
	Host = flag.String("host", "127.0.0.1", "Victord running on default host")
)

type Config struct {
	Host     string
	Port     string
	ApiKey   string
	Username string
	Password string
}
