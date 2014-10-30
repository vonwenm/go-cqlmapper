package cqlmapper_test

import (
	"fmt"
	"os/exec"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/probkiizokna/gocql"

	cqlmapper "."
)

var cassandraConfig = struct {
	Hosts    []string
	Keyspace string
}{
	Hosts:    []string{"localhost"},
	Keyspace: "cqlmapper_test",
}

var upTemplate = template.Must(template.New("up").Parse(`
DROP KEYSPACE IF EXISTS {{.Keyspace}};
CREATE KEYSPACE IF NOT EXISTS {{.Keyspace}} WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };
USE {{.Keyspace}};
CREATE TABLE IF NOT EXISTS example_table(id UUID, value TEXT, value2 INT, value3 int, PRIMARY KEY(id));
INSERT INTO example_table(id, value, value2, value3) VALUES({{.Id}}, 'Test', 123, null);
`))

var downTempalate = template.Must(template.New("down").Parse(`
DROP KEYSPACE IF EXISTS {{.Keyspace}};
`))

type ExampleTable struct {
	Id     gocql.UUID
	Value  string
	Value2 *int
	Value3 *uint
}

type ExamplesSuite struct {
	suite.Suite
	session *gocql.Session
	id      gocql.UUID
}

func (suite *ExamplesSuite) SetupSuite() {
	suite.id = gocql.TimeUUID()
	suite.executeTemplate(upTemplate)

	cluster := gocql.NewCluster(cassandraConfig.Hosts...)
	cluster.Keyspace = cassandraConfig.Keyspace
	if session, sessionErr := cluster.CreateSession(); nil == sessionErr {
		suite.session = session
	} else {
		panic(sessionErr.Error())
	}
}

func (suite *ExamplesSuite) TearDownSuite() {
	suite.executeTemplate(downTempalate)
}

func (suite *ExamplesSuite) executeTemplate(tmpl *template.Template) {
	cqlshCmd := exec.Command("cqlsh")

	inPipe, inPipeErr := cqlshCmd.StdinPipe()
	if nil != inPipeErr {
		panic(inPipeErr.Error())
	}

	if startErr := cqlshCmd.Start(); nil != startErr {
		panic(startErr.Error())
	}

	data := struct {
		Keyspace string
		Id       gocql.UUID
	}{
		Keyspace: cassandraConfig.Keyspace,
		Id:       suite.id,
	}

	if renderErr := tmpl.Execute(inPipe, data); nil != renderErr {
		panic(renderErr.Error())
	}
	inPipe.Close()

	if waitErr := cqlshCmd.Wait(); nil != waitErr {
		panic(waitErr.Error())
	}
}

func (suite *ExamplesSuite) TestInstanceMapper_SelectQuery() {
	exampleTable := &ExampleTable{}
	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(exampleTable)

	query := suite.session.Query(
		fmt.Sprintf(
			"%s WHERE id = ?",
			mapper.SelectQuery(),
		),
		suite.id,
	)

	if scanErr := query.Scan(mapper.FieldAddresses()...); nil != scanErr {
		panic(scanErr.Error())
	}

	assert.Equal(suite.T(), suite.id, exampleTable.Id)
	assert.Equal(suite.T(), "Test", exampleTable.Value)
	if assert.NotNil(suite.T(), exampleTable.Value2) {
		assert.Equal(suite.T(), 123, *exampleTable.Value2)
	}
	assert.Nil(suite.T(), exampleTable.Value3)
}

func (suite *ExamplesSuite) TestInstanceMapper_InsertQuery() {
	exampleTable := &ExampleTable{
		Id:    gocql.TimeUUID(),
		Value: "NewValue",
	}

	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(exampleTable)

	insertQuery := suite.session.Query(mapper.InsertQuery(), mapper.FieldValues()...)

	if execErr := insertQuery.Exec(); nil != execErr {
		panic(execErr.Error())
	}

}

func TestExamples(t *testing.T) {
	suite.Run(t, new(ExamplesSuite))
}
