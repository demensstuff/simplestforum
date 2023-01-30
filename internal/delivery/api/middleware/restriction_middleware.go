package middleware

import (
	"net/http"
	"simplestforum/internal/delivery"
	"simplestforum/internal/domain"
	"simplestforum/internal/domain/entity"
)

// Restriction represents a Restriction middleware.
type Restriction struct{}

// NewRestriction instantiates an Restriction middleware.
func NewRestriction() *Restriction {
	return &Restriction{}
}

// Handler creates a new callback that is run to check if the user is banned.
func (m *Restriction) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sess := entity.GetSession(ctx)

			if sess.IsAuthorized() && sess.Restriction.AtLeast(entity.UserRestrictionBanned) {
				_, _ = w.Write(delivery.BuildErrorResponse(sess, domain.NewError(domain.ErrCodeRestricted,
					"You are banned"), true))

				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
