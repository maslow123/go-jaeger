package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/maslow123/go-jaeger/internal/pkg/storage"
	"github.com/maslow123/go-jaeger/internal/pkg/trace"
	"github.com/maslow123/go-jaeger/internal/users"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()

	// Bootstrap tracer.
	prv, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: "http://localhost:14268/api/traces",
		ServiceName:    "client",
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Disabled:       false,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer prv.Close(ctx)

	// Bootstrap database.
	dtb, err := sql.Open("mysql", "user:pass@tcp(:3306)/client")
	if err != nil {
		log.Fatalln(err)
	}
	defer dtb.Close()

	// Bootstrap API.
	usr := users.New(storage.NewUserStorage(dtb))

	// Bootstrap HTTP router.
	rtr := http.DefaultServeMux
	rtr.HandleFunc("/api/v1/users", trace.HTTPHandlerFunc(usr.Create, "users_create"))

	// Start HTTP server.
	log.Println("Served on :8888")
	if err := http.ListenAndServe(":8888", rtr); err != nil {
		log.Fatalln(err)
	}

}
