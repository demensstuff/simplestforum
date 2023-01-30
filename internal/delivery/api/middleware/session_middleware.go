package middleware

import (
	"net/http"
	"simplestforum/internal/domain/entity"

	"github.com/google/uuid"
)

// Session represents a Session middleware.
type Session struct{}

// NewSession instantiates a Session middleware.
func NewSession() *Session {
	return &Session{}
}

// Handler creates a session.
func (m *Session) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sess := entity.Session{
				ID:  uuid.NewString(),
				Ctx: r.Context(),
			}

			r = r.WithContext(entity.PutSession(r.Context(), sess))

			next.ServeHTTP(w, r)
		},
	)
}
