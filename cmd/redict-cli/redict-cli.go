package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ellezio/redict/internal/resp"
)

func main() {
	b := &bytes.Buffer{}
	resp.SimpleString("OK").Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.SimpleError("Error message").Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.Integer(123).Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.Integer(-123).Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.BulkString("hello").Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.Array().Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	arr := resp.Array()
	arr.Add(resp.BulkString("hello"))
	arr.Add(resp.BulkString("world"))
	arr.Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	arr = resp.Array()
	arr.Add(resp.Integer(1))
	arr.Add(resp.Integer(2))
	arr.Add(resp.Integer(3))
	arr.Add(resp.Integer(4))
	arr.Add(resp.BulkString("hello"))
	arr.Marshal(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))
}
