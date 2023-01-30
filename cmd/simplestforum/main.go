package main

import (
	"simplestforum/internal/delivery/api"
	"simplestforum/internal/delivery/api/middleware"
	"simplestforum/internal/delivery/gql/resolvers"
	"simplestforum/internal/domain/service"
	"simplestforum/internal/domain/usecase"
	"simplestforum/internal/infrastructure/repository"

	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simplestforum/internal/bootstrap"
	"syscall"
	"time"

	"github.com/gocraft/dbr"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const shutdownTimeout = 5 * time.Second

const pathToMigrations = "migrations"

func main() {
	// Showing timestamps in the log
	log.SetFlags(log.Lmsgprefix | log.LstdFlags)

	// Loading the configuration
	c, err := bootstrap.NewConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	log.Println("Configuration loaded")

	// Listening to the OS interruptions (Ctrl + C)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Initializing the DB connection pool, deferring release function
	dbPool, err := bootstrap.NewDBConnPg(c.DB.Username, c.DB.Password, c.DB.Name, c.DB.Host, c.DB.Port)
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}

	log.Println("Connected to the database")

	defer func(dbPool *dbr.Connection) {
		err := dbPool.Close()
		if err != nil {
			log.Println("Error closing the database connection pool:", err)
		} else {
			log.Println("Database connection pool closed")
		}
	}(dbPool)

	// Applying the migrations
	anyMigrations, err := bootstrap.UpMigrationsPg(dbPool.DB, c.DB.Name, pathToMigrations)
	if err != nil {
		log.Println("Error processing migrations:", err)

		return
	}

	if anyMigrations {
		log.Println("New migrations were applied")
	} else {
		log.Println("No new migrations found")
	}

	// Initializing the layers
	storage := repository.NewRepository(dbPool)
	adapters := service.NewServices(storage)
	interactors := usecase.NewAdapters(adapters)
	middlewares := middleware.NewMiddlewares(adapters.User)
	gqlHandler := resolvers.NewGQLHandler(interactors)

	// Creating the server
	srv := api.NewServer(
		c.HTTPPort,
		gqlHandler,
		middlewares,
	)

	// Running the server and handling the possible error
	go func() {
		err := srv.Start()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()
	<-done

	log.Println("Shutting down")

	// Giving some time for a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutting down
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println("Error on server shutdown:", err)
	}
}
