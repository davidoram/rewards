package context

import (
	//	"database/sql"
	httpcontext "github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Shared context that is available in web request handlers
// provides a db transaction context
type DatabaseContext struct {
	db       *sqlx.DB
	tx       *sqlx.Tx
	commitTx bool
}

// Create a new DatabaseContext for the db provided
func NewDatabaseContext(db *sqlx.DB) *DatabaseContext {
	return &DatabaseContext{db: db, commitTx: true}
}

// Create a new Transactional context, or return the existing Transaction context
// if one has already been created
// Panic on error
func (c *DatabaseContext) Begin() *sqlx.Tx {
	if c.tx == nil {
		c.tx = c.db.MustBegin()
	}
	return c.tx
}

// Mark the Transaction in the context for rollback when EndDatabaseContext is called.
// If this is never called
// then it is assumed that the transaction will be committed on EndDatabaseContext
func (c *DatabaseContext) Rollback() {
	c.commitTx = false
}

// Panic on err != nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Close the Transactional context, by either committing or rolling back
// Panic on error
func (c *DatabaseContext) End() {
	var err error = nil
	if c.tx != nil {
		if c.commitTx {
			err = c.tx.Commit()
		} else {
			err = c.tx.Rollback()
		}
	}
	check(err)
}

type key int

// Key used to store a database Connection in the request
const dbkey key = 87865544073

// GetDatabase returns the DatabaseContext from the request values.
func GetDatabase(r *http.Request) *DatabaseContext {
	if rv := httpcontext.Get(r, dbkey); rv != nil {
		return rv.(*DatabaseContext)
	}
	panic("Missing dbkey in http context")
}

// SetDatabase sets a DatabaseContext for this package in the request values.
func SetDatabase(r *http.Request, val *DatabaseContext) {
	httpcontext.Set(r, dbkey, val)
}

// DeleteDatabase completes the request by committing or rolling back the tx, and
// then deletes the DatabaseContext from the request values.
func EndDatabase(r *http.Request) {
	defer httpcontext.Delete(r, dbkey)
	GetDatabase(r).End()
}

// Handler that provides Database access to downstream handlers
// Use as follows
func DatabaseHandler(h http.Handler, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := NewDatabaseContext(db)
		SetDatabase(r, context)
		defer EndDatabase(r)
		h.ServeHTTP(w, r)
	})
}
