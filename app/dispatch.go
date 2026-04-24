package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func DispatchCommand(RESP resp.RESP) ([]byte, error) {
	result := make([]byte, 0)
	values := []string{}

	RESP.ForEach(func(r resp.RESP) bool {
		values = append(values, r.String())
		return true
	})
	cmd := values[0]
	args := values[1:]

	switch cmd {
	case "PING":
		result = resp.AppendString(result, "PONG")
	case "ECHO":
		result = resp.AppendBulkString(result, args[0])
	default:
		return nil, fmt.Errorf("unknown command %s", cmd)
	}

	return result, nil
}
