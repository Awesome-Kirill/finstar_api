package test

import (
	"bytes"
	"context"
	"errors"
	mock_transport "finstar/internal/mock"
	"finstar/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDeposit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	r := gin.Default()
	repository := mock_transport.NewMockRepository(ctrl)
	// mocking//
	repository.EXPECT().FindUser(context.TODO(), 1).Return(true, nil)
	repository.EXPECT().Deposited(context.TODO(), 1, float32(2.0)).Return(nil)

	h := transport.NewHttp(transport.Options{
		Addr:       ":8080",
		Log:        log,
		Repository: repository,
	})
	r.POST("/user/deposit", h.Deposit)
	var body = []byte(`{"to" : 1,"total" : 2}`)
	req, err := http.NewRequest(http.MethodPost, "/user/deposit", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	repository.EXPECT().FindUser(context.TODO(), 11111).Return(false, nil)
	body = []byte(`{"to" : 11111,"total" : 2}`)
	req, err = http.NewRequest(http.MethodPost, "/user/deposit", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	repository.EXPECT().FindUser(context.TODO(), 11111).Return(false, errors.New(""))
	body = []byte(`{"to" : 11111,"total" : 2}`)
	req, err = http.NewRequest(http.MethodPost, "/user/deposit", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}
