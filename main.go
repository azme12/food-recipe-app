package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq" // PostgreSQL driver
    "backend-app/controllers"
    "backend-app/repository"
    "backend-app/routes"
)

// Config represents the application configuration structure
type Config struct {
    DatabaseURL string
    Port        string
}

func main() {
    // Load configuration from environment variables
    cfg := Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Port:        os.Getenv("PORT"),
    }

    // Initialize database connection
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    recipeRepo := repository.NewRecipeRepository(db)

    // Initialize controllers
    authController := controllers.NewAuthController(userRepo)
    userController := controllers.NewUserController(userRepo)
    recipeController := controllers.NewRecipeController(recipeRepo)

    // Initialize middleware
    authMiddleware := controllers.NewAuthMiddleware()

    // Initialize router
    router := mux.NewRouter()

    // Register routes
    routes.RegisterRoutes(router, authMiddleware, authController, userController, recipeController)

    // Serve static files (images)
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    // Start server
    port := cfg.Port
    fmt.Printf("Server running on port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}
