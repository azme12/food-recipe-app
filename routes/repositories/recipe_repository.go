package repository

import (
	"database/sql"
	"log"

	"backend-app/models"
)

type RecipeRepository struct {
	DB *sql.DB
}

// NewRecipeRepository initializes a new RecipeRepository
func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{
		DB: db,
	}
}

// CreateRecipe creates a new recipe in the database
func (rr *RecipeRepository) CreateRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	query := `
			INSERT INTO recipes (title, description, ingredients, steps, prep_time, category_id, creator_id, images)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
	err := rr.DB.QueryRow(
		query,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Steps,
		recipe.PrepTime,
		recipe.CategoryID,
		recipe.CreatorID,
		recipe.Images,
	).Scan(&recipe.ID)
	if err != nil {
		log.Println("Error creating recipe:", err)
		return nil, err
	}
	return recipe, nil
}

// UpdateRecipe updates an existing recipe in the database
func (rr *RecipeRepository) UpdateRecipe(recipe *models.Recipe) error {
	query := `
		UPDATE recipes
		SET title = $1, description = $2, ingredients = $3, steps = $4, prep_time = $5, category_id = $6, images = $7
		WHERE id = $8
	`
	_, err := rr.DB.Exec(
		query,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Steps,
		recipe.PrepTime,
		recipe.CategoryID,
		recipe.Images,
		recipe.ID,
	)
	if err != nil {
		log.Println("Error updating recipe:", err)
		return err
	}
	return nil
}

// DeleteRecipe deletes a recipe from the database by ID
func (rr *RecipeRepository) DeleteRecipe(recipeID int64) error {
	query := `
		DELETE FROM recipes
		WHERE id = $1
	`
	_, err := rr.DB.Exec(query, recipeID)
	if err != nil {
		log.Println("Error deleting recipe:", err)
		return err
	}
	return nil
}

// GetRecipeByID retrieves a recipe from the database by ID
func (rr *RecipeRepository) GetRecipeByID(recipeID int64) (*models.Recipe, error) {
	var recipe models.Recipe
	query := `
		SELECT id, title, description, ingredients, steps, prep_time, category_id, creator_id, images
		FROM recipes
		WHERE id = $1
	`
	err := rr.DB.QueryRow(query, recipeID).Scan(
		&recipe.ID,
		&recipe.Title,
		&recipe.Description,
		&recipe.Ingredients,
		&recipe.Steps,
		&recipe.PrepTime,
		&recipe.CategoryID,
		&recipe.CreatorID,
		&recipe.Images,
	)
	if err != nil {
		log.Println("Error retrieving recipe:", err)
		return nil, err
	}
	return &recipe, nil
}

// GetAllRecipes retrieves all recipes from the database
func (rr *RecipeRepository) GetAllRecipes() ([]*models.Recipe, error) {
	var recipes []*models.Recipe
	query := `
		SELECT id, title, description, ingredients, steps, prep_time, category_id, creator_id, images
		FROM recipes
	`
	rows, err := rr.DB.Query(query)
	if err != nil {
		log.Println("Error retrieving recipes:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			&recipe.Ingredients,
			&recipe.Steps,
			&recipe.PrepTime,
			&recipe.CategoryID,
			&recipe.CreatorID,
			&recipe.Images,
		)
		if err != nil {
			log.Println("Error scanning recipe row:", err)
			continue
		}
		recipes = append(recipes, &recipe)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over recipe rows:", err)
		return nil, err
	}

	return recipes, nil
}
