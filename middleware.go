package goracoon

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (gr *Goracoon) SessionLoad(next http.Handler) http.Handler {
	return gr.SessionManager.LoadAndSave(next)
}

func (gr *Goracoon) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(gr.config.cookie.secure)

	csrfHandler.ExemptRegexp("/api/.*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   gr.config.cookie.domain,
	})

	return csrfHandler
}
