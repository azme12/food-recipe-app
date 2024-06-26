package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"food-recipe-site-backend/models"
	"food-recipe-site-backend/repository"
)

type UserController struct {
	UserRepository *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{
		UserRepository: userRepo,
	}
}

// GetUser retrieves user data by ID
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Retrieve user data from the database
	user, err := uc.UserRepository.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Marshal user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

// UpdateUser updates user data
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if updatedUser.ID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Implement logic to update user data in the database
	err = uc.UserRepository.UpdateUser(&updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User data updated successfully")
}

// DeleteUser deletes a user by ID
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// TODO: Implement logic to delete user from the database
	err = uc.UserRepository.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}
