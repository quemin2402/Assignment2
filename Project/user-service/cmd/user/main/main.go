package main

import (
	"context"
	"github.com/quemin2402/user-service/internal/infrastructure/db"
	"github.com/quemin2402/user-service/internal/infrastructure/repository"
	"github.com/quemin2402/user-service/internal/infrastructure/transport"
	"github.com/quemin2402/user-service/internal/usecase"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dsn := env("USR_DB_DSN", "postgres://postgres:0000@localhost:5432/user_service?sslmode=disable")

	pool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()
	repo, err := repository.New(pool)
	if err != nil {
		log.Fatal(err)
	}

	uc := usecase.New(repo)
	grpcSrv := transport.NewServer(uc)
	lis, _ := net.Listen("tcp", ":9003")
	go func() { log.Println("User gRPC :9003"); _ = grpcSrv.Serve(lis) }()

	<-ctx.Done()
	grpcSrv.GracefulStop()
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
