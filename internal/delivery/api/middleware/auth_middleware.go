package middleware

import (
	"net/http"
	"simplestforum/internal/delivery"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/domain/usecase"
)

// Auth represents an Auth middleware.
type Auth struct {
	userService usecase.UserAdapter
}

// NewAuth instantiates an Auth middleware.
func NewAuth(userService usecase.UserAdapter) *Auth {
	return &Auth{userService: userService}
}

// Handler creates a new callback that is run to check the credentials.
func (m *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			nickname, password, ok := r.BasicAuth()

			if !ok {
				next.ServeHTTP(w, r)

				return
			}

			ctx := r.Context()
			sess := entity.GetSession(ctx)
			user, err := m.userService.ByLoginAndPassword(sess, nickname, password)

			if err != nil {
				_, _ = w.Write(delivery.BuildErrorResponse(sess, err, true))

				return
			}

			sess.UserID = user.ID
			sess.Level = user.Level
			sess.Restriction = user.Restriction
			sess.Ctx = ctx

			r = r.WithContext(entity.PutSession(ctx, sess))

			next.ServeHTTP(w, r)
		},
	)
}
