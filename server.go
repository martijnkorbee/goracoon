package goracoon

import (
	"fmt"
	"net/http"
	"time"
)

// ListenAndServe starts the webserver
func (gr *Goracoon) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", gr.config.host, gr.config.port),
		Handler:      gr.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	// defer close db and/or cache connections if open
	if gr.DB.ConnectionPool != nil {
		defer gr.DB.ConnectionPool.Close()
	}
	if redisPool != nil {
		defer redisPool.Close()
	}
	if badgerConnection != nil {
		defer badgerConnection.Close()
	}

	// defer close log file
	if !gr.Debug {
		defer logOut.Close()
	}

	// start RPC server
	go gr.listenRPC()

	gr.Log.Info().Msg(fmt.Sprintf("Starting webserver on %s:%s", gr.config.host, gr.config.port))
	err := srv.ListenAndServe()
	gr.Log.Fatal().Err(err)

	return nil
}
