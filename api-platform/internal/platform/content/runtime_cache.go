package content

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// runtimeCacheClient only touches this application's Redis namespace. It must
// never issue FLUSHDB/FLUSHALL because Redis is also used by IM and other apps.
type runtimeCacheClient struct {
	addr     string
	password string
	db       int
}

func (c runtimeCacheClient) deletePattern(ctx context.Context, pattern string) (int64, error) {
	if strings.TrimSpace(c.addr) == "" {
		return 0, errors.New("redis address is empty")
	}
	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	if err := redisAuthenticate(reader, conn, c.password, c.db); err != nil {
		return 0, err
	}

	var deleted int64
	cursor := "0"
	for {
		result, err := redisCommand(ctx, reader, conn, "SCAN", cursor, "MATCH", pattern, "COUNT", "500")
		if err != nil {
			return 0, err
		}
		values, ok := result.([]any)
		if !ok || len(values) != 2 {
			return 0, errors.New("unexpected redis scan response")
		}
		cursor, ok = values[0].(string)
		if !ok {
			return 0, errors.New("unexpected redis cursor")
		}
		keys, ok := values[1].([]any)
		if !ok {
			return 0, errors.New("unexpected redis key list")
		}
		if len(keys) > 0 {
			args := make([]string, 0, len(keys)+1)
			args = append(args, "UNLINK")
			for _, item := range keys {
				key, ok := item.(string)
				if !ok {
					return 0, errors.New("unexpected redis key")
				}
				args = append(args, key)
			}
			removed, err := redisCommand(ctx, reader, conn, args...)
			if err != nil {
				return 0, err
			}
			count, ok := removed.(int64)
			if !ok {
				return 0, errors.New("unexpected redis unlink response")
			}
			deleted += count
		}
		if cursor == "0" {
			return deleted, nil
		}
	}
}

func redisAuthenticate(reader *bufio.Reader, conn net.Conn, password string, db int) error {
	if password != "" {
		if _, err := redisCommand(context.Background(), reader, conn, "AUTH", password); err != nil {
			return err
		}
	}
	if db > 0 {
		if _, err := redisCommand(context.Background(), reader, conn, "SELECT", strconv.Itoa(db)); err != nil {
			return err
		}
	}
	return nil
}

func redisCommand(ctx context.Context, reader *bufio.Reader, conn net.Conn, args ...string) (any, error) {
	deadline := time.Now().Add(3 * time.Second)
	if at, ok := ctx.Deadline(); ok && at.Before(deadline) {
		deadline = at
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(conn, "*%d\r\n", len(args)); err != nil {
		return nil, err
	}
	for _, arg := range args {
		if _, err := fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return nil, err
		}
	}
	return readRESP(reader)
}

func readRESP(reader *bufio.Reader) (any, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	line, err := readRESPLine(reader)
	if err != nil {
		return nil, err
	}
	switch prefix {
	case '+':
		return line, nil
	case '-':
		return nil, errors.New("redis: " + line)
	case ':':
		return strconv.ParseInt(line, 10, 64)
	case '$':
		length, err := strconv.Atoi(line)
		if err != nil || length < -1 {
			return nil, errors.New("invalid redis bulk length")
		}
		if length == -1 {
			return "", nil
		}
		data := make([]byte, length+2)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		return string(data[:length]), nil
	case '*':
		length, err := strconv.Atoi(line)
		if err != nil || length < -1 {
			return nil, errors.New("invalid redis array length")
		}
		if length == -1 {
			return []any{}, nil
		}
		items := make([]any, length)
		for index := range items {
			if items[index], err = readRESP(reader); err != nil {
				return nil, err
			}
		}
		return items, nil
	default:
		return nil, errors.New("invalid redis response")
	}
}

func readRESPLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", errors.New("invalid redis line")
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}
