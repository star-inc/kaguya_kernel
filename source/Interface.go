package source

import "gopkg.in/rethinkdb/rethinkdb-go.v6"

type Interface interface {
	CheckReady() bool
	GetTerm() rethinkdb.Term
	GetSession() *rethinkdb.Session
}
