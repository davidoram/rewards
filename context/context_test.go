package context

import (
	//	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func openDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_CONNECT"))
	if err != nil {
		t.Fatalf("Error creating db with connection string '%v'. Should set envar DB_CONNECT=postgres://user@password/host/dbname?sslmode=disable ", os.Getenv("DB_CONNECT"))
		return nil
	}
	db.MustExec("CREATE TABLE IF NOT EXISTS people ( name varchar(255) )")
	return db

}

func closeDB(db *sqlx.DB) {
	if db != nil {
		db.Close()
	}
}

func countPeople(t *testing.T, db *sqlx.DB) int {
	var count int
	db.QueryRowx("select count(*)  from people").Scan(&count)
	return count
}

func Test_Ping(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
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

	tx := context.Begin()
	tx.MustExec("INSERT INTO people(name) VALUES ( 'Dave')")
	context.End()

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

	tx := context.Begin()
	tx.MustExec("INSERT INTO people(name) VALUES ( 'Kerry')")
	context.Rollback()

	context.End()

	after := countPeople(t, db)
	if after == before {
		t.Log("Rollback passed")
	} else {
		t.Error("Expected to rollback insert, but inserted ", (after - before))
	}
}

func Test_MutiStatementCommit(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
	context := NewContext(db)

	tx := context.Begin()
	tx.MustExec("INSERT INTO people(name) VALUES ( 'Bob'), ('Rosa')")
	db.MustExec("DELETE FROM people")
	tx.MustExec("INSERT INTO people(name) VALUES ('Dave'),('Kerry'),('Jack'),('Tom')")
	context.End()

	after := countPeople(t, db)
	if after != 4 {
		t.Log("Commit passed")
	} else {
		t.Error("Expected 4 rows, actually have ", after)
	}
}

func Test_MutiStatementRollback(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)
	context := NewContext(db)

	// Start with an empty db
	db.MustExec("DELETE FROM people")

	tx := context.Begin()
	tx.MustExec("INSERT INTO people(name) VALUES ( 'Bob'), ('Rosa')")
	tx.MustExec("INSERT INTO people(name) VALUES ('Dave'),('Kerry'),('Jack'),('Tom')")
	context.Rollback()

	context.End()

	after := countPeople(t, db)
	if after == 0 {
		t.Log("Rollback passed")
	} else {
		t.Error("Expected 0 rows, actually have ", after)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	tx := GetDatabase(r).Begin()
	tx.MustExec("SELECT 1")
	fmt.Fprintln(w, "ok")
}

func Test_HttpPing(t *testing.T) {
	db := openDB(t)
	defer closeDB(db)

	ts := httptest.NewServer(DatabaseHandler(http.HandlerFunc(pingHandler), db))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	expected := "ok\n"
	if body != expected {
		t.Errorf("Response mismatch, got '%q', expected '%q'", body, expected)
	}

}
