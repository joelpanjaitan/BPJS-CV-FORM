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

type Exp struct {
	ID        int    `json:"id"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type Employment struct {
	ID        int    `json:"id"`
	Employer  string `json:"employer"`
	JobTitle  string `json:"job_title"`
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

type PhotoData struct {
	PhotoURL string `json:"photo_url"`
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

func GetExperience(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profile_id"]
	rows, err := database.DB.Query("SELECT id, company, position, start_date, end_date FROM experience WHERE profile_id = ?", profileID)
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

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(experiences)
}

func UpdateExperience(w http.ResponseWriter, r *http.Request) {
	var exp Exp
	
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE experience SET company = ?, position = ?, start_date = ?, end_date = ? WHERE id = ?",
		exp.Company, exp.Position, exp.StartDate, exp.EndDate, exp.ID)
	if err != nil {
		http.Error(w, "Failed to update experience", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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