package context

import (
	"database/sql"
)

// Shared context that is available in web request handlers
// provides a database transaction context
type Context struct {
  sql.DB *database
  sql.Tx *tx
  bool shouldCommit
}

// Create a new Context for the database provided
func NewContext(sql.DB *database) Context {
  Context c = new(Context)
  c.database = database
  c.shouldCommit = true
  return &c
}

// Create a new Transactional context, or return the existing Transaction context
// if one has already been created
func (c Context) Begin (*Tx, error) {
  if c.tx == nil {
    c.tx, err := c.database.Begin()
    if err != nil {
      return nil, err
    }
  }
  return c.tx, nil
}

// Mark the Transaction in the context for rollback when EndContext is called.
// If this is never called
// then it is assumed that the transaction will be committed on EndContext
func (c Context) Rollback {
  c.shouldCommit = false
}

// Close the Transactional context, by either committing or rolling back
func (c Context) EndContext (error) {
  var err error = nil
  if c.tx != nil {
    if c.shouldCommit {
      err = c.tx.Commit()
    }
    else {
      err = c.tx.Rollback()
    }
  }
  return err
}

type ContextHandler struct {
	c context
	f contextFunc
}

type ContextFunc func(c context, w http.ResponseWriter, r *http.Request)

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
