package main

import (
	"bytes"
	"fmt"
	"net"

	"github.com/ellezio/redict/internal/redict"
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
			var res resp.DataType
			switch string(cmd.Value) {
			case "SET":
				res = cmdSet(arr)
			case "GET":
				res = cmdGet(arr)
			default:
				fmt.Println("unknown command")
				return
			}
			fmt.Println("success")
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

func cmdSet(args *resp.Array) resp.DataType {
	var key string
	var value []byte

	if str, ok := args.Values[1].(*resp.BulkString); ok {
		key = string(str.Value)
	}

	if str, ok := args.Values[2].(*resp.BulkString); ok {
		value = str.Value
	}

	if err := db.Set(key, value); err != nil {
		fmt.Printf("cannot set value for store: %s\n", err)
		return resp.NewSimpleError(fmt.Sprintf("Err - cannot set value for store: %s\n", err))
	}

	return resp.NewSimpleString("OK")
}

func cmdGet(args *resp.Array) resp.DataType {
	var key string

	if str, ok := args.Values[1].(*resp.BulkString); ok {
		key = string(str.Value)
	}

	data, err := db.Get(key)
	if err != nil {
		fmt.Printf("cannot get data: %s", err)
		return resp.NewSimpleError(fmt.Sprintf("cannot get data: %s", err))
	}

	return resp.NewBulkString(string(data))
}
