package transport

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type HTTP struct {
	srv        *http.Server
	repository Repository
	log        zerolog.Logger
}

func (h *HTTP) Start() error {
	return h.srv.ListenAndServe()
}

func (h *HTTP) Stop(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}

func NewHttp(options Options) *HTTP {
	h := &HTTP{}
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/user/deposit", h.Deposit)
	api.POST("/user/transfer", h.Transfer)

	h.srv = &http.Server{
		Addr:    options.Addr,
		Handler: router,
	}

	h.repository = options.Repository
	h.log = options.Log
	return h
}
