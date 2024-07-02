// controllers/auth_controller.go

package controllers

import (
	"log"
	"net/http"
	"time"

	"backend-app/models"
	"backend-app/repository"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
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
func (ac *AuthController) SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error decoding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate input
	if user.Username == "" || user.PasswordHash == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Create user in database
	if err := ac.UserRepository.CreateUser(c.Request.Context(), &user); err != nil {
		log.Println("Error creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Clear sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var creds models.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("Error decoding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Retrieve user from database by email
	storedUser, err := ac.UserRepository.GetUserByEmail(creds.Email)
	if err != nil {
		log.Println("Error retrieving user:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(creds.Password)); err != nil {
		log.Println("Invalid password:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
