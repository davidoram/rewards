package context

import (
	"testing"
  "database/sql"
  _ "github.com/lib/pq"
	"os"
)

func openDB(t *testing.T) (*sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DB_CONNECT"))
	if err != nil {
		t.Fatalf("Error creating database with connection string '%v'. Should set envar DB_CONNECT=postgres://user@password/host/dbname?sslmode=disable ", os.Getenv("DB_CONNECT"))
		return nil
	}
	return db

}

func closeDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func Test_Ping(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
	err := db.Ping()
	if err != nil {
		t.Fatalf("Ping returned error, %v", err)
	}

	t.Log("ping test passed.")
}

func Test_Fail(t *testing.T) {
	t.Error("this is just hardcoded as an error.")
}
