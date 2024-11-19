package handlers

import (
	"bpjs-cv-form/database"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetEducation(w http.ResponseWriter, r *http.Request) {
	var record Record
	profileID := r.URL.Query().Get("profile_id")
	if profileID == "" {
		http.Error(w, "Profile ID is required", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query("SELECT id, school_name, degree, year_of_graduation FROM education WHERE profile_id = ?", profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var educationRecords []map[string]interface{}

	for rows.Next() {
		if err := rows.Scan(&record.ID, &record.SchoolName, &record.Degree, &record.YearOfGraduation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		educationRecords = append(educationRecords, map[string]interface{}{
			"id":                  record.ID,
			"school_name":         record.SchoolName,
			"degree":              record.Degree,
			"year_of_graduation":  record.YearOfGraduation,
		})
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(educationRecords)
}

func CreateEducation(w http.ResponseWriter, r *http.Request) {
	var record Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err := database.DB.Exec("INSERT INTO education (profile_id, school_name, degree, year_of_graduation) VALUES (?, ?, ?, ?)",
		record.ID, record.SchoolName, record.Degree, record.YearOfGraduation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteEducation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.DB.Exec("DELETE FROM education WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}