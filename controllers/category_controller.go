package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"backend-app/models"
	"backend-app/repository"
)

type CategoryController struct {
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryController(categoryRepo *repository.CategoryRepository) *CategoryController {
	return &CategoryController{
		CategoryRepository: categoryRepo,
	}
}

// CreateCategory creates a new recipe category
func (cc *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var newCategory models.Category
    err := json.NewDecoder(r.Body).Decode(&newCategory)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if newCategory.Name == "" {
        http.Error(w, "Category name is required", http.StatusBadRequest)
        return
    }


	// Create the category in the database
	err = cc.CategoryRepository.CreateCategory(&newCategory)
    if err != nil {
        log.Println("Error creating category:", err)
        http.Error(w, "Failed to create category", http.StatusInternalServerError)
        return
    }

	// Return success response with the created category
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newCategory) // Adjust based on how you handle the created category
}

// UpdateCategory updates an existing recipe category
func (cc *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var updatedCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if updatedCategory.ID == 0 {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	// Update the category in the database
	err = cc.CategoryRepository.UpdateCategory(&updatedCategory)
	if err != nil {
		log.Println("Error updating category:", err)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category updated successfully")
}

// DeleteCategory deletes a recipe category by ID
func (cc *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("id")
	if categoryIDStr == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Delete the category from the database
	err = cc.CategoryRepository.DeleteCategory(categoryID)
	if err != nil {
		log.Println("Error deleting category:", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category deleted successfully")
}

// GetAllCategories retrieves all recipe categories
func (cc *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Retrieve all categories from the database
	categories, err := cc.CategoryRepository.GetAllCategories()
	if err != nil {
		log.Println("Error retrieving categories:", err)
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}

	// Return categories as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
