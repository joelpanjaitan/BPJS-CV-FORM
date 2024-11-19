package handlers

import (
	"bpjs-cv-form/database"
	"strconv"

	// "database/sql"
	"encoding/json"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"
)

type Exp struct {
	ID        int    `json:"id"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func GetExperience(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profile_id"]
	rows, err := database.DB.Query("SELECT id, company_name, role, start_date, end_date FROM working_experience WHERE profile_id = ?", profileID)
	if err != nil {
		http.Error(w, "Error retrieving experiences", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var experiences []map[string]interface{}

	for rows.Next() {
		var exp Exp
		
		if err := rows.Scan(&exp.ID, &exp.Company, &exp.Position, &exp.StartDate, &exp.EndDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		experiences = append(experiences, map[string]interface{}{
			"id":         exp.ID,
			"company":    exp.Company,
			"position":   exp.Position,
			"start_date": exp.StartDate,
			"end_date":   exp.EndDate,
		})
	}

	if len(experiences) == 0 {
		http.Error(w, "Experiences not found for this profile", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(experiences)
}

func UpdateExperience(w http.ResponseWriter, r *http.Request) {
	var exp Exp
		vars := mux.Vars(r)
	expID, err := strconv.Atoi(vars["experience_id"])
	if err != nil {
		http.Error(w, "Invalid experience ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if exp.Company == "" || exp.Position == "" {
		http.Error(w, "Company and Position are required fields", http.StatusBadRequest)
		return
	}
	
	_, err = database.DB.Exec("UPDATE experience SET company = ?, position = ?, start_date = ?, end_date = ? WHERE id = ?",
		exp.Company, exp.Position, exp.StartDate, exp.EndDate, expID)
	if err != nil {
		http.Error(w, "Failed to update experience", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Experience updated successfully",
	})
}