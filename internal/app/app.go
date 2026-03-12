package appapp

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"

	"time"
)

type App struct{
	log *slog.Logger
	GRPCSrv *grpcapp.App
	port int
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App{

	grpcApp := grpcapp.New(log, grpcPort)
	
	// authgrpc.Register(grpcServer)

	return &App{
		GRPCSrv: grpcApp,
	}
}
