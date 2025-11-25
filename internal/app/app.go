package app

import (
	"clean-architecture/config"
	inboundadapterecho "clean-architecture/internal/adapter/inbound/echo"
	outboundadapterkafka "clean-architecture/internal/adapter/outbound/kafka"
	outboundadapterminio "clean-architecture/internal/adapter/outbound/minio"
	outboundadapterpostgres "clean-architecture/internal/adapter/outbound/postgres/repository"
	"clean-architecture/internal/domain/service"
	"clean-architecture/utils/validator"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	cfg := config.NewConfig()
	redisConfig := cfg.RedisConfig()
	kafkaConfig := cfg.NewKafkaConfig()

	initMinio, err := cfg.InitMinio()
	if err != nil {
		log.Fatalf("[RunServer-1] Failed to connect Minio: %v", err)
		return
	}

	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-2] Failed to connect Postgres: %v", err)
		return
	}

	if cfg.App.AppPort == "" {
		cfg.App.AppPort = os.Getenv("APP_PORT")
	}
	appPort := ":" + cfg.App.AppPort

	publisher, err := outboundadapterkafka.NewKafkaProducer(cfg.Kafka.Brokers, kafkaConfig)
	if err != nil {
		log.Fatalf("[RunServer-3] Failed to init Kafka: %v", err)
	}

	minioClient := outboundadapterminio.NewMinioStorage(initMinio, cfg.Minio.Bucket)
	userRepo := outboundadapterpostgres.NewUserRepository(db.DB)
	verificationTokenRepo := outboundadapterpostgres.NewVerificationTokenRepository(db.DB)
	roleRepo := outboundadapterpostgres.NewRoleRepository(db.DB)

	jwtService := service.NewJwtService(cfg)
	kafkaService := service.NewKafkaService(cfg, publisher)
	userService := service.NewUserService(userRepo, cfg, jwtService, verificationTokenRepo, kafkaService, redisConfig)
	roleService := service.NewRoleService(roleRepo)

	e := echo.New()
	e.Use(middleware.CORS())
	e.HideBanner = true
	e.Use(middleware.Recover())

	customValidator := validator.NewValidator(db.DB)
	if err := en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator); err != nil {
		log.Fatalf("[RunServer-4] %v", err)
		return
	}
	e.Validator = customValidator

	mid := inboundadapterecho.NewMiddlewareAdapter(cfg, redisConfig, jwtService)

	pingHandler := inboundadapterecho.NewPingHandler()
	userHandler := inboundadapterecho.NewUserHandler(userService)
	roleHandler := inboundadapterecho.NewRoleHandler(roleService)
	uploadImageHandler := inboundadapterecho.NewUploadImageHandler(minioClient)

	inboundadapterecho.InitRoutes(e, mid, pingHandler, userHandler, roleHandler, uploadImageHandler)

	go func() {
		log.Infof("[RunServer-5] Server starting at %s", appPort)
		if err := e.Start(appPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[RunServer-6] Server start failed: %v", err)
		}
	}()

	// === GRACEFUL SHUTDOWN ===
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Infof("[RunServer-7] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("[RunServer-8] Server forced to shutdown: %v", err)
	}

	log.Infof("[RunServer-9] Server exited properly")
}
