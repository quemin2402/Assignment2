package main

import (
	"github.com/quemin2402/api-gateway/internal"
	"log"
)

func main() {
	r := internal.NewRouter()
	log.Println("API‑Gateway listening on :8080")
	log.Fatal(r.Run(":8080"))
}
