package main

import (
    "log"
	"fmt"
	"os"
	"encoding/hex"

	"net/http"
	"database/sql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"godoit/database"
	"godoit/frontend"
	"godoit/utils"
)

var store *sessions.CookieStore


func crudDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !utils.CheckLogin(r, store) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	session, _ := store.Get(r, "user-session")
	
	// Get POST form values
	formResults, err := utils.ExtractTaskFormValues(r)
	fmt.Println(formResults.Id)
	// Change error type
	if err != nil {
		http.Error(w, "Oops! Something went wrong.", http.StatusNotFound)
		return
	}
	database.DeleteTask(db, formResults, session.Values["userID"].(int))
}

func crudCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !utils.CheckLogin(r, store) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	session, _ := store.Get(r, "user-session")

	// Get POST form values
	formResults, err := utils.ExtractTaskFormValues(r)
	// Change error type
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Oops! Something went wrong.", http.StatusNotFound)
		return
	}
	database.InsertIntoTasks(db, formResults, session.Values["userID"].(int))
	http.Redirect(w, r, "/tasks", http.StatusFound)
	return
}


func crudUpdate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !utils.CheckLogin(r, store) {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	session, _ := store.Get(r, "user-session")

	// Get POST form values
	formResults, err := utils.ExtractTaskFormValues(r)
	if err != nil {
		fmt.Println("Error happened in first check", err)
		http.Error(w, "Oops! Something went wrong.", http.StatusNotFound)
		return
	}
	currentData, err := database.SelectSingleQuery(db, formResults.Id, session.Values["userID"].(int))
	if err != nil {
		fmt.Println("Error happened in second check", err)
		http.Error(w, `Oops! Something went wrong.\
		Did you try editing a task that does not belong to you?`, http.StatusNotFound)
		return
	}
	currentData = utils.UpdateCurrentData(formResults, currentData)
	database.UpdateTask(db, currentData, session.Values["userID"].(int))
	http.Redirect(w, r, "/tasks", http.StatusFound)
	return
}

func crudRead(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !utils.CheckLogin(r, store) {
		fmt.Println("Useri s not logged in!")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	session, _ := store.Get(r, "user-session")
	// Loads tasks from db to memory
	tasks := database.SelectManyQuery(db, session.Values["userID"].(int))

	// Defines data for HTML template
	data := map[string]interface{}{
		"Tasks": tasks,
	}

	// Loads HTML template file
	tmpl, err := frontend.LoadHTML("tasks.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}

func register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Must check that the user is not logged in first!!

	// Should get values from form.
	// Should check that both passwords are equal
	// Should hash passwords
	// Should send insert query into database


	// Get POST form values
	formData := utils.ExtractRegisterFormData(r)
	if formData.Email != "" || formData.Password1 != "" || formData.Password2 != "" {
		if formData.Password1 != formData.Password2 {
			fmt.Fprintf(w, "Your passwords don't match!")
			return
		}
		var err error
		formData.Password1, err = utils.HashPassword(formData.Password1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		formData.Password2 = formData.Password1
		// Must check for errors
		database.InsertIntoUsers(db, formData)
	}

	// Loads the HTML template file
	tmpl, err := frontend.LoadHTML("register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func about(w http.ResponseWriter, r *http.Request) {
	tmpl, err := frontend.LoadHTML("about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Must check that the user is not logged in first!!

	// Get POST form values
	loginData := utils.ExtractLoginFormData(r)
	if loginData.Email != "" || loginData.Password != "" {
		fmt.Println("Data:", loginData.Email, loginData.Password)
		queryResults, err := database.SelectLoginQuery(db, loginData, 1)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Something went wrong in the initial login query")
			return
		}
		hashedPassword := queryResults.Password
		if !utils.ValidatePassword(hashedPassword, loginData.Password) {
			http.Error(w, "Email or password incorrect", http.StatusInternalServerError)
			return
		}

		session, _ := store.Get(r, "user-session")
		session.Values["userID"] = queryResults.Id
		session.Save(r, w)
		fmt.Println("User logged in!")
		http.Redirect(w, r, "/tasks", http.StatusFound)
		return
	}

	// Loads the HTML template file
	tmpl, err := frontend.LoadHTML("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1 // Expire the session immediately
	session.Save(r, w)
	fmt.Println("User ID after logout:", session.Values["userID"])
}

// Creates session data
func init() {
	key := securecookie.GenerateRandomKey(32)
	os.Setenv("SESSION_KEY", hex.EncodeToString(key))
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func main() {
	db := database.Connect()
	defer db.Close()

	// Serve static files (including styles.css)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, db)
	})
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/tasks/update", func(w http.ResponseWriter, r *http.Request) {
		crudUpdate(w, r, db)
	})
	http.HandleFunc("/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		crudDelete(w, r, db)
	})
	http.HandleFunc("/tasks/create", func(w http.ResponseWriter, r *http.Request) {
		crudCreate(w, r, db)
	})
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		crudRead(w, r, db)
	})
	http.HandleFunc("/about", about)
    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
