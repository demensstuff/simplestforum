# Simplest forum developed in Golang

This is a tiny example (or even POC) of the simplest forum engine one can imagine, which represents a clean architecture and GraphQL API. The basic idea is to demonstrate an example of circular dependency between services: this allows to handle chained GraphQL queries recursively. As long as the invocation of the services eventually depends on the requested fields, this provides quite an optimal storage usage.

The business logic is encapsulated in UseCases. All four layers (GraphQL resolvers, UseCases, Services, Repositories) communicate via interfaces.

The application uses PostgreSQL database as a storage. Both the database and the app can be run in Docker. [GQLGen](https://gqlgen.com/v0.9.3/) library is used to generate the relevant code based on the schema files.

## Quick start

1) Create the `/deployment/.env` file; `/deployment/.env.example` contains an example configuration.
2) Use `make gen` to generate a GraphQL execution runtime.
3) Use `make up-db` to launch a database container.
4) Execute `make run` to run the app. You may also use `make build-docker` to build a Docker image. Try `go mod tidy` in case any dependencies are missing.
5) Check out the GraphQL Playground at http://localhost:8080/playground (by default) to test some queries and mutations.
