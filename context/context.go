package context

import (
	//	"database/sql"
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

type ContextHandler struct {
	c Context
	f ContextFunc
}

type ContextFunc func(c Context, w http.ResponseWriter, r *http.Request)

func (h ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f(h.c, w, r)
}

/*
func f1(c context, w http.ResponseWriter, r *http.Request) {
    // Implement the handler here
}

func main() {
    c := NewContext(...) // Set up your context here
    http.Handle("/", ContextHandler{c, f1})
}
*/
