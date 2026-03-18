package suite

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	authgrpc "sso/internal/grpc/auth"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"

	ssov1 "github.com/pavelfire/protostu/gen/go/sso"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

type Suite struct {
	*testing.T                  // for access to testing.T methods
	Cfg        *Config          // for access to config
	AuthClient ssov1.AuthClient // client for interacting with gRPC server
}

type Config struct {
	TokenTTL time.Duration
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := &Config{
		TokenTTL: time.Hour,
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 6*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	storagePath := filepath.Join(t.TempDir(), "sso_test.db")
	requireMigrationsApplied(t, storagePath)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	storage, err := sqlite.New(storagePath)
	if err != nil {
		t.Fatalf("failed to create sqlite storage: %v", err)
	}
	t.Cleanup(func() {
		_ = storage.Close()
	})

	authService := auth.New(log, storage, storage, storage, cfg.TokenTTL)

	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:0", grpcHost))
	if err != nil {
		t.Fatalf("failed to listen on random port: %v", err)
	}
	t.Cleanup(func() { _ = l.Close() })

	go func() {
		_ = gRPCServer.Serve(l)
	}()
	t.Cleanup(func() { gRPCServer.GracefulStop() })

	cc, err := grpc.DialContext(
		context.Background(),
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %v", err)
	}
	t.Cleanup(func() { _ = cc.Close() })

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}

}

func requireMigrationsApplied(t *testing.T, storagePath string) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate suite.go path")
	}
	// thisFile: .../ssotu/tests/suite/suite.go
	testsDir := filepath.Dir(filepath.Dir(thisFile))       // .../ssotu/tests
	rootDir := filepath.Dir(testsDir)                      // .../ssotu
	mainMigrations := filepath.Join(rootDir, "migrations")  // .../ssotu/migrations
	testMigrations := filepath.Join(testsDir, "migrations") // .../ssotu/tests/migrations

	apply := func(migrationsDir string, table string) {
		absMigrations, err := filepath.Abs(migrationsDir)
		if err != nil {
			t.Fatalf("failed to abs migrations path: %v", err)
		}
		dbURL := fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, table)
		m, err := migrate.New("file://"+filepath.ToSlash(absMigrations), dbURL)
		if err != nil {
			t.Fatalf("failed to init migrator: %v", err)
		}
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("failed to apply migrations (%s): %v", absMigrations, err)
		}
	}

	// Keep separate tables so the same DB can apply both sets.
	apply(mainMigrations, "schema_migrations")
	apply(testMigrations, "test_schema_migrations")
}
