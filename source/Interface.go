package source

import "gopkg.in/rethinkdb/rethinkdb-go.v6"

type Interface interface {
	CheckReady() bool
	GetSession() *rethinkdb.Session
	GetTerm() rethinkdb.Term
	checkTerm() bool
	CreateTerm() error
	DropTerm() error
}
