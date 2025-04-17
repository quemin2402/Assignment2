package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/quemin2402/order-service/internal/infrastructure/db"
	"github.com/quemin2402/order-service/internal/infrastructure/repository"
	"github.com/quemin2402/order-service/internal/infrastructure/transport"
	"github.com/quemin2402/order-service/internal/usecase"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dsn := env("ORD_DB_DSN", "postgres://postgres:0000@localhost:5432/order_service?sslmode=disable")

	pool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	uc := usecase.New(repository.New(pool))
	grpcSrv := transport.NewServer(uc)

	lis, _ := net.Listen("tcp", ":9002")
	go func() { log.Println("Order gRPC :9002"); _ = grpcSrv.Serve(lis) }()

	<-ctx.Done()
	grpcSrv.GracefulStop()
	log.Println("shutdown")
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
