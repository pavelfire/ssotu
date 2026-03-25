package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/pavelfire/protostu/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "clients.sso.grpc.New"

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Client{
		api: ssov1.NewAuthClient(conn), log: log,
	}, nil
}
