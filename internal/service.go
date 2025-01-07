package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WishService struct {
	conn *pgxpool.Pool
}

func (s *WishService) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT * FROM wishes`
	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wishes []Wish
	for rows.Next() {
		var wish Wish

		if err := rows.Scan(&wish.ID, &wish.Name, &wish.Description, &wish.CreatedAt, &wish.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wishes = append(wishes, wish)
	}

	if err := json.NewEncoder(w).Encode(wishes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *WishService) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM wishes WHERE id = $1`
	row := s.conn.QueryRow(context.Background(), query, id)

	var wish Wish
	if err := row.Scan(&wish.ID, &wish.Name, &wish.Description, &wish.CreatedAt, &wish.UpdatedAt); err != nil {
		http.Error(w, "Wish not found!", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(wish); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *WishService) Save(w http.ResponseWriter, r *http.Request) {
	var payload Wish

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	newWish := Wish{
		ID:          uuid.New(),
		Name:        payload.Name,
		Description: payload.Description,
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	if err = saveWish(s.conn, newWish); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func saveWish(db *pgxpool.Pool, newWish Wish) error {
	if newWish.Name == "" {
		return errors.New("ERROR: name field is required")
	}

	if newWish.Description == "" {
		return errors.New("ERROR: description field is required")
	}

	query := `INSERT INTO wishes (id, name, description, created_at, updated_at) VALUES($1, $2, $3, $4, $5)`
	_, err := db.Exec(
		context.Background(),
		query,
		newWish.ID, newWish.Name, newWish.Description, newWish.CreatedAt, newWish.UpdatedAt,
	)

	if err != nil {
		log.Fatal("ERROR: ", err)
		return errors.New("ERROR: failed to insert new wish")
	}

	return nil
}

func (s *WishService) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM wishes WHERE id = $1`
	row := s.conn.QueryRow(context.Background(), query, id)

	var selectedWish Wish
	if err := row.Scan(
		&selectedWish.ID,
		&selectedWish.Name,
		&selectedWish.Description,
		&selectedWish.CreatedAt,
		&selectedWish.UpdatedAt,
	); err != nil {
		http.Error(w, "Wish not found!", http.StatusNotFound)
		return
	}

	var payload Wish
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	toBeUpdatedWish := Wish{
		ID:          selectedWish.ID,
		Name:        payload.Name,
		Description: payload.Description,
		CreatedAt:   selectedWish.CreatedAt,
		UpdatedAt:   &now,
	}

	if err := updateWish(s.conn, toBeUpdatedWish, id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateWish(db *pgxpool.Pool, toBeUpdatedWish Wish, id uuid.UUID) error {
	if toBeUpdatedWish.Name == "" {
		return errors.New("name field is required")
	}

	if toBeUpdatedWish.Description == "" {
		return errors.New("description field is required")
	}

	query := `
		UPDATE wishes 
		SET name = $1, description = $2 
		WHERE id = $3`
	_, err := db.Exec(context.Background(), query, toBeUpdatedWish.Name, toBeUpdatedWish.Description, id)
	if err != nil {
		log.Fatal("ERROR: ", err)
		return errors.New("ERROR: failed to updated selected wish")
	}

	return nil
}

func (s *WishService) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `DELETE FROM wishes WHERE id = $1`
	result, err := s.conn.Exec(context.Background(), query, id)
	if err != nil {
		http.Error(w, "Failed to delete wish", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Wish not found!", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
