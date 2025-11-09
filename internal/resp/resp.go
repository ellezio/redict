package resp

import (
	"bytes"
	"strconv"
)

const (
	termination string = "\r\n"

	simpleString byte = '+'
	simpleError  byte = '-'
	integer      byte = ':'
	bulkString   byte = '$'
	array        byte = '*'
)

type Marshaler interface {
	Marshal(b *bytes.Buffer)
}

type JustString struct{ value string }

func (l *JustString) Marshal(b *bytes.Buffer) {
	b.WriteString(l.value)
	b.WriteString(termination)
}

type Simple struct {
	prefix byte
	value  string
}

func (s *Simple) Marshal(b *bytes.Buffer) {
	b.WriteByte(s.prefix)
	b.WriteString(s.value)
	b.WriteString(termination)
}

type Aggregate struct {
	prefix byte
	length uint64
	values []Marshaler
}

func (a *Aggregate) Marshal(b *bytes.Buffer) {
	b.WriteByte(a.prefix)
	b.WriteString(strconv.FormatUint(a.length, 10))
	b.WriteString(termination)

	for _, v := range a.values {
		v.Marshal(b)
	}
}

func (a *Aggregate) Add(dataType Marshaler) {
	a.values = append(a.values, dataType)
	a.length++
}

func SimpleString(v string) *Simple {
	return &Simple{simpleString, v}
}

func SimpleError(v string) *Simple {
	return &Simple{simpleError, v}
}

func Integer(v int64) *Simple {
	return &Simple{integer, strconv.FormatInt(v, 10)}
}

func BulkString(v string) *Aggregate {
	s := make([]Marshaler, 1)
	s[0] = &JustString{v}
	return &Aggregate{bulkString, uint64(len(v)), s}
}

func Array() *Aggregate {
	s := make([]Marshaler, 0, 4)
	return &Aggregate{array, 0, s}
}
