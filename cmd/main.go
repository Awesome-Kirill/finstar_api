package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"finstar/internal/data"
	"finstar/internal/transport"
)

type Config struct {
	//LogLevel     string `envconfig:"log_level" default:"debug"`
	HttpAPI      string `envconfig:"http_api" default:":8000"`
	PostgresConn string `envconfig:"pg_conn" default:"postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable&search_path=billing"`
}

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("failed while reading config")
		return
	}
	/*
		logLevel, err := zerolog.ParseLevel(config.LogLevel)
		if err != nil {
			log.Warn().Err(err)
			logLevel = zerolog.InfoLevel
		}


		log = log.Level(logLevel)

	*/
	log = log.With().Str("service", "finstar-api").Logger()

	db, err := sql.Open("postgres", config.PostgresConn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed connection db")
	}
	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed ping db")
		return
	}

	server := transport.NewHttp(transport.Options{
		Addr:       config.HttpAPI,
		Log:        log,
		Repository: data.NewDbRepository(db),
	})

	go func() {
		if err := server.Start(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http api listen error")
			return
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	log.Info().Msg("Interrupt!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Stop(ctx)

	if err != nil {
		log.Error().Err(err).Msg("server stop error")
	}
	log.Info().Msg("serv stop")
	err = db.Close()
	if err != nil {
		log.Error().Err(err).Msg("db close error")
	}
	log.Info().Msg("db close")
	log.Info().Msg("stopped!")
}
