package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/ellezio/redict/internal/resp"
)

func main() {
	args := os.Args[1:]

	var cmd *resp.Array

	switch args[0] {
	case "SET":
		cmd = cmdSet(args[1:])
	case "GET":
		cmd = cmdGet(args[1:])
	case "LPUSH":
		cmd = cmdLPush(args[1:])
	case "RPUSH":
		cmd = cmdRPush(args[1:])
	case "LPOP":
		cmd = cmdLPop(args[1:])
	case "RPOP":
		cmd = cmdRPop(args[1:])
	case "LRANGE":
		cmd = cmdLRange(args[1:])
	case "LTRIM":
		cmd = cmdLTrim(args[1:])
	case "LLEN":
		cmd = cmdLLen(args[1:])
	default:
		fmt.Println("unknow command")
		return
	}

	var buf bytes.Buffer
	cmd.Encode(&buf)
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("command SET send to server")
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("waiting for response")
	res := make([]byte, 0, 1024)
	n, err := conn.Read(res[len(res):cap(res)])
	res = res[:len(res)+n]
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()
	fmt.Printf("res: %q\n", res)
}

func cmdSet(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("SET"))
	arr.Add(resp.NewBulkString(args[0]))
	arr.Add(resp.NewBulkString(args[1]))
	return arr
}

func cmdGet(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("GET"))
	arr.Add(resp.NewBulkString(args[0]))
	return arr
}

func cmdLPush(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("LPUSH"))
	arr.Add(resp.NewBulkString(args[0]))
	arr.Add(resp.NewBulkString(args[1]))
	return arr
}

func cmdRPush(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("RPUSH"))
	arr.Add(resp.NewBulkString(args[0]))
	arr.Add(resp.NewBulkString(args[1]))
	return arr
}

func cmdLPop(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("LPOP"))
	arr.Add(resp.NewBulkString(args[0]))
	return arr
}

func cmdRPop(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("RPOP"))
	arr.Add(resp.NewBulkString(args[0]))
	return arr
}

func cmdLRange(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("LRANGE"))
	arr.Add(resp.NewBulkString(args[0]))
	arr.Add(resp.NewBulkString(args[1]))
	arr.Add(resp.NewBulkString(args[2]))
	return arr
}

func cmdLTrim(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("LTRIM"))
	arr.Add(resp.NewBulkString(args[0]))
	arr.Add(resp.NewBulkString(args[1]))
	arr.Add(resp.NewBulkString(args[2]))
	return arr
}

func cmdLLen(args []string) *resp.Array {
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString("LLEN"))
	arr.Add(resp.NewBulkString(args[0]))
	return arr
}
