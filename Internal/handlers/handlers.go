package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Joshdike/stocks-api/Internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type handle struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *handle {
	return &handle{db}
}

func (h handle) GetAllStock(w http.ResponseWriter, r *http.Request) {
	query, params, err := sq.Select("*").From("stocks").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		http.Error(w, "internal sql error", http.StatusInternalServerError)
		return
	}

	fmt.Println(query, params)

	rows, err := h.db.Query(r.Context(), query, params...)
	if err != nil {
		http.Error(w, "error retrieving data", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	result := make([]models.Stock, 0)
	for rows.Next() {
		var s models.Stock

		if err = rows.Scan(&s.StockID, &s.Name, &s.Price, &s.Company); err != nil {
			http.Error(w, "error retrieving data", http.StatusInternalServerError)
			return
		}

		result = append(result, s)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h handle) CreateStock(w http.ResponseWriter, r *http.Request) {
	var s models.Stock
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "error decoding payload", http.StatusBadRequest)
		return
	}

	query, params, err := sq.Insert("stocks").Columns("stockid", "name", "price", "company").
		Values(s.StockID, s.Name, s.Price, s.Company).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal sql error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("success: query: %s, params: %v\n", query, params)

	if _, err := h.db.Exec(r.Context(), query, params...); err != nil {
		http.Error(w, "error inserting stock data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "successful"}`))

}
func (h handle) GetStockById(w http.ResponseWriter, r *http.Request) {

	var s models.Stock
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatalf("conversion error %v", err)
	}
	query, params, err := sq.Select("stockid", "name", "price", "company").
		From("stocks").Where("stockid = ?", id).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		http.Error(w, "internal sql error", http.StatusInternalServerError)
		return
	}

	fmt.Println(query, params)

	err = h.db.QueryRow(r.Context(), query, params...).Scan(&s.StockID, &s.Name, &s.Price, &s.Company)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "stock not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error retrieving data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)

}
func (h handle) UpdateStock(w http.ResponseWriter, r *http.Request) {
	var s models.Stock
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatalf("conversion error %v", err)
	}

	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "error decoding payload", http.StatusBadRequest)
		return
	}

	updateS := sq.Update("stocks").
		Where("stockid = ?", id).PlaceholderFormat(sq.Dollar)

	if s.Name != "" {
		updateS = updateS.Set("name", s.Name)
	}
	if s.Price != 0 {
		updateS = updateS.Set("price", s.Price)
	}
	if s.Company != "" {
		updateS = updateS.Set("company", s.Company)
	}

	query, params, err := updateS.ToSql()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal sql error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("success: query: %s, params: %v\n", query, params)

	if _, err := h.db.Exec(r.Context(), query, params...); err != nil {
		http.Error(w, "error updating stock data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "successful"}`))
}
func (h handle) DeleteStock(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatalf("conversion error %v", err)
	}
	query, params, err := sq.Delete("stocks").Where("stockid = ?", id).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		http.Error(w, "internal sql error", http.StatusInternalServerError)
		return
	}

	fmt.Println(query, params)

	if _, err := h.db.Exec(r.Context(), query, params...); err != nil {
		http.Error(w, "error deleting stock data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "successful"}`))
}
