package middleware

import (
	"simplestforum/internal/domain/usecase"

	"github.com/gorilla/mux"
)

// Middlewares provides list of all middlewares.
type Middlewares struct {
	Cors        *Cors
	Session     *Session
	Auth        *Auth
	Restriction *Restriction
}

func NewMiddlewares(userAdapter usecase.UserAdapter) *Middlewares {
	return &Middlewares{
		Cors:        NewCors(),
		Session:     NewSession(),
		Auth:        NewAuth(userAdapter),
		Restriction: NewRestriction(),
	}
}

func (m Middlewares) Handlers() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{m.Cors.Handler, m.Session.Handler, m.Auth.Handler, m.Restriction.Handler}
}
