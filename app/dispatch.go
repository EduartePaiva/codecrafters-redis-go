package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type CacheValue struct {
	Value    string
	ExpireAt *time.Time
}

var cache_mu sync.Mutex
var rpush_mu sync.Mutex
var CACHE map[string]CacheValue = make(map[string]CacheValue)
var RPUSHVALUES map[string][]string = make(map[string][]string)

func DispatchCommand(RESP resp.RESP) ([]byte, error) {
	result := make([]byte, 0)
	values := []string{}

	RESP.ForEach(func(r resp.RESP) bool {
		values = append(values, r.String())
		return true
	})
	cmd := strings.ToUpper(values[0])
	args := values[1:]

	switch cmd {
	case "PING":
		return resp.AppendString(result, "PONG"), nil
	case "ECHO":
		return resp.AppendBulkString(result, args[0]), nil
	case "SET":
		key := args[0]
		value := args[1]
		if len(args) == 2 {
			simpleSet(key, value)
			return resp.AppendOK(result), nil
		}
		setOptions := strings.ToUpper(args[2])
		delay, err := strconv.Atoi(args[3])
		if err != nil {
			return nil, err
		}

		switch setOptions {
		case "PX":
			delayedSet(key, value, time.Duration(delay)*time.Millisecond)
		case "EX":
			delayedSet(key, value, time.Duration(delay)*time.Second)
		}

		return resp.AppendOK(result), nil
	case "GET":
		v, ok := CACHE[args[0]]
		if ok {
			if v.ExpireAt == nil || v.ExpireAt.Compare(time.Now()) >= 0 {
				return resp.AppendBulkString(result, v.Value), nil
			}

			delete(CACHE, args[0])
		}
		return resp.AppendNull(result), nil
	case "RPUSH":
		rpush_mu.Lock()
		list, ok := RPUSHVALUES[args[0]]
		if !ok {
			list = make([]string, 1)
		}
		list = append(list, args[1])
		RPUSHVALUES[args[0]] = list
		rpush_mu.Unlock()
		return resp.AppendInt(result, int64(len(list))), nil
	default:
		return nil, fmt.Errorf("unknown command %s", cmd)
	}
}

func simpleSet(key, value string) {
	cache_mu.Lock()
	CACHE[key] = CacheValue{Value: value}
	cache_mu.Unlock()
}

func delayedSet(key, value string, delay time.Duration) {
	cache_mu.Lock()
	expireAt := time.Now().Add(delay)
	CACHE[key] = CacheValue{Value: value, ExpireAt: &expireAt}
	cache_mu.Unlock()
}
