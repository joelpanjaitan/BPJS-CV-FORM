package handlers

import (
	"bpjs-cv-form/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Profile struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	PhotoURL  string `json:"photo_url"`
	Summary string `json:"summary"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Record struct {
	ID               int    `json:"id"`
	SchoolName       string `json:"school_name"`
	Degree           string `json:"degree"`
	YearOfGraduation int    `json:"year_of_graduation"`
}

type Skill struct {
	ID       int    `json:"id"`
	SkillName string `json:"skill_name"`
}


func GetProfile(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email, phone, photo_url, summary, created_at, updated_at FROM profile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		if err := rows.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Phone, &profile.PhotoURL, &profile.Summary, &profile.CreatedAt, &profile.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profiles = append(profiles, profile)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if profile.Name == ""||profile.Email == "" {
		http.Error(w, "You must fill the name and email which are required", http.StatusBadRequest)
		return 
	}

	result, err := database.DB.Exec("INSERT INTO profile (name, email, phone, photo_url, summary) VALUES (?, ?, ?, ?, ?)",
		profile.Name, profile.Email, profile.Phone, profile.PhotoURL,profile.Summary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	profile.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "Input ID is required", http.StatusBadRequest)
		return
	}

	profileID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Input invalid for profile ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}


	if profile.Name == ""||profile.Email == "" {
		http.Error(w, "You must fill the name and email which are required", http.StatusBadRequest)
		return 
	}

	result, err := database.DB.Exec("UPDATE profile SET name = ?, email = ?, phone = ?, photo_url = ?, summary = ?, updated_at = NOW() WHERE id = ?",
		profile.Name, profile.Email, profile.Phone, profile.PhotoURL, profile.Summary, profileID)
	if err != nil {
		http.Error(w, "Failed to update profile data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Profile not found or no changes made", http.StatusNotFound)
		return
	}

	profile.ID = profileID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func GetProfileDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Input invalid for profile ID", http.StatusBadRequest)
		return
	}

	row := database.DB.QueryRow("SELECT id, name, email, phone, photo_url, summary, created_at, updated_at FROM profile WHERE id = ?", id)
	var profile Profile
	if err := row.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Phone, &profile.PhotoURL, &profile.Summary, &profile.CreatedAt, &profile.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Detailed profile not found", http.StatusNotFound)
		} else {
			http.Error(w, "Can't fetching profile details", http.StatusInternalServerError)
	
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}


func GetSkills(w http.ResponseWriter, r *http.Request) {
	var skill Skill
	profileID := r.URL.Query().Get("profile_id")
	if profileID == "" {
		http.Error(w, "Profile ID is required", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query("SELECT id, skill_name FROM skill WHERE profile_id = ?", profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var skills []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(&skill.ID, &skill.SkillName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		skills = append(skills, map[string]interface{}{
			"id":         skill.ID,
			"skill_name": skill.SkillName,
		})
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(skills)
}

func CreateSkill(w http.ResponseWriter, r *http.Request) {
	var skill Skill
	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("INSERT INTO skill (profile_id, skill_name) VALUES (?, ?)", skill.ID, skill.SkillName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteSkill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.DB.Exec("DELETE FROM skill WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}