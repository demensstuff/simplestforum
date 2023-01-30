package middleware

import (
	"net/http"
)

// Cors represents a CORS middleware.
type Cors struct{}

// NewCors instantiates a CORS middleware.
func NewCors() *Cors {
	return &Cors{}
}

// Handler creates a new callback that is run before handling any requests.
func (m *Cors) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			header := w.Header()

			header.Set("Accept", "application/json")
			header.Set("Content-Type", "application/json")
			header.Set("Access-Control-Allow-Origin", "*")

			if r.Method == http.MethodOptions {
				header.Set("Access-Control-Allow-Headers", "*")
			}

			next.ServeHTTP(w, r)
		},
	)
}
