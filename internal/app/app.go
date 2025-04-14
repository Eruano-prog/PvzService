package app

import (
	"AvitoPvz/internal/config"
	pvzGrpc "AvitoPvz/internal/grpc"
	pvzGrpcpb "AvitoPvz/internal/grpc/api"
	"AvitoPvz/internal/jwt"
	"AvitoPvz/internal/monitoring"
	"AvitoPvz/internal/postgres"
	"AvitoPvz/internal/postgres/repository"
	"AvitoPvz/internal/rest/controllers"
	"AvitoPvz/internal/rest/middleware"
	"AvitoPvz/internal/service"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	tokenService := jwt.NewJWTService(log, cfg.TokenSecret, cfg.TokenTTL)

	//Repository
	db, err := postgres.SetupDBConnection(log, cfg.DBAddress)
	if err != nil {
		log.Error("error connecting to database")
		return err
	}

	err = postgres.Migrate(log, db)
	if err != nil {
		log.Error("error migrating database")
		return err
	}

	userRepository := repository.NewUserRepo(log, db)
	pvzRepository := repository.NewPVZRepo(log, db)
	receptionRepository := repository.NewReceptionRepo(log, db)
	productRepository := repository.NewProductRepo(log, db)

	metrics := monitoring.NewPrometheusMetrics()

	// Service
	userService := service.NewUserService(log, tokenService, userRepository)
	pvzService := service.NewPVZService(log, pvzRepository, receptionRepository, productRepository, metrics)
	receptionService := service.NewReceptionService(log, receptionRepository, metrics)
	productService := service.NewProductService(log, productRepository, receptionRepository, metrics)
	// Controller
	authController := controllers.NewAuthController(log, userService)
	pvzController := controllers.NewPVZController(log, pvzService)
	receptionController := controllers.NewReceptionController(log, receptionService)
	productController := controllers.NewProductController(log, productService)

	// HTTP server
	httpMetrics := monitoring.NewPrometheusHTTPMetrics()
	mux := http.NewServeMux()
	wrappedMux := middleware.Metrics(mux, httpMetrics)

	authController.Register(mux)
	pvzController.Register(mux, tokenService)
	receptionController.Register(mux, tokenService)
	productController.Register(mux, tokenService)

	httpServer := http.Server{
		Addr:        cfg.HTTPConfig.Address,
		ReadTimeout: cfg.HTTPConfig.Timeout,
		Handler:     wrappedMux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// GRPC
	grpcListener, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pvzGrpcpb.RegisterPVZServiceServer(grpcServer, pvzGrpc.NewServer(pvzService))
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		grpcServer.GracefulStop()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	errGroup, _ := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.PrometheusAddress, nil); err != nil {
			log.Error("failed to start metrics server", "error", err)
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		log.Info("starting gRPC server", "address", cfg.GrpcAddress)
		if err := grpcServer.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC server error: %w", err)
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		log.Error("server error", "error", err)
		return err
	}

	log.Info("server shutdown completed")

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
