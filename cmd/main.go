package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Joshdike/stocks-api/Internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	connString := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil{
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := handlers.New(pool)

	r.Get("/api/stock", h.GetAllStock)
	r.Get("/api/stock/{id}", h.GetStockById)
	r.Post("/api/stock", h.CreateStock)
	r.Put("/api/stock/{id}", h.UpdateStock)
	r.Delete("/api/stock/{id}", h.DeleteStock)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}

}
