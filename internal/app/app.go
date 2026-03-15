package appapp

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"

	"time"
)

type App struct {
	log     *slog.Logger
	GRPCSrv *grpcapp.App
	port    int
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("error", err.Error()))
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
