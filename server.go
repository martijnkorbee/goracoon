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
		ErrorLog:     gr.ErrorLog,
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

	// start mail channels
	go gr.Mail.ListenForMail()

	// start RPC server
	go gr.listenRPC()

	gr.InfoLog.Printf("Starting webserver on %s:%s", gr.config.host, gr.config.port)
	err := srv.ListenAndServe()
	gr.ErrorLog.Fatal(err)

	return nil
}
