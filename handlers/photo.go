package handlers

import (
	"bpjs-cv-form/database"
	"encoding/json"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"
)

type PhotoData struct {
	PhotoURL string `json:"photo_url"`
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
	var photoData PhotoData
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