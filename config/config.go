package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	//LogLevel     string `envconfig:"log_level" default:"debug"`
	HttpAPI      string `envconfig:"http_api" default:":8000"`
	PostgresConn string `envconfig:"pg_conn" default:"postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable&search_path=billing"`
	ConnMig      string `envconfig:"pg_conn" default:"postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"`
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
