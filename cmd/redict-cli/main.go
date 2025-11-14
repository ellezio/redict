package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/ellezio/redict/internal/redict/command"
	"github.com/ellezio/redict/internal/resp"
)

func main() {
	args := os.Args[1:]

	arr := resp.NewArray()
	meta, ok := command.GetMeta(command.CmdName(args[0]))
	if !ok {
		fmt.Println("unknow command")
		return
	}
	argsNumber := meta.ArgsNumber + 1

	if uint(len(args)) < argsNumber {
		fmt.Printf("invalid arguments number - expected %d given %d\n", argsNumber, len(args))
		return
	}

	for i := range argsNumber {
		arr.Add(resp.NewBulkString(args[i]))
	}

	var buf bytes.Buffer
	arr.Encode(&buf)
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
