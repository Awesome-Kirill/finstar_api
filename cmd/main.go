package main

import (
	"context"
	"finstar/config"
	"finstar/internal/data"
	"finstar/internal/transport"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	conf := config.Get()
	log = log.With().Str("service", "finstar-api").Logger()
	conn, err := pgx.Connect(context.Background(), conf.PostgresConn)
	//db, err := sql.Open("postgres", config.PostgresConn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed connection db")
	}
	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed ping db")
		return
	}

	// toDO
	err = data.RunMigrations()
	if err != nil {
		log.Fatal().Err(err).Msg("failed up migrate")
	}
	log.Info().Msg("up migrate")

	server := transport.NewHttp(log, data.NewDbRepository(conn))

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
	err = conn.Close(ctx)
	if err != nil {
		log.Error().Err(err).Msg("db close error")
	}
	log.Info().Msg("db close")
	log.Info().Msg("stopped!")
}
