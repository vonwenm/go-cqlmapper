package cqlmapper_test

import (
	"os/exec"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/gocql/gocql"

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
CREATE TABLE IF NOT EXISTS example_table(id UUID, value TEXT, PRIMARY KEY(id));
INSERT INTO example_table(id, value) VALUES(uuid(), 'Test');
`))

var downTempalate = template.Must(template.New("down").Parse(`
DROP KEYSPACE IF EXISTS {{.Keyspace}};
`))

type ExampleTable struct {
	Id    gocql.UUID
	Value string
}

type ExamplesSuite struct {
	suite.Suite
	session *gocql.Session
}

func (suite *ExamplesSuite) SetupSuite() {
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

	go func() {
		if renderErr := tmpl.Execute(inPipe, cassandraConfig); nil != renderErr {
			panic(renderErr.Error())
		}
		inPipe.Close()
	}()

	if output, runErr := cqlshCmd.CombinedOutput(); nil != runErr {
		println(string(output))
		panic(runErr.Error())
	}
}

func (suite *ExamplesSuite) TestMapper_SelectQuery() {
	exampleTable := &ExampleTable{}
	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(exampleTable)

	query := suite.session.Query(mapper.SelectQuery())

	if scanErr := query.Scan(mapper.FieldPointers()...); nil != scanErr {
		panic(scanErr.Error())
	}

	assert.Equal(suite.T(), "Test", exampleTable.Value)
}

func (suite *ExamplesSuite) ExampleInstanceMapper_InsertQuery() {
	exampleTable := &ExampleTable{
		Id: gocql.TimeUUID(),
	}

	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(exampleTable)

	insertQuery := suite.session.Query(mapper.InsertQuery(), mapper.FieldPointers())

	if execErr := insertQuery.Exec(); nil != execErr {
		panic(execErr.Error())
	}
}

func TestExamples(t *testing.T) {
	suite.Run(t, new(ExamplesSuite))
}
