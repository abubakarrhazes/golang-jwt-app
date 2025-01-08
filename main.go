package main

import (
	"log"
	"net/http"

	"golang-jwt-app/controllers"
	"golang-jwt-app/models"
	"golang-jwt-app/routes"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Database Connection
	dsn := "root:a2144@tcp(127.0.0.1:3306)/new_db"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate Models
	models.MigrateUsers(db)

	// Initialize Controllers
	controllers.InitUserController(db)

	// Setup Routes
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
