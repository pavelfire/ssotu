package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	// "strconv"
	authgrpc "ssotu/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port, //":" + strconv.Itoa(port),
	}
}

func (a *App) Run() error {
	const op = "grppcpp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil{
		return fmt.Errorf("%s: %w", op, err)
	}
	
}
