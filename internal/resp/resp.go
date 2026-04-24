package resp

import (
	"strconv"
	"strings"
)

type Type byte

// Various RESP kinds
const (
	Integer Type = ':'
	String  Type = '+'
	Bulk    Type = '$'
	Array   Type = '*'
	Error   Type = '-'
)

type RESP struct {
	Type  Type
	Data  []byte
	Count int
}

func (r RESP) String() string {
	return string(r.Data)
}
func (r RESP) Int() int64 {
	v, _ := strconv.ParseInt(r.String(), 10, 64)
	return v
}

func (*RESP) Parse(data string) error {
	lines := make([]string, 1)
	for line := range strings.Lines(data) {
		lines = append(lines, line)
	}

	dataType := lines[0]
	numItems, err := strconv.Atoi(string(data[1:]))
	if err != nil {
		return err
	}

}
