package stuf

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xorium/logging-example/internal/store"
	"github.com/xorium/logging-example/pkg/log"
)

func loggerFrom(ctx context.Context) zerolog.Logger {
	return log.FromContext(ctx).With().Str("from", "Stuf").Logger()
}

type Stuf struct {
	store *store.Memory
	log   zerolog.Logger
}

func NewStuf(store *store.Memory) *Stuf {
	return &Stuf{
		store: store,
		log:   zerolog.Nop(),
	}
}

func (s *Stuf) DoSomethingWith(ctx context.Context, value string) error {
	lg := loggerFrom(ctx)

	lg.Trace().Str("value", value).Msg("Start doing something with value")

	s.store.Push(ctx, value)

	lg.Trace().Str("value", value).Msg("Do something with value has been finished")

	return nil
}

func (s *Stuf) GetSomeValue(ctx context.Context) (value string, err error) {
	lg := loggerFrom(ctx)

	lg.Trace().Msg("Start getting some value")

	val, err := s.store.Pop(ctx)
	if err != nil {
		// Внимание! Логировать ошибку (как в строки ниже) обычно НЕ нужно - избыточность
		// lg.Err(err).Msg("can't get some stuf value")
		return "", err
	}

	lg.Debug().Str("value", val).Msg("Getting some value has been finisheed")

	return val, nil
}
