package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func DispatchCommand(RESP resp.RESP) ([]byte, error) {
	result := make([]byte, 0)
	values := make([]string, 1)

	for r := range RESP.ForEach {
		values = append(values, r.String())
	}
	cmd := values[0]
	args := values[1:]

	switch cmd {
	case "PING":
		resp.AppendString(result, "PONG")
	case "ECHO":
		resp.AppendBulkString(result, args[0])
	default:
		return nil, fmt.Errorf("unknown command %s", cmd)
	}

	return result, nil
}
