package data

import "gopkg.in/rethinkdb/rethinkdb-go.v6"

type RethinkSource struct {
	Table   string
	Term    rethinkdb.Term
	Session *rethinkdb.Session
}
