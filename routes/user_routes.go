package routes

import (
	"golang-jwt-app/controllers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/auth/signup", controllers.SignupHandler).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", controllers.LoginHandler).Methods("POST")
	//r.HandleFunc("/api/v1/auth/welcome", controllers.WelcomeHandler).Methods("GET")
}
