package http

import (
	"context"
	"eau-de-go/internal/transport/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router         *mux.Router
	AppUserService AppUserService
	AuthService    AuthService
	Server         *http.Server
}

func NewHandler(appUserService AppUserService) *Handler {
	h := &Handler{
		AppUserService: appUserService,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(middleware.JSONMiddleware)
	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {

	h.Router.HandleFunc("/api/user/login/", h.Login).Methods("POST")
	h.Router.HandleFunc("/api/user/token-refresh/", h.TokenRefresh).Methods("POST")
	h.Router.HandleFunc("/api/user/{id}/", h.GetAppUserById).Methods("GET")
	h.Router.HandleFunc("/api/user/", h.CreateAppUser).Methods("POST")
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shut down gracefully")
	return nil
}
