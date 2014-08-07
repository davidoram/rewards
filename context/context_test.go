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
		t.Fatalf("Error creating db with connection string '%v'. Should set envar DB_CONNECT=postgres://user@password/host/dbname?sslmode=disable ", os.Getenv("DB_CONNECT"))
		return nil
	}
	sql := `CREATE TABLE IF NOT EXISTS people
				(
					name varchar(255)
				)`
	_, err = db.Exec(sql)
	if err != nil {
		t.Fatalf("Error creating 'people' table")
		return nil
	}
	return db

}

func closeDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func countPeople(t *testing.T, db *sql.DB) int {
	var count int
	err := db.QueryRow("select count(*) from people").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	return count
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

func Test_Commit(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
	before := countPeople(t, db)
	context := NewContext(db)
	if context == nil {
		t.Fatal("Cant create context")
	}

	tx, err := context.Begin()
	if err != nil {
		t.Fatal("Failed on context.Begin")
	}
	_, err = tx.Exec("INSERT INTO people(name) VALUES ( 'Dave')")
	if err != nil {
		t.Fatal("Failed on insert")
	}
	err = context.End()
	if err != nil {
		t.Fatal("Failed on context.End")
	}

	after := countPeople(t, db)
	if after == (before + 1) {
		t.Log("Commit passed")
	} else {
		t.Error("Expected to insert 1 row, actually inserted ", (after - before))
	}
}

func Test_Rollback(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
	before := countPeople(t, db)
	context := NewContext(db)
	if context == nil {
		t.Fatal("Cant create context")
	}

	tx, err := context.Begin()
	if err != nil {
		t.Fatal("Failed on context.Begin")
	}
	_, err = tx.Exec("INSERT INTO people(name) VALUES ( 'Kerry')")
	if err != nil {
		t.Fatal("Failed on insert")
	}
	context.Rollback()
	
	err = context.End()
	if err != nil {
		t.Fatal("Failed on context.End")
	}

	after := countPeople(t, db)
	if after == before {
		t.Log("Rollback passed")
	} else {
		t.Error("Expected to rollback insert, but inserted ", (after - before))
	}
}
