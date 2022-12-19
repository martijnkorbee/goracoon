package goracoon

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func (gr *Goracoon) CheckMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maintenanceMode {

			allowedURLS := strings.Split(os.Getenv("ALLOWED_URLS"), ",")
			allowedURLS = append(allowedURLS, "public/maintenance.html")

			allowed := func() bool {
				for _, url := range allowedURLS {
					if strings.Contains(r.URL.Path, url) {
						return true
					}
				}
				return false
			}()

			if !allowed {
				if strings.Contains(r.URL.Path, "/api/") {
					_ = gr.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
						"message": "under maintenance",
					})
					return
				}

				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", gr.RootPath))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
