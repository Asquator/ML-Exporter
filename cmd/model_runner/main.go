package main

import (
	"fmt"
	"goml/internal/config"
	mwLogger "goml/internal/http-server/middleware/logger"
	"goml/internal/http-server/tus/hooks"
	"goml/internal/logger/sl"
	"goml/internal/storage/local"
	"log/slog"
	"net/http"
	"os"

	ort "github.com/yalue/onnxruntime_go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"

	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(&cfg)
	log.Info("service started")
	log.Debug("debug messages are enabled")
	log.Debug("config loaded", slog.Any("config", cfg))

	_, err := local.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	store := filestore.FileStore{
		Path: cfg.Tus.UploadPath,
	}

	filelocker := filelocker.FileLocker{
		Path: cfg.Tus.UploadPath,
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
			hooks.CompleteUploadHook(event, &cfg)
		}
	}()

	if err != nil {
		log.Error("failed to create handler", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Mount("/upload/", http.StripPrefix("/upload/", handler))

	srv := &http.Server{
		Addr:           cfg.Address,
		Handler:        router,
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("starting server", slog.Attr{Key: "config", Value: slog.AnyValue(srv)})

	/*
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
	*/

	log.Error("server stopped")

	ort.SetSharedLibraryPath("/usr/lib/libonnxruntime.so")

	err = ort.InitializeEnvironment()

	if err != nil {
		panic(err)
	}
	defer ort.DestroyEnvironment()

	inp := []float32{-0.051237, -0.829306, 0.072428, -0.510460, -0.586643, -0.251280, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
		1.235929, -0.829306, -0.605820, 1.620986, -0.462119, 2.005455, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
		0.216922, -0.764031, -0.605820, -0.943410, -2.143185, 0.490219, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	shape := ort.NewShape(3, 20)

	tensor, err := ort.NewTensor(shape, inp)

	defer tensor.Destroy()

	outputTensor, err := ort.NewEmptyTensor[int64](ort.NewShape(3))
	defer outputTensor.Destroy()

	if err != nil {
		panic(err)
	}

	tensor.GetInternals()

	meta, err := ort.GetModelMetadata("/home/acf/Programming/goml/cmd/model_runner/models/logistic.onnx")

	fmt.Println(meta.GetCustomMetadataMapKeys())
	fmt.Println(meta.GetDescription())
	fmt.Println(meta.GetVersion())
	i, o, err := ort.GetInputOutputInfo("/home/acf/Programming/goml/cmd/model_runner/models/logistic.onnx")

	fmt.Println(i, o)

	session, err := ort.NewAdvancedSession("/home/acf/Programming/goml/cmd/model_runner/models/logistic.onnx",
		[]string{"float_input"}, []string{"output_label"},
		[]ort.ArbitraryTensor{tensor}, []ort.ArbitraryTensor{outputTensor}, nil)
	defer session.Destroy()

	session.Run()

	if err != nil {
		panic(err)
	}

	output := outputTensor.GetData()

	fmt.Println(output)

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
