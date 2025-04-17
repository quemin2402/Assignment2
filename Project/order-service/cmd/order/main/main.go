package main

import (
	"Assignment2/Project/order-service/internal/infrastructure/db"
	"Assignment2/Project/order-service/internal/infrastructure/repository"
	"Assignment2/Project/order-service/internal/infrastructure/transport"
	"Assignment2/Project/order-service/internal/usecase"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	dsn := env("ORD_DB_DSN", "postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable")

	sqlDB, _ := sql.Open("pgx", dsn)
	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	pool, _ := db.New(ctx, dsn)
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
