package main

import (
	"fmt"
	"goml/internal/config"
	"goml/internal/http-server/tus/hooks"
	"goml/internal/logger/sl"
	"goml/internal/storage/local"
	"log/slog"
	"net/http"
	"os"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"

	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(&cfg)
	log.Info("service started")
	log.Debug("debug messages are enabled")

	_, err := local.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	store := filestore.FileStore{
		Path: cfg.StoragePath,
	}

	filelocker := filelocker.FileLocker{
		Path: cfg.StoragePath,
	}

	composer := tusd.NewStoreComposer()

	filelocker.UseIn(composer)
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:                "/upload/",
		NotifyCompleteUploads:   true,
		StoreComposer:           composer,
		PreUploadCreateCallback: hooks.PreuploadHook,
	})

	go func() {
		for {
			event := <-handler.CompleteUploads
			hooks.CompleteUploadHook(event, cfg.StoragePath)
		}
	}()

	if err != nil {
		log.Error("failed to create handler", sl.Err(err))
		os.Exit(1)
	}

	http.Handle("/upload/", http.StripPrefix("/upload/", handler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("unable to listen: %s", err))
	}

}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	log = log.With(slog.String("env", cfg.Env))
	return log
}
