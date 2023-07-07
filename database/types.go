package database

import (
	"time"
)

// Define tasks table columns into a struct
type TaskQuery struct {
	Id int
	Title string 
	Description string
	DueDate time.Time
	User_id int
}

// Define tasks table columns into a struct
type LoginQuery struct {
	Id int
	Email string
	Password string 
}

type LoginFormData struct {
	Email string
	Password string
}

type RegisterFormData struct {
	Email string 
	Password1 string 
	Password2 string
}

type Task struct {
	ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
	DueDate     time.Time `json:"duedate"`
}

type CreateTaskResponse struct {
    Task    Task   `json:"task"`
    Message string `json:"message"`
}