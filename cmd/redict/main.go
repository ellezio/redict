package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/ellezio/redict/internal/redict"
	"github.com/ellezio/redict/internal/redict/command"
	"github.com/ellezio/redict/internal/resp"
)

var db *redict.Database

func main() {
	db = redict.NewDatabase()

	fmt.Println("Listening on :3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	handlers := map[command.CmdName]func(command.Cmd) (resp.DataType, error){
		command.SET: cmdSet,
		command.GET: cmdGet,

		command.LPUSH:  cmdLPush,
		command.RPUSH:  cmdRPush,
		command.LPOP:   cmdLPop,
		command.RPOP:   cmdRPop,
		command.LRANGE: cmdLRange,
		command.LTRIM:  cmdLTrim,
		command.LLEN:   cmdLLen,
	}

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
		} else {
			fmt.Println("executing command")

			var res resp.DataType = nil
			cmd, err := command.ParseCmd(arr)
			if err != nil {
				err = fmt.Errorf("cannot resolve command: %s", err)
				log.Println(err)
			} else if handler, ok := handlers[cmd.Name]; ok {
				res, err = handler(cmd)
			} else {
				log.Println("unknown command")
				return
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

func cmdSet(cmd command.Cmd) (resp.DataType, error) {
	if err := db.Set(cmd.Key, cmd.Value); err != nil {
		fmt.Printf("cannot set value for store: %s\n", err)
		return nil, fmt.Errorf("Err - cannot set value for store: %s\n", err)
	}

	return nil, nil
}

func cmdGet(cmd command.Cmd) (resp.DataType, error) {
	data, err := db.Get(cmd.Key)
	if err != nil {
		fmt.Printf("cannot get data: %s", err)
		return nil, fmt.Errorf("cannot get data: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdLPush(cmd command.Cmd) (resp.DataType, error) {
	if err := db.LPush(cmd.Key, cmd.Value); err != nil {
		log.Printf("cannot push value to list head: %s", err)
		return nil, fmt.Errorf("Err - cannot push value to list head: %s", err)
	}

	return nil, nil
}

func cmdRPush(cmd command.Cmd) (resp.DataType, error) {
	if err := db.RPush(cmd.Key, cmd.Value); err != nil {
		log.Printf("cannot push value to list tail: %s", err)
		return nil, fmt.Errorf("Err - cannot push value to list tail: %s", err)
	}

	return nil, nil
}

func cmdLPop(cmd command.Cmd) (resp.DataType, error) {
	data, err := db.LPop(cmd.Key)
	if err != nil {
		log.Printf("cannot pop list's head: %s", err)
		return nil, fmt.Errorf("Err - cannot pop list's head: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdRPop(cmd command.Cmd) (resp.DataType, error) {
	data, err := db.RPop(cmd.Key)
	if err != nil {
		log.Printf("cannot pop list's tail: %s", err)
		return nil, fmt.Errorf("Err - cannot pop list's tail: %s", err)
	}

	return resp.NewBulkString(string(data)), nil
}

func cmdLRange(cmd command.Cmd) (resp.DataType, error) {
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

func cmdLTrim(cmd command.Cmd) (resp.DataType, error) {
	if err := db.LTrim(cmd.Key, cmd.Start, cmd.End); err != nil {
		log.Printf("cannot trim list: %s", err)
		return nil, fmt.Errorf("Err - cannot trim list: %s", err)
	}

	return nil, nil
}

func cmdLLen(cmd command.Cmd) (resp.DataType, error) {
	data, err := db.LLen(cmd.Key)
	if err != nil {
		log.Printf("cannot trim list: %s", err)
		return nil, fmt.Errorf("Err - cannot trim list: %s", err)
	}

	return resp.NewInteger(int64(data)), nil
}
