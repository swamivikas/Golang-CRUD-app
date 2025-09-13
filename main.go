package main


import (
	"fmt"
	"encoding/json"
	"net/http"
	"database/sql"
	"log"
	"strings"
	"strconv"
	_ "github.com/lib/pq"
	
)


var db *sql.DB // Global for handlers to use

func initDB() {
    // Hardcoded connection string (REPLACE WITH YOUR ACTUAL STRING BELOW)
    connStr := "postgresql://neondb_owner:pass@ep-mute-brook-ad5mtu2p-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Open error: ", err)
    }

    // Ping to verify connection
    if err = db.Ping(); err != nil {
        log.Fatal("Ping error: ", err)
    }

    log.Println("Connected to Neon DB successfully")
}


func main() {
	initDB() 
    defer db.Close()
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/check", checkPage)
	http.HandleFunc("/delete", deleteStudent)
	http.HandleFunc("/update", updateStudent)


	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}


func landingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
	w.Write([]byte("welcome to the landing page"))
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	// here we will get name and roll no from the student and insert it into the database
	fmt.Println(r.Body)
	fmt.Println(r.Header)

	var student struct {
		Name string `json:"name"`
		RollNo int `json:"rollno"`
	}

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusBadRequest)
    	return
	}

	fmt.Println(student.Name, student.RollNo)
	_, err = db.Exec("INSERT INTO students (name, roll_no) VALUES ($1, $2)", student.Name, student.RollNo)
	if err != nil {
		// Handle duplicates
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			http.Error(w, "RollNo already exists", http.StatusConflict)
			return
		}
		log.Println("Insert error: ", err) // Debug log
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}


	fmt.Fprintf(w, "Student registered successfully")


}

func checkPage(w http.ResponseWriter, r *http.Request) {
	// here we will get roll no from the student and check if the student is present in the database
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rollNo := r.URL.Query().Get("rollno")
	if rollNo == "" {
		http.Error(w, "Roll no is required", http.StatusBadRequest)
		return
	}

	rollNoInt, err := strconv.Atoi(rollNo)
	if err != nil || rollNoInt <= 0 {
    	http.Error(w, "Invalid roll_no (must be positive integer)", http.StatusBadRequest)
    	return
	}
	row := db.QueryRow("SELECT name, roll_no FROM students WHERE roll_no = $1", rollNoInt)
	var student struct {
		Name string `json:"name"`
		RollNo int `json:"rollno"`
	}
	err = row.Scan(&student.Name, &student.RollNo)
	fmt.Println(student.Name, student.RollNo)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Student found: %s", student.Name)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	// here we will get roll no from the student and delete the student from the database
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rollNo := r.URL.Query().Get("rollno")
	if rollNo == "" {
		http.Error(w, "Roll no is required", http.StatusBadRequest)
		return
	}

	rollNoInt, err := strconv.Atoi(rollNo)
	if err != nil || rollNoInt <= 0 {
		http.Error(w, "Invalid roll_no (must be positive integer)", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM students WHERE roll_no = $1", rollNoInt)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Student deleted successfully")
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	// here we will get roll no and name from the student and update the student in the database
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rollNo := r.URL.Query().Get("rollno")
	if rollNo == "" {
		http.Error(w, "Roll no is required", http.StatusBadRequest)
		return
	}

	rollNoInt, err := strconv.Atoi(rollNo)
	if err != nil || rollNoInt <= 0 {
		http.Error(w, "Invalid roll_no (must be positive integer)", http.StatusBadRequest)
		return
	}
	var student struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE students SET name = $1 WHERE roll_no = $2", student.Name, rollNoInt)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Student updated successfully")
}

	
