package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ellezio/redict/internal/resp"
)

func main() {
	b := &bytes.Buffer{}
	resp.NewSimpleString([]byte("OK")).Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.NewSimpleError([]byte("Error message")).Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.NewInteger(123).Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.NewInteger(-123).Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.NewBulkString([]byte("hello")).Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	resp.NewArray().Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	arr := resp.NewArray()
	arr.Add(resp.NewBulkString([]byte("hello")))
	arr.Add(resp.NewBulkString([]byte("world")))
	arr.Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	b.Reset()
	arr = resp.NewArray()
	arr.Add(resp.NewInteger(1))
	arr.Add(resp.NewInteger(2))
	arr.Add(resp.NewInteger(3))
	arr.Add(resp.NewInteger(4))
	arr.Add(resp.NewBulkString([]byte("hello")))
	arr.Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))

	dt, i, err := resp.Decode([]byte("$0\r\n"))
	if err != nil {
		fmt.Println(i, err)
		return
	}
	b.Reset()
	dt.Encode(b)
	fmt.Println(strings.ReplaceAll(b.String(), "\r\n", `\r\n`))
}
