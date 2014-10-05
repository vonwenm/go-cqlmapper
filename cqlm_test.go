package cqlm_test

import (
	"testing"

	"github.com/gocql/gocql"

	cqlm "."
)

const (
	host     = "localhost"
	keyspace = "cqlm_test"
)

var session = mustCqlSession()

func mustCqlSession() *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	session, sessionErr := cluster.CreateSession()
	if nil != sessionErr {
		panic(sessionErr.Error())
	}
	return session
}

type RawStruct struct {
	Id        gocql.UUID `cqlm:"id"`
	TextField string     `cqlm:"text_field"`
	IntField  int        `cqlm:"int_field"`
}

func TestRawSelect(t *testing.T) {
	rawStruct := RawStruct{}
	query := cqlm.DefaultMapper.Select(&rawStruct)

	if scanErr := query.Scan(session); nil != scanErr {
		t.Error(scanErr)
	}

}

type UnderscoreStruct struct {
	Id        gocql.UUID
	TextField string
	IntField  int
}

func TestUnderscoreSelect(t *testing.T) {
	underscoreStruct := UnderscoreStruct{}
	query := cqlm.UnderscoreMapper.Select(&underscoreStruct)
	if scanErr := query.Scan(session); nil != scanErr {
		t.Error(scanErr)
	}
}
