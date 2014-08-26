package context

import (
	//	"database/sql"
	httpcontext "github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Shared context that is available in web request handlers
// provides a db transaction context
type Context struct {
	db       *sqlx.DB
	tx       *sqlx.Tx
	commitTx bool
}

// Create a new Context for the db provided
func NewContext(db *sqlx.DB) *Context {
	return &Context{db: db, commitTx: true}
}

// Create a new Transactional context, or return the existing Transaction context
// if one has already been created
// Panic on error
func (c *Context) Begin() *sqlx.Tx {
	if c.tx == nil {
		c.tx = c.db.MustBegin()
	}
	return c.tx
}

// Mark the Transaction in the context for rollback when EndContext is called.
// If this is never called
// then it is assumed that the transaction will be committed on EndContext
func (c *Context) Rollback() {
	c.commitTx = false
}

// Close the Transactional context, by either committing or rolling back
// Panic on error
func (c *Context) End() {
	var err error = nil
	if c.tx != nil {
		if c.commitTx {
			err = c.tx.Commit()
		} else {
			err = c.tx.Rollback()
		}
	}
	if err != nil {
		panic(err)
	}
}

type key int

// Key used to store a database Connection in the request
const dbkey key = 87865544073

// GetDatabase returns the Context from the request values.
func GetDatabase(r *http.Request) *Context {
	if rv := httpcontext.Get(r, dbkey); rv != nil {
		return rv.(*Context)
	}
	panic("Missing dbkey in http context")
}

// SetDatabase sets a Context for this package in the request values.
func SetDatabase(r *http.Request, val *Context) {
	httpcontext.Set(r, dbkey, val)
}

// DeleteDatabase completes the request by committing or rolling back the tx, and
// then deletes the Context from the request values.
func EndDatabase(r *http.Request) {
	defer httpcontext.Delete(r, dbkey)
	GetDatabase(r).End()
}

// Handler that provides Database access to downstream handlers
// Use as follows
func DatabaseHandler(h http.Handler, db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := NewContext(db)
		SetDatabase(r, context)
		defer EndDatabase(r)
		h.ServeHTTP(w, r)
	})
}
