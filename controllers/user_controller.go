package controllers

import (
	"encoding/json"
	"golang-jwt-app/models"
	"golang-jwt-app/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Assuming db is already initialized
var db *gorm.DB

// InitUserController initializes the db instance in controller
func InitUserController(database *gorm.DB) {
	db = database
}

// Signup Handler
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if user.Username == "" || user.Password == "" || user.Email == "" {
		utils.SendError(w, http.StatusBadRequest, "Missing required fields: username, password, and email are required")
		return
	}

	// Check if user already exists
	var existingUser models.User
	if db.Where("email = ?", user.Email).First(&existingUser).RowsAffected > 0 {
		utils.SendError(w, http.StatusBadRequest, "User with this email already exists")
		return
	}
	if db.Where("username = ?", user.Username).First(&existingUser).RowsAffected > 0 {
		utils.SendError(w, http.StatusBadRequest, "Username already taken")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save the new user in the database
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login Handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Username == "" || input.Password == "" {
		utils.SendError(w, http.StatusBadRequest, "Username and password are required")
		return
	}
	var user models.User
	if db.Where("username = ?", input.Username).First(&user).RowsAffected == 0 {
		utils.SendError(w, http.StatusUnauthorized, "User not found")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	// Send response with the JWT token and user data
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
