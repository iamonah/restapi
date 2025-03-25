package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/restapi/internal/middleware"
)

type Handler struct {
	Service CommentService
	Router  *mux.Router
	Server  *http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.Router.Use(middleware.JSONMiddle)
	h.Router.Use(middleware.LoginMiddle,middleware.TimeoutMiddleware)

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080", 
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

	h.Router.HandleFunc("/api/v1/comment", h.PostComment).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/comment/{id}", h.GetComment).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/comment/{id}", h.DeleteComment).Methods(http.MethodDelete)
	h.Router.HandleFunc("/api/v1/comment/{id}", h.UpdateComment).Methods(http.MethodPut)

}

func (h *Handler) Serve() error {
	shutDown := make(chan error)
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			shutDown <- err
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Printf("%s, signal recieved shutting down server\n", (<-quit).String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	h.Server.Shutdown(ctx)

	err := <-shutDown
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	fmt.Println("Server shutdown complete")
	return nil
}
