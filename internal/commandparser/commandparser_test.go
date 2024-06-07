package commandparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InvalidCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("not a command")
	assert.Equal(t, UnknownCommand, res.Type)
}

func Test_UnknownCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("some-command('abc')")
	assert.Equal(t, UnknownCommand, res.Type)
}

func Test_MalformedCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("some-command('abc")
	assert.Equal(t, UnknownCommand, res.Type)
}

func Test_AddItemCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("addItem('key','value')")
	assert.Equal(t, AddItemCommand, res.Type)
	assert.Equal(t, "key", res.Key)
	assert.Equal(t, "value", res.Value)
}

func Test_AddItemCommandWithSpace(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("addItem('key', 'value')")
	assert.Equal(t, AddItemCommand, res.Type)
	assert.Equal(t, "key", res.Key)
	assert.Equal(t, "value", res.Value)
}

func Test_DeleteItemCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("deleteItem('key')")
	assert.Equal(t, DeleteItemCommand, res.Type)
	assert.Equal(t, "key", res.Key)
}

func Test_DeleteItemCommandWithSpace(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("deleteItem('key' )")
	assert.Equal(t, DeleteItemCommand, res.Type)
	assert.Equal(t, "key", res.Key)
}

func Test_GetItemCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("getItem('key')")
	assert.Equal(t, GetItemCommand, res.Type)
	assert.Equal(t, "key", res.Key)
}

func Test_GetAllItemsCommand(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("getAllItems()")
	assert.Equal(t, GetAllItemsCommand, res.Type)
}

func Test_GetAllItemsCommandWithSpace(t *testing.T) {
	sut := NewCommandParser()
	res := sut.ParseCommand("getAllItems( )")
	assert.Equal(t, GetAllItemsCommand, res.Type)
}
