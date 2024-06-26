package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"models"
	"repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserRepository *repository.UserRepository
	JwtSecret      []byte // Secret key for JWT
}

func NewAuthController(userRepo *repository.UserRepository, jwtSecret []byte) *AuthController {
	return &AuthController{
		UserRepository: userRepo,
		JwtSecret:      jwtSecret,
	}
}

// SignUp handles user registration
func (ac *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Username == "" || user.PasswordHash == "" || user.Email == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Create user in database
	err = ac.UserRepository.CreateUser(r.Context(), &user)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Clear sensitive information
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Login handles user authentication
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve user from database by email
	storedUser, err := ac.UserRepository.GetUserByEmail(creds.Email)
	if err != nil {
		log.Println("Error retrieving user:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(creds.Password))
	if err != nil {
		log.Println("Invalid password:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    storedUser.ID,
		"email": storedUser.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	tokenString, err := token.SignedString(ac.JwtSecret)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Middleware to authenticate requests using JWT
func (ac *AuthController) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ac.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store user information in context
		claims, _ := token.Claims.(jwt.MapClaims)
		userID := int(claims["id"].(float64))
		context.Set(r, "userID", userID)

		next.ServeHTTP(w, r)
	})
}
