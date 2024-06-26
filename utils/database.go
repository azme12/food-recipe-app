package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "azme07@"
	dbname   = "food_recipes"
)

func main() {
	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables
	createUsersTable(db)
	createRecipesTable(db)
	createStepsTable(db)
	createIngredientsTable(db)
	createCategoriesTable(db)

	fmt.Println("Tables created successfully")
}

func createUsersTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func createRecipesTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			prep_time INT,
			category_id INT,
			creator_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT current_timestamp,
			updated_at TIMESTAMP DEFAULT current_timestamp,
			FOREIGN KEY (category_id) REFERENCES categories(id),
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func createStepsTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS steps (
			id SERIAL PRIMARY KEY,
			recipe_id INT NOT NULL,
			step_number INT NOT NULL,
			description TEXT,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func createIngredientsTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ingredients (
			id SERIAL PRIMARY KEY,
			recipe_id INT NOT NULL,
			name VARCHAR(255),
			quantity VARCHAR(50),
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func createCategoriesTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

// CREATE TABLE users (
// 	id SERIAL PRIMARY KEY,
// 	username TEXT UNIQUE NOT NULL,
// 	email TEXT UNIQUE NOT NULL,
// 	password TEXT NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//   );

//   CREATE TABLE categories (
// 	id SERIAL PRIMARY KEY,
// 	name TEXT UNIQUE NOT NULL
//   );

//   CREATE TABLE recipes (
// 	id SERIAL PRIMARY KEY,
// 	user_id INTEGER REFERENCES users(id),
// 	category_id INTEGER REFERENCES categories(id),
// 	title TEXT NOT NULL,
// 	description TEXT,
// 	preparation_time INTEGER,
// 	featured_image TEXT,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//   );

//   CREATE TABLE steps (
// 	id SERIAL PRIMARY KEY,
// 	recipe_id INTEGER REFERENCES recipes(id),
// 	step_number INTEGER NOT NULL,
// 	description TEXT NOT NULL
//   );

//   CREATE TABLE ingredients (
// 	id SERIAL PRIMARY KEY,
// 	recipe_id INTEGER REFERENCES recipes(id),
// 	name TEXT NOT NULL,
// 	quantity TEXT NOT NULL
//   );

//   CREATE TABLE likes (
// 	id SERIAL PRIMARY KEY,
// 	user_id INTEGER REFERENCES users(id),
// 	recipe_id INTEGER REFERENCES recipes(id)
//   );

//   CREATE TABLE bookmarks (
// 	id SERIAL PRIMARY KEY,
// 	user_id INTEGER REFERENCES users(id),
// 	recipe_id INTEGER REFERENCES recipes(id)
//   );

//   CREATE TABLE comments (
// 	id SERIAL PRIMARY KEY,
// 	user_id INTEGER REFERENCES users(id),
// 	recipe_id INTEGER REFERENCES recipes(id),
// 	content TEXT NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//   );

//   CREATE TABLE ratings (
// 	id SERIAL PRIMARY KEY,
// 	user_id INTEGER REFERENCES users(id),
// 	recipe_id INTEGER REFERENCES recipes(id),
// 	rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5)
//   );
