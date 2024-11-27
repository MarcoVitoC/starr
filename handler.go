package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
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

func SaveWishHandler(w http.ResponseWriter, r *http.Request) {
	var payload Wish

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	newWish := Wish{
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