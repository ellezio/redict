package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/ellezio/redict/internal/redict"
	"github.com/ellezio/redict/internal/resp"
)

var db *redict.Database

type Patterns map[string]string

var patterns = Patterns{
	"SET": "KEY VALUE",
	"GET": "KEY",

	"LPUSH":  "KEY VALUE",
	"RPUSH":  "KEY VALUE",
	"LPOP":   "KEY",
	"RPOP":   "KEY",
	"LRANGE": "KEY START END",
	"LTRIM":  "KEY START END",
	"LLEN":   "KEY",
}

func (p Patterns) match(args *resp.Array) (Cmd, error) {
	cmd := Cmd{}

	cmdName, ok := args.Values[0].(*resp.BulkString)
	if !ok {
		return Cmd{}, errors.New("invalid type of command name")
	}

	pattern, ok := p[string(cmdName.Value)]
	if !ok {
		return Cmd{}, errors.New("pattern not found")
	}

	parts := strings.Split(pattern, " ")
	argsLen := len(args.Values)
	idx := 1
	for _, part := range parts {
		if idx >= argsLen && p.isPartRequired(part) {
			return Cmd{}, fmt.Errorf("cannot read %q - invalid arguments number", part)
		}

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
			cmd.Start, err = p.toInt64(bstr)
		case "END":
			cmd.End, err = p.toInt64(bstr)
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

func (p Patterns) isPartRequired(part string) bool {
	return part[0] != '['
}

func (p Patterns) toInt64(v *resp.BulkString) (int64, error) {
	i, err := strconv.ParseInt(string(v.Value), 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot convert BulkString to int64: %s", err)
	}
	return i, err
}

type Cmd struct {
	Key   string
	Value []byte
	Start int64
	End   int64
}

func main() {
	db = redict.NewDatabase()

	fmt.Println("Listening on :3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		fmt.Println("command received")

		fmt.Println("reading payload")
		payload := make([]byte, 0, 1024)
		n, err := conn.Read(payload[len(payload):cap(payload)])
		payload = payload[:len(payload)+n]
		if err != nil {
			fmt.Println(err)
			return
		}

		dt, _, err := resp.Decode(payload)
		if err != nil {
			fmt.Println("decode -", err)
			return
		}

		if arr, ok := dt.(*resp.Array); !ok {
			fmt.Println("invalid payload type")
			continue
		} else if cmd, ok := arr.Values[0].(*resp.BulkString); !ok {
			fmt.Println("invalid payload type")
			continue
		} else {
			fmt.Println("executing command")

			var res resp.DataType = nil
			cmdArgs, err := patterns.match(arr)
			if err != nil {
				err = fmt.Errorf("cannot resolve command: %s", err)
				log.Println(err)
			} else {
				switch string(cmd.Value) {
				case "SET":
					err = cmdSet(cmdArgs)
				case "GET":
					res, err = cmdGet(cmdArgs)
				case "LPUSH":
					err = cmdLPush(cmdArgs)
				case "RPUSH":
					err = cmdRPush(cmdArgs)
				case "LPOP":
					res, err = cmdLPop(cmdArgs)
				case "RPOP":
					res, err = cmdRPop(cmdArgs)
				case "LRANGE":
					res, err = cmdLRange(cmdArgs)
				case "LTRIM":
					err = cmdLTrim(cmdArgs)
				case "LLEN":
					res, err = cmdLLen(cmdArgs)
				default:
					log.Println("unknown command")
					return
				}
			}

			if err != nil {
				res = resp.NewSimpleError(err.Error())
			} else if res == nil {
				res = resp.NewSimpleString("OK")
			}

			var buf bytes.Buffer
			res.Encode(&buf)
			_, err = conn.Write(buf.Bytes())
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("response sent")
		}
	}
}

func cmdSet(cmd Cmd) error {
	if err := db.Set(cmd.Key, cmd.Value); err != nil {
		fmt.Printf("cannot set value for store: %s\n", err)
		return fmt.Errorf("Err - cannot set value for store: %s\n", err)
	}

	return nil
}

func cmdGet(cmd Cmd) (resp.DataType, error) {
	data, err := db.Get(cmd.Key)
	if err != nil {
		fmt.Printf("cannot get data: %s", err)
		return nil, fmt.Errorf("cannot get data: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdLPush(cmd Cmd) error {
	if err := db.LPush(cmd.Key, cmd.Value); err != nil {
		log.Printf("cannot push value to list head: %s", err)
		return fmt.Errorf("Err - cannot push value to list head: %s", err)
	}

	return nil
}

func cmdRPush(cmd Cmd) error {
	if err := db.RPush(cmd.Key, cmd.Value); err != nil {
		log.Printf("cannot push value to list tail: %s", err)
		return fmt.Errorf("Err - cannot push value to list tail: %s", err)
	}

	return nil
}

func cmdLPop(cmd Cmd) (resp.DataType, error) {
	data, err := db.LPop(cmd.Key)
	if err != nil {
		log.Printf("cannot pop list's head: %s", err)
		return nil, fmt.Errorf("Err - cannot pop list's head: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdRPop(cmd Cmd) (resp.DataType, error) {
	data, err := db.RPop(cmd.Key)
	if err != nil {
		log.Printf("cannot pop list's tail: %s", err)
		return nil, fmt.Errorf("Err - cannot pop list's tail: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdLRange(cmd Cmd) (resp.DataType, error) {
	data, err := db.LRange(cmd.Key, cmd.Start, cmd.End)
	if err != nil {
		log.Printf("cannot get range of list: %s", err)
		return nil, fmt.Errorf("Err - cannot get range of list: %s", err)
	}

	arr := resp.NewArray()
	for _, d := range data {
		arr.Add(resp.NewBulkString(string(d)))
	}

	return arr, nil
}

func cmdLTrim(cmd Cmd) error {
	if err := db.LTrim(cmd.Key, cmd.Start, cmd.End); err != nil {
		log.Printf("cannot trim list: %s", err)
		return fmt.Errorf("Err - cannot trim list: %s", err)
	}

	return nil
}

func cmdLLen(cmd Cmd) (resp.DataType, error) {
	data, err := db.LLen(cmd.Key)
	if err != nil {
		log.Printf("cannot trim list: %s", err)
		return nil, fmt.Errorf("Err - cannot trim list: %s", err)
	}

	return resp.NewInteger(int64(data)), nil
}
