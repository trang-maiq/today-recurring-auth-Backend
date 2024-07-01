package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func Migrate() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		sslMode  = os.Getenv("DB_SSL_MODE")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS identities (
        id SERIAL PRIMARY KEY,
        email VARCHAR(100) NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert initial data with hashed passwords only if it does not exist
	users := []struct {
		email    string
		password string
	}{
		{"user1@example.com", "password1"},
		{"user2@example.com", "password2"},
		{"user3@example.com", "password3"},
	}

	for _, user := range users {
		hashedPassword, err := hashPassword(user.password)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		insertUserSQL := `INSERT INTO identities (email, password)
                          VALUES ($1, $2)
                          ON CONFLICT (email) DO NOTHING`
		_, err = db.Exec(insertUserSQL, user.email, hashedPassword)
		if err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}
	}
}
