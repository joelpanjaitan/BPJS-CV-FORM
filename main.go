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
	router.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	router.HandleFunc("/profile", handlers.CreateProfile).Methods("POST")
	router.HandleFunc("/profile", handlers.UpdateProfile).Methods("PUT")

	router.HandleFunc("/photo", handlers.GetPhoto).Methods("GET")
	router.HandleFunc("/photo", handlers.UpdatePhoto).Methods("PUT")
	router.HandleFunc("/photo", handlers.DeletePhoto).Methods("DELETE")

	router.HandleFunc("/working-experience", handlers.GetExperience).Methods("GET")
	router.HandleFunc("/working-experience", handlers.UpdateExperience).Methods("PUT")

	router.HandleFunc("/employment", handlers.GetEmployment).Methods("GET")
	router.HandleFunc("/employment", handlers.CreateEmployment).Methods("POST")
	router.HandleFunc("/employment", handlers.DeleteEmployment).Methods("DELETE")

	router.HandleFunc("/education", handlers.GetEducation).Methods("GET")
	router.HandleFunc("/education", handlers.CreateEducation).Methods("POST")
	router.HandleFunc("/education", handlers.DeleteEducation).Methods("DELETE")

	router.HandleFunc("/skill", handlers.GetSkills).Methods("GET")
	router.HandleFunc("/skill", handlers.CreateSkill).Methods("POST")
	router.HandleFunc("/skill", handlers.DeleteSkill).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}