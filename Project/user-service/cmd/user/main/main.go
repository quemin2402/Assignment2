package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user-service/internal/infrastructure/db"
	"user-service/internal/infrastructure/repository"
	"user-service/internal/infrastructure/transport"
	"user-service/internal/usecase"
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
