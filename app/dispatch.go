package main

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

var mu sync.Mutex
var CACHE map[string]string = make(map[string]string)

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
	case "SET":
		mu.Lock()
		CACHE[args[0]] = args[1]
		mu.Unlock()
		result = resp.AppendOK(result)
	case "GET":
		value, ok := CACHE[args[0]]
		if ok {
			result = resp.AppendBulkString(result, value)
		} else {
			result = resp.AppendNull(result)
		}
	default:
		return nil, fmt.Errorf("unknown command %s", cmd)
	}

	return result, nil
}
