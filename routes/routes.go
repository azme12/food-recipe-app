package routes

import (
    "net/http"

    "github.com/gorilla/mux"

   "github.com/azme12/backend-app/controllers"
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(router *mux.Router, authMiddleware *controllers.AuthMiddleware,
    authController *controllers.AuthController, userController *controllers.UserController,
    recipeController *controllers.RecipeController) {

    // Auth routes
    router.HandleFunc("/signup", authController.SignUp).Methods("POST")
    router.HandleFunc("/login", authController.Login).Methods("POST")

    // User routes
    router.HandleFunc("/user", userController.GetUser).Methods("GET")
    router.HandleFunc("/user/update", authMiddleware.Authenticate(userController.UpdateUser)).Methods("PUT")
    router.HandleFunc("/user/delete", authMiddleware.Authenticate(userController.DeleteUser)).Methods("DELETE")

    // Recipe routes
    router.HandleFunc("/recipe/create", authMiddleware.Authenticate(recipeController.CreateRecipe)).Methods("POST")
    router.HandleFunc("/recipe/update", authMiddleware.Authenticate(recipeController.UpdateRecipe)).Methods("PUT")
    router.HandleFunc("/recipe/delete", authMiddleware.Authenticate(recipeController.DeleteRecipe)).Methods("DELETE")

    // Serve static files (images)
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}
