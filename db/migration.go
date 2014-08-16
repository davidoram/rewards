package db

import (
	//	"database/sql"
	"github.com/jmoiron/sqlx"
	//	"log"
	"path/filepath"
	"time"
)

type Migration struct {
	Filename  string
	CreatedAt time.Time `db:"created_at"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Runs each migration in order provided it hasnt been run prior, panics if the migrations fail.
// Pass the migration files in order. The filename component is stored so that we
// can check for future runs
func MustMigrate(db *sqlx.DB, paths *[]string) {
	db.MustExec(` CREATE TABLE IF NOT EXISTS migrations( 
                    filename      varchar(2048)   not null PRIMARY KEY, 
                    created_at    timestamp       not null default CURRENT_TIMESTAMP
                )`)
	for _, path := range *paths {
		if !IsMigrated(db, path) {
			//log.Logger.Println("Executing migration '%v'...", path)
			_, err := sqlx.LoadFile(db, path)
			check(err)
			setMigrated(db, path)
		}
	}
}

// Has the migration filename been applied? Note that all the
// file path information is stripped
func IsMigrated(db *sqlx.DB, filename string) bool {
	m := []Migration{}
	err := db.Select(&m,
		`SELECT   *
     FROM     migrations 
     WHERE    filename = $1`,
		filepath.Base(filename))
	check(err)
	return len(m) == 1
}

// Record filename as having been migrated, note that all the
// file path information is stripped
func setMigrated(db *sqlx.DB, filename string) {
	db.MustExec(`INSERT INTO migrations ( 
                              filename 
                          ) VALUES ( 
                              $1 
                          )`,
		filepath.Base(filename))
}
