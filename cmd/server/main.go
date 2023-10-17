package main

import (
	"log"

	"github.com/leofartes-2/proglog/internal/server"
)

func main() {
	srv := server.CreateServer(":8080")
	log.Fatal(srv.ListenAndServe())

	// POST:
	// curl -X POST localhost:8080/log/write -d '{"record": {"value": "TGV0J3MgR28gIzEK"}}'
	// GET:
	// curl -X GET localhost:8080/log/read -d '{"offset": 0}'
}
