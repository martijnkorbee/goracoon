package goracoon

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
)

func (gr *Goracoon) routes() http.Handler {
	mux := chi.NewRouter()

	// default middleware
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if gr.Debug {
		mux.Use(httplog.RequestLogger(*gr.Log))
	}
	mux.Use(middleware.Recoverer)
	mux.Use(gr.CheckMaintenanceMode)

	// added middleware
	mux.Use(gr.SessionLoad)
	mux.Use(gr.NoSurf)

	// default home handler
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "silence is golden.")
	})

	return mux
}
