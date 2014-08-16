package db

import (
	//	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"testing"
)

func openDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_CONNECT"))
	if err != nil {
		t.Fatalf("Error creating db with connection string '%v'. Should set envar DB_CONNECT=postgres://user@password/host/dbname?sslmode=disable ", os.Getenv("DB_CONNECT"))
		return nil
	}
	// Allways start with clean db!
	db.MustExec("DROP TABLE IF EXISTS migrations")

	empty := []string{}
	MustMigrate(db, &empty)
	return db
}

func writeFile(path string, text string) {
	b := []byte(text)
	err := ioutil.WriteFile(path, b, 0644)
	check(err)
}

// func countPeople(t *testing.T, db *sqlx.DB) int {
// 	var count int
// 	db.QueryRowx("select count(*)  from people").Scan(&count)
// 	return count
// }

func Test_Ping(t *testing.T) {
	db := openDB(t)
	defer db.Close()
	t.Log("ping test passed.")
}

func Test_FirstMigration(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	filename := "/tmp/001_migration_test.sql"
	writeFile(filename, "select 1")
	defer os.Remove(filename)

	if IsMigrated(db, filename) {
		t.Error("Expected '", filename, "' not to have been applied yet")
	}

	MustMigrate(db, &[]string{filename})

	if !IsMigrated(db, filename) {
		t.Error("Expected '", filename, "' migration to have been applied")
	}

}
