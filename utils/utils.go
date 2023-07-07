package utils

import (
	"net/http"
	"godoit/database"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/sessions"
	"strconv"
)

// Will return true if the user is logged in, false if not
func CheckLogin(r *http.Request, store *sessions.CookieStore) bool {
	session, _ := store.Get(r, "user-session")
	if session.Values["userID"] != nil {
		return true
	}
	return false
}

func ExtractRegisterFormData(r *http.Request) database.RegisterFormData {
	r.ParseForm()
	var formValues database.RegisterFormData
	formValues.Email = r.PostForm.Get("email") 
	formValues.Password1 = r.PostForm.Get("password1")
	formValues.Password2 = r.PostForm.Get("password2")
	return formValues
}

func ValidatePassword(hashedPassword string, password string) bool {
	// Compare a password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Handle error
		return "", err
	}

	// Convert the hashed password to a string
	hashedPasswordString := string(hashedPassword)
	return hashedPasswordString, nil
}

func ExtractLoginFormData(r *http.Request) database.LoginFormData {
	r.ParseForm()
	var formValues database.LoginFormData
	formValues.Email = r.PostForm.Get("email")
	formValues.Password = r.PostForm.Get("password")
	return formValues
}


func ExtractTaskFormValues(r *http.Request) (database.TaskQuery, error) {
	r.ParseForm()
	var formValues database.TaskQuery

	//convert id to int
	idStr := r.PostForm.Get("id")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		formValues.Id = id
	}

	//convert duedate to time
	dueDateString := r.PostForm.Get("dueDate")
	dueDate, err := time.Parse("2006-01-02T15:04", dueDateString)
	if err != nil {
		// Example: Set a default value
		formValues.DueDate = time.Time{}
	}
	formValues.DueDate = dueDate

	formValues.Id = id
	formValues.Title = r.PostForm.Get("title")
	formValues.Description = r.PostForm.Get("description")
	return formValues, nil
}


func UpdateCurrentData(formResults database.TaskQuery, currentData database.TaskQuery) database.TaskQuery {
	if len(formResults.Title) > 0 {
		currentData.Title = formResults.Title
	}
	if len(formResults.Description) > 0 {
		currentData.Description = formResults.Description
	}
	if formResults.DueDate != (time.Time{}) && currentData.DueDate != formResults.DueDate {
		currentData.DueDate = formResults.DueDate
	}
	return currentData
}