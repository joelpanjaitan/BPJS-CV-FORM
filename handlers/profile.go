package handlers

import (
	"bpjs-cv-form/database"
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
	Summary string `json:"summary"`
}

type Photo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Summary string `json:"summary"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email, phone, summary FROM profile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		if err := rows.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Phone, &profile.Summary); err != nil {
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

	_, err := database.DB.Exec("INSERT INTO profile (name, email, phone, summary) VALUES (?, ?, ?, ?)",
		profile.Name, profile.Email, profile.Phone, profile.Summary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE profile SET name = ?, email = ?, phone = ?, summary = ? WHERE id = ?",
		profile.Name, profile.Email, profile.Phone, profile.Summary, profile.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetProfileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	row := database.DB.QueryRow("SELECT id, name, email, phone, summary FROM profile WHERE id = ?", id)
	var profile Profile
	if err := row.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Phone, &profile.Summary); err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profile_id"]
	if profileID == "" {
		http.Error(w, "Profile ID is required", http.StatusBadRequest)
		return
	}

	var photoURL string
	err := database.DB.QueryRow("SELECT photo_url FROM profile WHERE id = ?", profileID).Scan(&photoURL)
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"photo_url": photoURL})
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profile_id"]

	var photoData struct {
		PhotoURL string `json:"photo_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&photoData); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE profile SET photo_url = ? WHERE id = ?", photoData.PhotoURL, profileID)
	if err != nil {
		http.Error(w, "Failed to update photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profile_id"]
	_, err := database.DB.Exec("UPDATE profile SET photo_url = NULL WHERE id = ?", profileID)
	if err != nil {
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetExperience(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profile_id"]
	rows, err := database.DB.Query("SELECT id, company, position, start_date, end_date FROM experience WHERE profile_id = ?", profileID)
	if err != nil {
		http.Error(w, "Error retrieving experiences", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
}

func UpdateExperience(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE profile SET name = ?, email = ?, phone = ?, summary = ? WHERE id = ?",
		profile.Name, profile.Email, profile.Phone, profile.Summary, profile.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}