package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	HttpAPI      string `envconfig:"HTTP_API" default:":8000"`
	PostgresConn string `envconfig:"PG_CONN" default:"postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"`
}

var config Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
	})
	return &config
}
