package transport

import (
	"context"
	apierror "finstar/internal/error"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type TransferRequest struct {
	From  int     `json:"from"`
	To    int     `json:"to"`
	Total float32 `json:"total"`
}

type DepositRequest struct {
	To    int     `json:"to"`
	Total float32 `json:"total"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type Repository interface {
	Deposited(ctx context.Context, userId int, total float32) error
	Transfer(ctx context.Context, userIdFrom int, userIdTo int, total float32) error
	FindUser(ctx context.Context, userId int) (bool, error)
}

func (h *HTTP) Deposit(context *gin.Context) {
	request := DepositRequest{}
	err := context.ShouldBindJSON(&request)

	if err != nil {
		log.Error().Err(err).Msg("bind request error")
		context.JSON(http.StatusBadRequest, apierror.APIError{Message: err.Error(), ExtCode: apierror.BindingError})
		return

	}

	ok, err := h.repository.FindUser(context.Request.Context(), request.To)
	if err != nil {
		log.Error().Err(err).Msg("find user error")
		context.JSON(http.StatusInternalServerError, apierror.APIError{Message: err.Error()})
		return
	}

	if !ok {
		context.JSON(http.StatusBadRequest, apierror.APIError{ExtCode: apierror.NotFound, Message: "User not found"})
		return
	}

	err = h.repository.Deposited(context.Request.Context(), request.To, request.Total)
	if err != nil {
		log.Error().Err(err).Msg("Deposited error")
		context.JSON(http.StatusInternalServerError, apierror.APIError{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, SuccessResponse{Success: true})

}

func (h *HTTP) Transfer(context *gin.Context) {
	request := TransferRequest{}
	err := context.ShouldBindJSON(&request)

	if err != nil {
		log.Error().Err(err).Msg("bind request error")
		context.JSON(http.StatusBadRequest, apierror.APIError{Message: err.Error(), ExtCode: apierror.BindingError})
		return
	}

	ok, err := h.repository.FindUser(context.Request.Context(), request.To)
	if err != nil {
		log.Error().Err(err).Msg("find user error")
		context.JSON(http.StatusInternalServerError, apierror.APIError{Message: err.Error()})
		return
	}

	if !ok {
		context.JSON(http.StatusBadRequest, apierror.APIError{ExtCode: apierror.NotFound, Message: "User not found"})
		return
	}

	ok, errs := h.repository.FindUser(context.Request.Context(), request.To)
	if errs != nil {
		log.Error().Err(err).Msg("find user error")
		context.JSON(http.StatusInternalServerError, apierror.APIError{Message: err.Error()})
		return
	}

	if !ok {
		context.JSON(http.StatusBadRequest, apierror.APIError{ExtCode: apierror.NotFound, Message: "User not found"})
		return
	}

	err = h.repository.Transfer(context.Request.Context(), request.From, request.To, request.Total)
	if err != nil {
		if err == apierror.LowBalance {
			context.JSON(http.StatusBadRequest, apierror.APIError{Message: err.Error(), ExtCode: apierror.NotEnoughMoney})
			return
		}
		log.Error().Err(err).Msg("Transfer error")
		context.JSON(http.StatusInternalServerError, apierror.APIError{Message: err.Error()})
		return
	}
	context.JSON(http.StatusOK, SuccessResponse{Success: true})
}
