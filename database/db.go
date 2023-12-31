package database

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type connection struct {
	host string
	user string 
	port string 
	password string 
	dbname string 
}


// Connects to the database
func Connect() *sql.DB {
	var psqlInfo string
	// Tries connecting to render.com db
	if os.Getenv("RENDER_POSTGRES_URL") != "" {
		psqlInfo = os.Getenv("RENDER_POSTGRES_URL")
	} else { // or to local postgresql db
		fmt.Println("Connecting to local db")
		var conn connection
		conn.host = os.Getenv("DB_HOST")
		conn.user = os.Getenv("DB_USER")
		conn.port = os.Getenv("DB_PORT")
		conn.password = os.Getenv("DB_PW")
		conn.dbname = os.Getenv("DB_NAME")

		psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conn.host, conn.port, conn.user, conn.password, conn.dbname)

	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected successfully!")

	return db
}


func SelectManyQuery(db *sql.DB, userID int) []TaskQuery {
	var results []TaskQuery

	// rows is the result of our query
	rows, err := db.Query("SELECT * FROM tasks WHERE user_id = $1;", userID)
	if err != nil {
	  // handle this error better than this
	  fmt.Println("if error appears here...")
	  panic(err)
	}
	defer rows.Close()

	// iterates results
	for rows.Next() {
		var row TaskQuery
		err = rows.Scan(&row.Id, &row.Title, &row.Description, &row.DueDate, &row.User_id)

		if err != nil {
			// handle this error
			fmt.Println("or if it appears here...")
			panic(err)
		}
		results = append(results, row)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return results
}

// Selects single row from query
func SelectSingleQuery(db *sql.DB, id int, userId int) (TaskQuery, error) {
	var row TaskQuery
	query := `SELECT * FROM tasks WHERE id = $1 AND user_id = $2;`
	fmt.Println("DEBUGGING QUERY: ", query, id, userId)
	qr := db.QueryRow(query, id, userId)
	err := qr.Scan(&row.Id,	&row.Title,	&row.Description, &row.DueDate,	&row.User_id)
	if err != nil {
		return row, err
	}
	return row, err
}


func SelectLoginQuery(db *sql.DB, loginForm LoginFormData, requestType int) (LoginQuery, error) {
	var row LoginQuery

	switch requestType {
		case 1:
			query := `SELECT * FROM users WHERE email=$1`
			qr := db.QueryRow(query, loginForm.Email)
			err := qr.Scan(&row.Id, &row.Email, &row.Password)
			if err != nil {
				return row, err
			}
		case 2:
			query := `SELECT * FROM users WHERE email=$1 AND password=$2`
			qr := db.QueryRow(query, loginForm.Email, loginForm.Password)
			err := qr.Scan(&row.Id, &row.Email, &row.Password)
			if err != nil {
				return row, err
			}
	}
	return row, nil
}

func InsertIntoUsers(db *sql.DB, data RegisterFormData) {
	query := `INSERT INTO users
	(email, password)
	VALUES ($1, $2);`
	_, err := db.Exec(query, 
		data.Email,
		data.Password1,
	)
	if err != nil {
		fmt.Println("ERROR INSERTING:", err)
	}
}


func InsertIntoTasks(db *sql.DB, data TaskQuery, userId int) error {
	query := `INSERT INTO tasks
	(title, description, dueDate, user_id)
	VALUES ($1, $2, $3, $4);`
	if data.Title == "" {
		data.Title = "Unnamed task"
	}
	_, err := db.Exec(query, 
		data.Title,
		data.Description,
		data.DueDate,
		userId,
	)
	if err != nil {
		fmt.Println("ERROR INSERTING:", err)
		return err
	}
	return nil
}

func DeleteTask(db *sql.DB, data TaskQuery, userId int) { //data TaskQuery) {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	_, err := db.Exec(query, data.Id, userId)
	if err != nil{
		fmt.Println(err)
	}
}

func UpdateTask(db *sql.DB, updateData TaskQuery, userId int) {
	/*
	UPDATE tasks SET
	title = 'iloveyou',
	description = 'somuch'
	WHERE id = 5;
	*/
	query := `
		UPDATE tasks SET
		title = $2,
		description = $3,
		duedate = $4
		WHERE id = $1
		AND user_id = $5;
	`
	_, err := db.Exec(query,
		updateData.Id,
		updateData.Title,
		updateData.Description,
		updateData.DueDate,
		userId)

	if err != nil {
		fmt.Println(err)
	}
}