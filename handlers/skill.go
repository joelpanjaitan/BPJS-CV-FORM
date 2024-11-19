package handlers

import (
	"bpjs-cv-form/database"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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