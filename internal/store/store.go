package store

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/xorium/logging-example/pkg/log"
)

var ErrStoreIsEmpty = errors.New("store has no elements")

func loggerFrom(ctx context.Context) zerolog.Logger {
	return log.FromContext(ctx).With().Str("from", "Memory Store").Logger()
}

type Memory struct {
	mu     *sync.Mutex
	values []string
}

func NewMemory() *Memory {
	return &Memory{
		mu:     new(sync.Mutex),
		values: make([]string, 0),
	}
}

func (m *Memory) Push(ctx context.Context, value string) {
	lg := loggerFrom(ctx)
	lg.Trace().Str("value", value).Msg("Pushing value")

	m.mu.Lock()
	m.values = append(m.values, value)
	m.mu.Unlock()
}

func (m *Memory) Pop(ctx context.Context) (string, error) {
	lg := loggerFrom(ctx)

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.values) == 0 {
		err := errors.Wrap(ErrStoreIsEmpty, "can't pop value from empty memory store")
		lg.Err(err).Msg("can't get some stuf value")

		return "", err
	}

	lastIndex := len(m.values) - 1
	value := m.values[lastIndex]
	m.values = m.values[:lastIndex]

	return value, nil
}
