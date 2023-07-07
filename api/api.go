package api

import (
	"fmt"
	"encoding/json"
)

func main() {
	fmt.Println("Hello World!")
}

func APICreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    // Parse the JSON request body into a struct
    var task database.Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Failed to parse JSON request", http.StatusBadRequest)
        return
    }
	fmt.Println("TESTING:", task)

    // Insert the task into the database
	// Must update userID
    err = database.InsertIntoTasks(db, task, 1)
    if err != nil {
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    // Return a success response
    response := database.CreateTaskResponse{
		Task: task,
        Message: "Task created successfully",
    }

    // Serialize the response data to JSON
    jsonBytes, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
        return
    }

    // Set the response headers and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}