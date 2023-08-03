package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const logTimeFormat = "2006-01-02T15:04:05.999"

func init() {
	// временная метка для всех логгеров и логов должна быть в едином формате
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// DefayltLogger возвращает логгер по умолчанию для проекта.
func DefaultLogger() zerolog.Logger {
	return NewLogger(zerolog.InfoLevel.String())
}

// NewLogger инициализирует и возвращает логгер с установленным уровнем level.
// Значения level см. https://pkg.go.dev/github.com/rs/zerolog#LevelFieldName
func NewLogger(level string) zerolog.Logger {
	loggerLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		loggerLevel = zerolog.InfoLevel
	}

	// подготавливаем writer, который будет писать логи уровня error и выше в STDERR.
	// это полезно для демонов, которые собирают логи для аггрегаторов, чтобы не нужно
	// было настраивать парсеры для определния логов ошибок и обычных логов, а для
	// человеческого глаза различия не будет.
	multi := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logTimeFormat},
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logTimeFormat},
	)

	logger := zerolog.New(multi).Level(loggerLevel).With().Timestamp().Logger()

	return logger
}

type loggerKey struct{}

// FromContext возвращает логгер из контекста.
// Если логгера не было установлено в контексте, то возвращается дефолтный.
func FromContext(ctx context.Context) zerolog.Logger {
	lg, ok := ctx.Value(loggerKey{}).(zerolog.Logger)
	if !ok {
		return DefaultLogger()
	}

	return lg
}

// UpdateContext добавляет логгер в контекст.
func UpdateContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
