package resp

import (
	"bytes"
	"errors"
	"io"
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

func Decode(b []byte) (DataType, int, error) {
	var dt DataType

	l := len(b)
	if l < 2 {
		return nil, 0, io.EOF
	}

	prefix := b[0]
	i := 1

	switch prefix {
	case simpleString:
		var buf bytes.Buffer
		for ; i < l; i++ {
			if b[i] == '\r' {
				break
			}

			buf.WriteByte(b[i])
		}

		if i+1 >= l || b[i+1] != '\n' {
			return nil, i, errors.New("invalid bytes")
		}
		i += 2

		dt = &SimpleString{buf.Bytes()}

	case simpleError:
		var buf bytes.Buffer
		for ; i < l; i++ {
			if b[i] == '\r' {
				break
			}
			buf.WriteByte(b[i])
		}

		if i+1 >= l || b[i+1] != '\n' {
			return nil, i, errors.New("invalid bytes")
		}
		i += 2

		dt = &SimpleError{buf.Bytes()}

	case integer:
		var buf bytes.Buffer
		for ; i < l; i++ {
			if b[i] == '\r' {
				break
			}
			buf.WriteByte(b[i])
		}

		if i+1 >= l || b[i+1] != '\n' {
			return nil, i, errors.New("invalid bytes")
		}
		i += 2

		v, err := strconv.ParseInt(buf.String(), 10, 64)
		if err != nil {
			return nil, i, err
		}
		dt = &Integer{v}

	case bulkString:
		length := 0
		for _, char := range b[i:] {
			if char == '\r' {
				break
			}

			length = (length * 10) + int(char-'0')
			i++
		}

		if i+1 >= l || b[i+1] != '\n' {
			return nil, i, errors.New("unexpected end of file")
		}
		i += 2
		if i > l {
			return nil, i, errors.New("unexpected end of file")
		}

		var buf bytes.Buffer
		if length > 0 {
			for j := 0; j < length; j++ {
				buf.WriteByte(b[i])
				i++
			}

			if i+1 > l || b[i] != '\r' || b[i+1] != '\n' {
				return nil, i, errors.New("unexpected end of file")
			}
			i += 2

		}

		dt = &BulkString{buf.Bytes()}
	case array:
		length := 0
		for _, char := range b[i:] {
			if char == '\r' {
				break
			}

			length = (length * 10) + int(char-'0')
			i++
		}

		if i+1 >= l || b[i+1] != '\n' {
			return nil, i, errors.New("unexpected end of file")
		}
		i += 2
		if i > l {
			return nil, i, errors.New("unexpected end of file")
		}

		arr := NewArray()
		for j := 0; j < length; j++ {
			elem, n, err := Decode(b[i:])
			i += n
			if err != nil {
				return nil, i, err
			}
			arr.Add(elem)
		}
		dt = arr
	}

	return dt, i, nil
}

type DataType interface {
	Encode(*bytes.Buffer)
}

func NewSimpleString(v string) *SimpleString {
	return &SimpleString{[]byte(v)}
}

type SimpleString struct {
	Value []byte
}

func (ss *SimpleString) Encode(b *bytes.Buffer) {
	b.WriteByte(simpleString)
	b.Write(ss.Value)
	b.WriteString(termination)
}

func NewSimpleError(v string) *SimpleError {
	return &SimpleError{[]byte(v)}
}

type SimpleError struct {
	Value []byte
}

func (se *SimpleError) Encode(b *bytes.Buffer) {
	b.WriteByte(simpleError)
	b.Write(se.Value)
	b.WriteString(termination)
}

func NewInteger(v int64) *Integer {
	return &Integer{v}
}

type Integer struct {
	Value int64
}

func (i *Integer) Encode(b *bytes.Buffer) {
	b.WriteByte(integer)
	b.WriteString(strconv.FormatInt(i.Value, 10))
	b.WriteString(termination)
}

func NewBulkString(v string) *BulkString {
	return &BulkString{[]byte(v)}
}

type BulkString struct {
	Value []byte
}

func (bs *BulkString) Encode(b *bytes.Buffer) {
	b.WriteByte(bulkString)
	b.WriteString(strconv.FormatUint(uint64(len(bs.Value)), 10))
	b.WriteString(termination)
	if len(bs.Value) > 0 {
		b.Write(bs.Value)
		b.WriteString(termination)
	}
}

func NewArray() *Array {
	s := make([]DataType, 0, 4)
	return &Array{s}
}

type Array struct {
	Values []DataType
}

func (a *Array) Encode(b *bytes.Buffer) {
	b.WriteByte(array)
	b.WriteString(strconv.FormatUint(uint64(len(a.Values)), 10))
	b.WriteString(termination)

	for _, v := range a.Values {
		v.Encode(b)
	}
}

func (a *Array) Add(dataType DataType) {
	a.Values = append(a.Values, dataType)
}
