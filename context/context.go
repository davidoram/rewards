package context

import (
	"database/sql"
	"net/http"
)

// Shared context that is available in web request handlers
// provides a db transaction context
type Context struct {
  db *sql.DB
  tx *sql.Tx
  commitTx bool
}

// Create a new Context for the db provided
func NewContext(db *sql.DB) (*Context) {
  return &Context{db: db, commitTx: true}
}


// Create a new Transactional context, or return the existing Transaction context
// if one has already been created
func (c *Context) Begin() (*sql.Tx, error) {
  if c.tx == nil {
		var err error
    c.tx, err = c.db.Begin()
    if err != nil {
      return nil, err
    }
  }
  return c.tx, nil
}

// Mark the Transaction in the context for rollback when EndContext is called.
// If this is never called
// then it is assumed that the transaction will be committed on EndContext
func (c *Context) Rollback() {
  c.commitTx = false
}

// Close the Transactional context, by either committing or rolling back
func (c *Context) End() (error) {
  var err error = nil
  if c.tx != nil {
    if c.commitTx {
      err = c.tx.Commit()
    } else {
      err = c.tx.Rollback()
    }
  }
  return err
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
