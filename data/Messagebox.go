package data

import (
	"log"
)

type Messagebox struct {
	source      RethinkSource
	CreatedTime int64       `rethinkdb:"createdTime" json:"createdTime"`
	LastSeen    int64       `rethinkdb:"lastSeen" json:"lastSeen"`
	Metadata    interface{} `rethinkdb:"metadata" json:"metadata"`
	Origin      string      `rethinkdb:"origin" json:"origin"`
	Target      string      `rethinkdb:"target" json:"target"`
}

func (m *Messagebox) Get(listenerID string) error {
	cursor, err := m.source.Term.Table(listenerID).Get(m.Target).Run(m.source.Session)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	return cursor.One(m)
}

func (m *Messagebox) Delete(listenerID string) error {
	return m.source.Term.Table(listenerID).Get(m.Target).Delete().Exec(m.source.Session)
}
