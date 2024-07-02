// routes/routes.go

package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gin-gonic/gin" // Import the gin package
	"backend-app/controllers"
	"backend-app/middleware" // Import the package that contains AuthMiddleware
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(router *mux.Router, authController *controllers.AuthController,
	userController *controllers.UserController, recipeController *controllers.RecipeController) {

	// Auth routes
router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
	ctx := &gin.Context{
		Request: r,
		Writer:  gin.ResponseWriter(w.(http.ResponseWriter)),
	}
	authController.SignUp(ctx)
}).Methods("POST")
router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	ctx := &gin.Context{
		Request: r,
		Writer:  w,
	}
	authController.Login(ctx)
}).Methods("POST")


	// User routes
	router.HandleFunc("/user", userController.GetUser).Methods("GET")
	router.HandleFunc("/user/update", middleware.AuthMiddleware(userController.UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/delete", middleware.AuthMiddleware(userController.DeleteUser)).Methods("DELETE")

	// Recipe routes
	router.HandleFunc("/recipe/create", middleware.AuthMiddleware(recipeController.CreateRecipe)).Methods("POST")
	router.HandleFunc("/recipe/update", middleware.AuthMiddleware(recipeController.UpdateRecipe)).Methods("PUT")
	router.HandleFunc("/recipe/delete", middleware.AuthMiddleware(recipeController.DeleteRecipe)).Methods("DELETE")

	// Serve static files (images)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}
