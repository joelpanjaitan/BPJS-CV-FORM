package main

import (
	"bpjs-cv-form/handlers"
	"log"
	"net/http"

	"bpjs-cv-form/database"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDatabase()

	router := mux.NewRouter()
	router.HandleFunc("/api/profile", handlers.GetProfile).Methods("GET")
	router.HandleFunc("/api/profile/{id}", handlers.GetProfileDetail).Methods("GET")
	router.HandleFunc("/api/profile", handlers.CreateProfile).Methods("POST")
	router.HandleFunc("/api/profile/{id}", handlers.UpdateProfile).Methods("PUT")

	router.HandleFunc("/api/photo", handlers.GetPhoto).Methods("GET")
	router.HandleFunc("/api/photo", handlers.UpdatePhoto).Methods("PUT")
	router.HandleFunc("/api/photo", handlers.DeletePhoto).Methods("DELETE")

	router.HandleFunc("/api/working-experience", handlers.GetExperience).Methods("GET")
	router.HandleFunc("/api/working-experience/{id}", handlers.GetExperience).Methods("GET")
	router.HandleFunc("/api/working-experience", handlers.UpdateExperience).Methods("PUT")

	router.HandleFunc("/api/employment", handlers.GetEmployment).Methods("GET")
	router.HandleFunc("/api/employment", handlers.CreateEmployment).Methods("POST")
	router.HandleFunc("/api/employment/{id}", handlers.DeleteEmployment).Methods("DELETE")

	router.HandleFunc("/api/education", handlers.GetEducation).Methods("GET")
	router.HandleFunc("/api/education", handlers.CreateEducation).Methods("POST")
	router.HandleFunc("/api/education/{id}", handlers.DeleteEducation).Methods("DELETE")

	router.HandleFunc("/api/skill", handlers.GetSkills).Methods("GET")
	router.HandleFunc("/api/skill", handlers.CreateSkill).Methods("POST")
	router.HandleFunc("/api/skill/{id}", handlers.DeleteSkill).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}