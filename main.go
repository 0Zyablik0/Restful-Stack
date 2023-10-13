package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"restful_stack/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "1.0.0"

func main() {
	cfg := configRead()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := setupDatabase(&cfg)
	if err != nil {
		panic(err)
	}

	app := &app{config: cfg, logger: logger, db: db}

	router := app.router()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	logger.Printf("%s: Starting server on %s\n", cfg.environment, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func configRead() config {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "port of the application")
	flag.StringVar(&cfg.environment, "env", "dev", "Environment: (dev|prod)")
	flag.BoolVar(&cfg.enableLogging, "enable-logging", false, "Enable logging")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("RESTFUL_STACK_DB_DSN"), "PostgreSQL DSL")
	flag.Parse()
	return cfg
}

func setupDatabase(cfg *config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.StackEntity{}, &models.UserStacks{}, &models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
