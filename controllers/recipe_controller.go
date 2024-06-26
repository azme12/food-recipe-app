package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"food-recipe-site-backend/models"
	"food-recipe-site-backend/repository"
)

type RecipeController struct {
    RecipeRepository *repository.RecipeRepository
}

func NewRecipeController(recipeRepo *repository.RecipeRepository) *RecipeController {
    return &RecipeController{
        RecipeRepository: recipeRepo,
    }
}

// CreateRecipe creates a new recipe with image upload
func (rc *RecipeController) CreateRecipe(w http.ResponseWriter, r *http.Request) {
    // Parse request body for recipe details
    var newRecipe models.Recipe
    err := json.NewDecoder(r.Body).Decode(&newRecipe)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if newRecipe.Title == "" || newRecipe.CreatorID == 0 {
        http.Error(w, "Title and creator ID are required", http.StatusBadRequest)
        return
    }

    // Handle image upload
    images, err := uploadImages(r)
    if err != nil {
        http.Error(w, "Failed to upload images", http.StatusInternalServerError)
        return
    }
    newRecipe.Images = images

    // Create recipe in the database
    createdRecipe, err := rc.RecipeRepository.CreateRecipe(&newRecipe)
    if err != nil {
        log.Println("Error creating recipe:", err)
        http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
        return
    }

    // Return created recipe as JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdRecipe)
}

// UpdateRecipe updates an existing recipe
func (rc *RecipeController) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
    // Parse request body for updated recipe details
    var updatedRecipe models.Recipe
    err := json.NewDecoder(r.Body).Decode(&updatedRecipe)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if updatedRecipe.ID == 0 {
        http.Error(w, "Recipe ID is required", http.StatusBadRequest)
        return
    }

    // Handle image upload
    images, err := uploadImages(r)
    if err != nil {
        http.Error(w, "Failed to upload images", http.StatusInternalServerError)
        return
    }
    updatedRecipe.Images = images

    // Update recipe in the database
    err = rc.RecipeRepository.UpdateRecipe(&updatedRecipe)
    if err != nil {
        log.Println("Error updating recipe:", err)
        http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
        return
    }

    // Return success response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Recipe updated successfully")
}

// DeleteRecipe deletes a recipe
func (rc *RecipeController) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
    // Parse request parameters for recipe ID
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Recipe ID is required", http.StatusBadRequest)
        return
    }

    // Convert ID to int64
    recipeID := convertToInt64(id)
    if recipeID == 0 {
        http.Error(w, "Invalid Recipe ID", http.StatusBadRequest)
        return
    }

    // Delete recipe from the database
    err := rc.RecipeRepository.DeleteRecipe(recipeID)
    if err != nil {
        log.Println("Error deleting recipe:", err)
        http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
        return
    }

    // Return success response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Recipe deleted successfully")
}

// GetRecipe retrieves a single recipe by ID
func (rc *RecipeController) GetRecipe(w http.ResponseWriter, r *http.Request) {
    // Parse request parameters for recipe ID
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Recipe ID is required", http.StatusBadRequest)
        return
    }

    // Convert ID to int64
    recipeID := convertToInt64(id)
    if recipeID == 0 {
        http.Error(w, "Invalid Recipe ID", http.StatusBadRequest)
        return
    }

    // Retrieve recipe from the database
    recipe, err := rc.RecipeRepository.GetRecipeByID(recipeID)
    if err != nil {
        log.Println("Error retrieving recipe:", err)
        http.Error(w, "Failed to retrieve recipe", http.StatusInternalServerError)
        return
    }

    // Return recipe as JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(recipe)
}

// GetAllRecipes retrieves all recipes
func (rc *RecipeController) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
    // Retrieve all recipes from the database
    recipes, err := rc.RecipeRepository.GetAllRecipes()
    if err != nil {
        log.Println("Error retrieving recipes:", err)
        http.Error(w, "Failed to retrieve recipes", http.StatusInternalServerError)
        return
    }

    // Return recipes as JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(recipes)
}

// Helper function to convert string to int64
func convertToInt64(str string) int64 {
    num, _ := strconv.ParseInt(str, 10, 64)
    return num
}

// uploadImages handles image uploads
func uploadImages(r *http.Request) ([]string, error) {
    var images []string

    err := r.ParseMultipartForm(10 << 20) // 10 MB maximum file size
    if err != nil {
        return images, err
    }

    // Retrieve files from form data
    files := r.MultipartForm.File["images"]
    for _, file := range files {
        // Open uploaded file
        src, err := file.Open()
        if err != nil {
            return images, err
        }
        defer src.Close()

        // Create destination file
        dst, err := os.Create(filepath.Join("uploads", file.Filename))
        if err != nil {
            return images, err
        }
        defer dst.Close()

        // Copy file to destination
        if _, err := io.Copy(dst, src); err != nil {
            return images, err
        }

        // Store image URL or path
        images = append(images, filepath.Join("uploads", file.Filename))
    }

    return images, nil
}
