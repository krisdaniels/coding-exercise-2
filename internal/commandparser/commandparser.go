package commandparser

import (
	"regexp"
)

type CommandParser interface {
	ParseCommand(cmd string) *Command
}

var _ CommandParser = (*commandParser)(nil)

type commandParser struct{}

type CommandType int

const (
	UnknownCommand CommandType = iota
	AddItemCommand
	DeleteItemCommand
	GetItemCommand
	GetAllItemsCommand
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
}

func NewCommandParser() CommandParser {
	return &commandParser{}
}

// regex parses function name, and 1 of 3 param group combinations, 2, 1 or no params,
// allowing for whitespaces to be present between the params, ( ) and ,
var functionParser = regexp.MustCompile(`(.*?)\((\s*?'(.*?)'\s*?,\s*?'(.*?)'\s*?|\s*?'(.*?)'\s*?|\s*?)\)`)

func (s *commandParser) ParseCommand(cmd string) *Command {

	match := functionParser.FindStringSubmatch(cmd)

	res := &Command{}

	if len(match) < 6 {
		return res
	}

	switch match[1] {
	case "addItem":
		res.Type = AddItemCommand
		res.Key = match[3]
		res.Value = match[4]

	case "deleteItem":
		res.Type = DeleteItemCommand
		res.Key = match[5]

	case "getItem":
		res.Type = GetItemCommand
		res.Key = match[5]

	case "getAllItems":
		res.Type = GetAllItemsCommand

	}

	return res
}
