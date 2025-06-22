package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/api"
	"github.com/mtavano/admoai-takehome/internal/metrics"
	"github.com/mtavano/admoai-takehome/internal/store"
	_ "github.com/mtavano/admoai-takehome/migrations"
	"github.com/pressly/goose/v3"
)

func main() {
	fmt.Println("admoai-take-home-test initialization")

	// Initialize metrics collector
	metrics.Init()

	// Initialize database
	dbDriver := os.Getenv("DB_DRIVER")
	dbDSN := os.Getenv("DB_DSN")

	if dbDriver == "" {
		dbDriver = "sqlite3"
	}
	if dbDSN == "" {
		dbDSN = "./data/admoai.db"
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		panic(fmt.Sprintf("Failed to create data directory: %v", err))
	}

	// Initialize database store
	dbStore, err := store.NewSqlStore(dbDriver, dbDSN)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// Run migrations
	if err := goose.SetDialect(dbDriver); err != nil {
		panic(fmt.Sprintf("Failed to set dialect: %v", err))
	}

	if err := goose.Up(dbStore.DB.DB, "./migrations"); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	fmt.Println("Database initialized and migrations completed")

	// api server specifics
	apiCtx := &api.Context{
		Db: dbStore,
	}
	router := gin.Default()
	api.RegisterRoutes(apiCtx, router)

	port := os.Getenv("API_PORT")

	// Configure and execute the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println(fmt.Sprintf("Start listening now on port %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-sigs
}
