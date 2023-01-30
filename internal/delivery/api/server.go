package api

import (
	"context"
	"log"
	"net/http"
	assets "simplestforum"
	"simplestforum/internal/delivery/api/middleware"

	"github.com/gorilla/mux"
)

const (
	gqlEndpoint        = "/v1/public"
	playgroundEndpoint = "/playground"
)

// Server is a structure which contains everything needed for the REST server.
type Server struct {
	srv    *http.Server
	router *mux.Router

	gqlHandler http.Handler

	middleware *middleware.Middlewares
}

// NewServer instantiates a new Server object.
func NewServer(port string, gh http.Handler, m *middleware.Middlewares) *Server {
	r := mux.NewRouter()

	srv := Server{
		srv: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
		router:     r,
		gqlHandler: gh,
		middleware: m,
	}

	return &srv
}

// setGraphQLRoutes defines the GraphQL endpoint.
func (srv *Server) setGraphQLRoutes() {
	srv.router.Handle(gqlEndpoint, srv.gqlHandler)
}

// setMiscRoutes defines miscellaneous helpful routes.
func (srv *Server) setMiscRoutes() {
	srv.router.HandleFunc(playgroundEndpoint, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, _ = w.Write(assets.GQLPlaygroundHTML)
	})
}

// Start prepares all routes required and listens to the incoming requests.
func (srv *Server) Start() error {
	srv.router.Use(srv.middleware.Handlers()...)
	srv.setGraphQLRoutes()
	srv.setMiscRoutes()

	// Preparing the GQL Playground
	assets.InitGQLPlaygroundHTML([]byte(srv.srv.Addr), []byte(gqlEndpoint))

	log.Println("Starting the server at", srv.srv.Addr)
	log.Println("Check the playground at http://localhost" + srv.srv.Addr + playgroundEndpoint)

	return srv.srv.ListenAndServe()
}

// Shutdown stops the server.
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.srv.Shutdown(ctx)
}
