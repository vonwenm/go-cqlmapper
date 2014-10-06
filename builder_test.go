package cqlmapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cqlmapper "."
)

type RawStruct struct {
	Id    uint
	Value string
}

func TestDefaultRaw(t *testing.T) {
	rawStruct := &RawStruct{}

	rawMapper, mapperErr := cqlmapper.Default.NewInstanceMapper(rawStruct)

	assert.Nil(t, mapperErr)
	assert.Equal(t, "RawStruct", rawMapper.TableName())
	assert.Equal(t, []string{"Id", "Value"}, rawMapper.ColumnNames())
}

type TaggedStruct struct {
	Id    uint   `cqlm:"id"`
	Value string `cqlm:"value"`
}

func TestDefaultTagged(t *testing.T) {
	taggedStruct := &TaggedStruct{}

	taggedMapper, mapperErr := cqlmapper.Default.NewInstanceMapper(taggedStruct)

	assert.Nil(t, mapperErr)
	assert.Equal(t, []string{"id", "value"}, taggedMapper.ColumnNames())
}

type UnderscoreStruct struct {
	Id    uint
	Value string
}

func TestUnderscore(t *testing.T) {
	underscoreStruct := &UnderscoreStruct{}

	underscoreMapper, mapperErr := cqlmapper.Underscore.NewInstanceMapper(underscoreStruct)

	assert.Nil(t, mapperErr)
	assert.Equal(t, "underscore_struct", underscoreMapper.TableName())
	assert.Equal(t, []string{"id", "value"}, underscoreMapper.ColumnNames())
}
