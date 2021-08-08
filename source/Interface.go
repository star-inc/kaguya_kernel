package source

import "gopkg.in/rethinkdb/rethinkdb-go.v6"

type Interface interface {
	GetTerm() rethinkdb.Term
	GetSession() *rethinkdb.Session
}
