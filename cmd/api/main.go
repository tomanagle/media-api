package main

import (
	"context"
	"fmt"
	"go-media/internal/db"
	"go-media/internal/handlers"
	"go-media/internal/pkg/config"
	"go-media/internal/storage/s3"
	"go-media/internal/store/dbstore"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

var defaultConfig = config.Config{
	Port:                8080,
	Host:                "localhost",
	DBConnectionTimeout: 5 * time.Second,
	DBName:              "go-media",

	S3BucketName: "go-media",
	S3AWSRegion:  "us-east-1",
}

func main() {

	conf := config.Must(defaultConfig)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: conf.LogLevel,
	}))

	r := chi.NewRouter()

	server := &http.Server{
		Addr:         conf.Host + ":" + fmt.Sprintf("%d", conf.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := db.Must(ctx, db.CreateMongoClientParams{
		ConnectionString:  conf.DBConnectionString,
		ConnectionTimeout: conf.DBConnectionTimeout,
		Logger:            logger,
	})

	defer db.Disconnect(context.Background())

	mediaCollection := db.Database(conf.DBName).Collection("media")
	mediaStore := dbstore.NewMediaStore(dbstore.NewMediaStoreParams{
		Collection: mediaCollection,
	})

	s3 := s3.New(s3.NewS3Params{
		BucketName: conf.S3BucketName,
		AWSRegion:  conf.S3AWSRegion,
	})

	r.Post("/upload", handlers.NewUploadFileHandler(handlers.NewUploadFileHandlerParams{
		S3:         s3,
		MediaStore: mediaStore,
	}).ServeHTTP)

	// handlers
	r.Get("/healthcheck", handlers.NewHealthCheckHandler().ServeHTTP)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("error starting server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("ready for work")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("shutting down server")
	// give the server 10 seconds to shut down
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("error shutting down server", slog.Any("error", err))
		os.Exit(1)
	}
}
