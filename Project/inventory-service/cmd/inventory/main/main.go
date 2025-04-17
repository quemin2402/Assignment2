package main

import (
	"context"
	"database/sql"
	"inventory-service/internal/infrastructure/db"
	"inventory-service/internal/infrastructure/repository"
	"inventory-service/internal/infrastructure/transport"
	"inventory-service/internal/usecase"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dsn := getenv("INV_DB_DSN", "postgres://postgres:postgres@localhost:5432/inventory_service?sslmode=disable")

	// migrate
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(sqlDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// pgx pool
	pool, err := db.New(ctx, dsn)
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
