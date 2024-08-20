package tarantool

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"time"
)

type Storage struct {
	connection  *tarantool.Connection
	spaceName   string
	primaryName string
}

func NewTarantoolStorage(
	storagePath string,
	spaceName string,
) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address: storagePath,
		User:    "guest",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return nil, fmt.Errorf("%w", err)
	}

	return &Storage{
		connection: conn,
		spaceName:  spaceName,
	}, nil
}

func (s *Storage) Close() error {
	return s.connection.Close()
}

func (s *Storage) Data(keys []string) (map[string]any, error) {
	const op = "storage.tarantool.Data"

	result := make(map[string]any)

	for _, key := range keys {
		resp, err := s.connection.Do(
			tarantool.NewSelectRequest(s.spaceName).
				Limit(1).
				Iterator(tarantool.IterEq).
				Key([]any{key}),
		).Get()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if len(resp) > 0 {
			result[key] = resp[0].([]any)[1]
		}
	}

	return result, nil
}

func (s *Storage) SaveData(data map[string]any) error {
	const op = "storage.tarantool.SaveData"

	var futures []*tarantool.Future
	for key, value := range data {
		request := tarantool.NewInsertRequest(s.spaceName).Tuple([]any{key, value})
		futures = append(futures, s.connection.Do(request))
	}

	for _, future := range futures {
		_, err := future.Get()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
