package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/quemin2402/inventory-service/internal/infrastructure/db"
	"github.com/quemin2402/inventory-service/internal/infrastructure/repository"
	"github.com/quemin2402/inventory-service/internal/infrastructure/transport"
	"github.com/quemin2402/inventory-service/internal/usecase"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dsn := getenv("INV_DB_DSN", "postgres://postgres:0000@localhost:5432/inventory_service?sslmode=disable")

	pool, err := db.New(ctx, dsn)

	if _, err := pool.Exec(ctx, `
	  CREATE TABLE IF NOT EXISTS products(
		id text PRIMARY KEY,
		name text NOT NULL,
		category text NOT NULL,
		price numeric(12,2) NOT NULL,
		stock integer NOT NULL
	  );
	`); err != nil {
		log.Fatalf("init table: %v", err)
	}

	// pgx pool
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	uc := usecase.New(repository.New(pool))
	grpcSrv := transport.NewServer(uc)

	lis, _ := net.Listen("tcp", ":9001")
	go func() {
		log.Println("Inventory gRPC listening on :9001")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	grpcSrv.GracefulStop()
	log.Println("shutdown completed")
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
