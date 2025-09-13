# Go CRUD API for School Database

Hey there! This is a simple CRUD (Create, Read, Update, Delete) API I built while learning Go. It's for managing student records in a school database—think registering kids with names and roll numbers, checking if they're enrolled, updating names, or unregistering them. I used Neon for the Postgres DB because it's free and easy for beginners. Nothing fancy, just a solid little project to practice HTTP handlers, JSON, and SQL queries.

## Features
- **Register Student**: POST a name and roll no to add a new student (checks for duplicate roll nos).
- **Check Student**: GET with roll no to see if they're in the DB (returns details if found).
- **Update Name**: PUT with roll no and new name to change it.
- **Delete Student**: DELETE with roll no to remove them.
- Basic error handling (e.g., 404 if not found, 409 for duplicates).

It's all in one `main.go` file for simplicity— no frameworks, just standard Go libs.

## Tech Stack
- Go (net/http for server, database/sql for DB)
- PostgreSQL (hosted on Neon—free tier works great)
- JSON for requests/responses

## Setup
1. **Clone the Repo**: `git clone https://github.com/yourusername/your-repo-name.git` (replace with your repo URL).
2. **Install Dependencies**: Run `go mod tidy` to get lib/pq for Postgres.
3. **Database**: Sign up for Neon (free), create a DB, and update the connection string in `main.go` (line 21) with your own (keep it secret!).
4. **Run Locally**: `go run main.go` — server starts on port 8080.

Note: For production, move the conn string to env vars (e.g., `os.Getenv("POSTGRES_URL")`) to avoid hardcoding secrets.

## Usage
Use Postman or curl to test the endpoints. Base URL: `http://localhost:8080`.

- **Register (POST /register)**:
  ```
  curl -X POST -H "Content-Type: application/json" -d '{"name": "Aditi", "rollno": 42}' http://localhost:8080/register
  ```
  - Success: "Student registered successfully" (201 if you add status).

- **Check (GET /check?rollno=42)**:
  ```
  curl http://localhost:8080/check?rollno=42
  ```
  - If found: "Student found: Aditi" (or JSON if updated).

- **Update (PUT /update?rollno=42)**:
  ```
  curl -X PUT -H "Content-Type: application/json" -d '{"name": "Aditi Sharma"}' http://localhost:8080/update?rollno=42
  ```
  - Success: "Student updated successfully".

- **Delete (DELETE /delete?rollno=42)**:
  ```
  curl -X DELETE http://localhost:8080/delete?rollno=42
  ```
  - Success: "Student deleted successfully" (404 if not found).

Check your Neon dashboard to see the data changes.


## Lessons Learned
This was my first Go API—loved how simple http and sql are, but handling errors and conversions took some trial and error. If you're learning, start with the register handler and build from there!


