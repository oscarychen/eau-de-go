package http

import (
	"context"
	"eau-de-go/internal/settings"
	"eau-de-go/internal/transport/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router          *mux.Router
	ProtectedRouter *mux.Router
	AppUserService  AppUserService
	AuthService     AuthService
	Server          *http.Server
}

func NewHandler(appUserService AppUserService) *Handler {
	h := &Handler{
		AppUserService: appUserService,
	}
	h.Router = mux.NewRouter()
	h.ProtectedRouter = h.Router.PathPrefix("/api").Subrouter()
	h.ProtectedRouter.Use(middleware.JwtAuthMiddleware)

	h.mapRoutes()
	h.Router.Use(middleware.JSONMiddleware)

	h.Server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", settings.ServerPort),
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {

	h.Router.HandleFunc("/auth/login/", h.Login).Methods("POST")
	h.Router.HandleFunc("/auth/token-refresh/", h.TokenRefresh).Methods("POST")
	h.Router.HandleFunc("/auth/sign-up/", h.CreateAppUser).Methods("POST")

	h.ProtectedRouter.HandleFunc("/user/{id}/", h.GetAppUserById).Methods("GET")
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
