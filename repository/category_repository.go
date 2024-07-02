package repository

import (
	"backend-app/models"
	"database/sql"
	"log"
)

type CategoryRepository struct {
	DB *sql.DB
}

// NewCategoryRepository initializes a new CategoryRepository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

// CreateCategory creates a new category in the database
func (cr *CategoryRepository) CreateCategory(category *models.Category) error {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id
	`
	err := cr.DB.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		log.Println("Error creating category:", err)
		return err
	}
	return nil
}

// UpdateCategory updates an existing category in the database
func (cr *CategoryRepository) UpdateCategory(category *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
	`
	_, err := cr.DB.Exec(query, category.Name, category.ID)
	if err != nil {
		log.Println("Error updating category:", err)
		return err
	}
	return nil
}

// DeleteCategory deletes a category from the database by ID
func (cr *CategoryRepository) DeleteCategory(categoryID int64) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`
	_, err := cr.DB.Exec(query, categoryID)
	if err != nil {
		log.Println("Error deleting category:", err)
		return err
	}
	return nil
}

// GetCategoryByID retrieves a category from the database by ID
func (cr *CategoryRepository) GetCategoryByID(categoryID int64) (*models.Category, error) {
	var category models.Category
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`
	err := cr.DB.QueryRow(query, categoryID).Scan(&category.ID, &category.Name)
	if err != nil {
		log.Println("Error retrieving category:", err)
		return nil, err
	}
	return &category, nil
}

// GetAllCategories retrieves all categories from the database
func (cr *CategoryRepository) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	query := `
		SELECT id, name
		FROM categories
	`
	rows, err := cr.DB.Query(query)
	if err != nil {
		log.Println("Error retrieving categories:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Println("Error scanning category row:", err)
			continue
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over category rows:", err)
		return nil, err
	}

	return categories, nil
}
