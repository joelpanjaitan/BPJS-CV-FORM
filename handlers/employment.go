package handlers

import (
	"bpjs-cv-form/database"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Employment struct {
	ID        int    `json:"id"`
	Employer  string `json:"employer"`
	JobTitle  string `json:"job_title"`
}

func GetEmployment(w http.ResponseWriter, r *http.Request) {
	var employment Employment
	profileID := r.URL.Query().Get("profile_id")
	if profileID == "" {
		http.Error(w, "Profile ID is required", http.StatusBadRequest)
		return
	}
	rows, err := database.DB.Query("SELECT id, employer, job_title FROM employment WHERE profile_id = ?", profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employments []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(&employment.ID, &employment.Employer, &employment.JobTitle); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employments = append(employments, map[string]interface{}{
			"id":        employment.ID,
			"employer":  employment.Employer,
			"job_title": employment.JobTitle,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employments)
}

func CreateEmployment(w http.ResponseWriter, r *http.Request) {
	var employment Employment
	if err := json.NewDecoder(r.Body).Decode(&employment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err := database.DB.Exec("INSERT INTO employment (id, employer, job_title) VALUES (?, ?, ?)",
		employment.ID, employment.Employer, employment.JobTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteEmployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.DB.Exec("DELETE FROM employment WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}