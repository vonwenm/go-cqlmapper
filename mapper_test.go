package cqlmapper_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	cqlmapper "."
)

type MyTable struct {
	Id    uint
	Value string
}

type InstanceMapperTestSuite struct {
	suite.Suite
}

func (suite *InstanceMapperTestSuite) TestNewInstanceMapper() {
	myTable := &MyTable{}

	mapper, mapperErr := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.NotNil(suite.T(), mapper)
	assert.Nil(suite.T(), mapperErr)

	invalidMyTable := MyTable{}
	invalidMapper, invalidMapperError := cqlmapper.Underscore.NewInstanceMapper(invalidMyTable)

	assert.Nil(suite.T(), invalidMapper)
	assert.Equal(suite.T(), cqlmapper.NotPointerToStructError, invalidMapperError)
}

func (suite *InstanceMapperTestSuite) TestSelectQuery() {
	myTable := &MyTable{}
	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.Equal(suite.T(), "SELECT id, value FROM my_table", mapper.SelectQuery())
}

func (suite *InstanceMapperTestSuite) TestInsertQuery() {
	myTable := &MyTable{}

	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.Equal(suite.T(), "INSERT INTO my_table (id, value) VALUES(?, ?)", mapper.InsertQuery())
}

func (suite *InstanceMapperTestSuite) TestFields() {
	myTable := &MyTable{}
	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.Equal(suite.T(), []interface{}{&myTable.Id, &myTable.Value}, mapper.FieldPointers())
}

func TestInstanceMapper(t *testing.T) {
	suite.Run(t, new(InstanceMapperTestSuite))
}
