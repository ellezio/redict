package command

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ellezio/redict/internal/resp"
)

type CmdName string

const (
	SET CmdName = "SET"
	GET CmdName = "GET"

	LPUSH  CmdName = "LPUSH"
	RPUSH  CmdName = "RPUSH"
	LPOP   CmdName = "LPOP"
	RPOP   CmdName = "RPOP"
	LRANGE CmdName = "LRANGE"
	LTRIM  CmdName = "LTRIM"
	LLEN   CmdName = "LLEN"
)

type CmdMeta struct {
	Pattern    []string
	ArgsNumber uint
}

var commands = map[CmdName]CmdMeta{
	SET: {[]string{"KEY", "VALUE"}, 2},
	GET: {[]string{"KEY"}, 1},

	LPUSH:  {[]string{"KEY", "VALUE"}, 2},
	RPUSH:  {[]string{"KEY", "VALUE"}, 2},
	LPOP:   {[]string{"KEY"}, 1},
	RPOP:   {[]string{"KEY"}, 1},
	LRANGE: {[]string{"KEY", "START", "END"}, 3},
	LTRIM:  {[]string{"KEY", "START", "END"}, 3},
	LLEN:   {[]string{"KEY"}, 1},
}

func GetMeta(name CmdName) (CmdMeta, bool) {
	meta, ok := commands[name]
	return meta, ok
}

type Cmd struct {
	Name  CmdName
	Key   string
	Value []byte
	Start int64
	End   int64
}

func ParseCmd(args *resp.Array) (Cmd, error) {
	name, ok := args.Values[0].(*resp.BulkString)
	if !ok {
		return Cmd{}, errors.New("invalid type of command name")
	}

	meta, ok := commands[CmdName(name.Value)]
	if !ok {
		return Cmd{}, errors.New("pattern not found")
	}

	if len(args.Values) < int(meta.ArgsNumber) {
		return Cmd{}, errors.New("invalid arguments number")
	}

	cmd := Cmd{}
	cmd.Name = CmdName(name.Value)

	idx := 1
	for _, part := range meta.Pattern {
		bstr, ok := args.Values[idx].(*resp.BulkString)
		if !ok {
			return Cmd{}, fmt.Errorf("invalid type for %q", part)
		}

		var err error
		switch part {
		case "KEY":
			cmd.Key = string(bstr.Value)
		case "VALUE":
			cmd.Value = bstr.Value
		case "START":
			cmd.Start, err = toInt64(bstr)
		case "END":
			cmd.End, err = toInt64(bstr)
		default:
			return Cmd{}, fmt.Errorf("unknown pattern part %q", part)
		}

		if err != nil {
			return Cmd{}, err
		}

		idx++
	}

	return cmd, nil
}

func toInt64(v *resp.BulkString) (int64, error) {
	i, err := strconv.ParseInt(string(v.Value), 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot convert BulkString to int64: %s", err)
	}
	return i, err
}
