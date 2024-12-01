package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func GetWishesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(wishlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetWishHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, wish := range wishlist {
		if wish.ID != id { continue }
		if err := json.NewEncoder(w).Encode(wish); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Wish not found!", http.StatusNotFound)
}

func SaveWishHandler(w http.ResponseWriter, r *http.Request) {
	var payload Wish

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	newWish := Wish{
		ID: uuid.New(),
		Name: payload.Name,
		Description: payload.Description,
		CreatedAt: &now,
		UpdatedAt: nil,
	}

	if err = saveWish(newWish); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

func saveWish(newWish Wish) error {
	if newWish.Name == "" {
		return errors.New("name field is required")
	}

	if newWish.Description == "" {
		return errors.New("description field is required")
	}

	wishlist = append(wishlist, newWish)
	return nil;
}