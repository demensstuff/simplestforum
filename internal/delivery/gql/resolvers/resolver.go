package resolvers

//go:generate go run github.com/99designs/gqlgen generate

import (
	"net/http"
	"simplestforum/internal/delivery/gql"
	"simplestforum/internal/delivery/gql/directives"
	"simplestforum/internal/delivery/gql/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Interactors is a list of all abstract usecases.
type Interactors struct {
	User         UserInteractor
	Topic        TopicInteractor
	Section      SectionInteractor
	Post         PostInteractor
	Notification NotificationInteractor
}

type Resolver = Interactors

// NewGQLHandler sets up the GraphQL schema, resolvers and directives.
func NewGQLHandler(interactors *Interactors) http.Handler {
	srv := handler.NewDefaultServer(
		gql.NewExecutableSchema(
			gql.Config{
				Resolvers: interactors,

				Directives: gql.DirectiveRoot{
					Range:     directives.Range,
					Normalise: directives.Normalise,
				},
			},
		),
	)

	srv.SetRecoverFunc(middleware.Recover)
	srv.AroundResponses(middleware.WrapResponse)
	srv.AroundFields(middleware.CollectFields)

	return srv
}
