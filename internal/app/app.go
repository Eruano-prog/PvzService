package app

import (
	"AvitoPvz/internal/config"
	"AvitoPvz/internal/jwt"
	"AvitoPvz/internal/postgres"
	"AvitoPvz/internal/rest/controllers"
	"AvitoPvz/internal/service"
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func Run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	tokenService := jwt.NewJWTService(log, cfg.TokenSecret)

	//Repository
	db, err := postgres.SetupDBConnection(log, cfg.Database.Address, cfg.Database.DBName, cfg.Database.Username, cfg.Database.Password)
	if err != nil {
		log.Error("error connecting to database")
		return err
	}

	userRepository := postgres.NewUserRepo(log, db)
	pvzRepository := postgres.NewPVZRepo(log, db)
	receptionRepository := postgres.NewReceptionRepo(log, db)
	productRepository := postgres.NewProductRepo(log, db)

	// Service
	userService := service.NewUserService(log, tokenService, userRepository)
	pvzService := service.NewPVZService(log, pvzRepository, receptionRepository, productRepository)
	receptionService := service.NewReceptionService(log, receptionRepository)
	productService := service.NewProductService(log, productRepository, receptionRepository)

	// Controller
	authController := controllers.NewAuthController(log, userService)
	pvzController := controllers.NewPVZController(log, pvzService)
	receptionController := controllers.NewReceptionController(log, receptionService)
	productController := controllers.NewProductController(log, productService)

	// HTTP server
	mux := http.NewServeMux()

	authController.Register(mux)
	pvzController.Register(mux)
	receptionController.Register(mux)
	productController.Register(mux)

	server := http.Server{
		Addr:        cfg.HTTPConfig.Address,
		ReadTimeout: cfg.HTTPConfig.Timeout,
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	log.Info("Running HTTP server", "address", cfg.HTTPConfig.Address)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return err
		}
	}

	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level, AddSource: true})
	return slog.New(handler)
}
