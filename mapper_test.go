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

	assert.Equal(suite.T(), `SELECT "id", "value" FROM "my_table"`, mapper.SelectQuery())
}

func (suite *InstanceMapperTestSuite) TestInsertQuery() {
	myTable := &MyTable{}

	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.Equal(suite.T(), `INSERT INTO "my_table" ("id", "value") VALUES(?, ?)`, mapper.InsertQuery())
}

func (suite *InstanceMapperTestSuite) TestFields() {
	myTable := &MyTable{}
	mapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)

	assert.Equal(suite.T(), []interface{}{myTable.Id, myTable.Value}, mapper.FieldValues())
	assert.Equal(suite.T(), []interface{}{&myTable.Id, &myTable.Value}, mapper.FieldAddresses())
}

type RenamedTable struct {
}

func (*RenamedTable) TableName() string {
	return "new_table"
}

func (suite *InstanceMapperTestSuite) TestInstanceMapper_TableName() {
	myTable := new(MyTable)
	myMapper, _ := cqlmapper.Underscore.NewInstanceMapper(myTable)
	suite.Assertions.Equal(`"my_table"`, myMapper.TableName())

	renamedTable := new(RenamedTable)
	newMapper, _ := cqlmapper.Underscore.NewInstanceMapper(renamedTable)
	suite.Assertions.Equal(`"new_table"`, newMapper.TableName())
}

func TestInstanceMapper(t *testing.T) {
	suite.Run(t, new(InstanceMapperTestSuite))
}

type InnerStruct struct {
	Inner string
}

type OuterStruct struct {
	InnerStruct
	Outer string
}

func (suite *InstanceMapperTestSuite) TestFieldEmbedding() {
	outerStruct := &OuterStruct{}
	mapper, mapperErr := cqlmapper.Underscore.NewInstanceMapper(outerStruct)
	if nil != mapperErr {
		suite.T().Error(mapperErr)
	}

	assert.Equal(suite.T(), []string{`"inner"`, `"outer"`}, mapper.ColumnNames())
	assert.Equal(suite.T(), []interface{}{outerStruct.Inner, outerStruct.Outer}, mapper.FieldValues())
	assert.Equal(suite.T(), []interface{}{&outerStruct.Inner, &outerStruct.Outer}, mapper.FieldAddresses())
}
