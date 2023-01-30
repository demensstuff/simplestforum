.DEFAULT_GOAL := default
-include deployment/.env

# Tool urls
MIGRATE_TOOL_URL := https://github.com/golang-migrate/migrate/releases/download/v4.13.0/migrate.linux-amd64.tar.gz
LINTER_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

# Default command to cleanup the code
default: gen format lint-fix

### Build and run ###

# Compiles the application
build-local:
	go build -o ./build/simplestforum ./cmd/simplestforum

# Builds docker image
build-docker:
	docker build -f ./deployment/Dockerfile -t simplestforum:1.0 .

# Runs the build application
run:
	go run ./cmd/simplestforum

### Database up/down (via docker) ###

# Deploys the database
up-db:
	docker-compose -f ./deployment/dc.local.yml --env-file ./deployment/.env up -d database

# Shuts down the database
down-db:
	cd deployment && \
    docker-compose -f dc.local.yml down

### Code generation ###

# Generates graphql resolvers
gen:
	go generate ./...

### Linter ###

# Installs golangci-lint
lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.0

# Inspects the code
lint:
	golangci-lint run

# Inspects the code and automatically solves problems if possible
lint-fix:
	golangci-lint run --fix

### Migrations ###

# Installs migrate tool
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1

# Creates a new migration
migrate-new:
	migrate create -ext sql -dir ./migrations "$(name)"

# Forcefully sets the desired migrations version
migrate-force:
	migrate -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" -path ./migrations force $(v)

# Applies the migrations
migrate-up:
	migrate -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" -path ./migrations up

# Rollbacks last migration
migrate-down:
	migrate -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" -path ./migrations down 1

# Purges the database
migrate-drop:
	migrate -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" -path ./migrations drop -f

# Drops the database and creates a fresh new version
migrate-rebase: migrate-drop migrate-up

# Auto-formats code
format:
	gofmt -s -w . && \
	go vet ./... && \
	go mod tidy