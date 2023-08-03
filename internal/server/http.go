package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/xorium/logging-example/internal/stuf"
	"github.com/xorium/logging-example/pkg/log"
)

// ServerOption функция опции, которая принимает в качестве аргумента сервер HTTP.
type ServerOption func(srv *HTTP)

// OptLogger возвращает функцию опции, которая устанавливает логгер для сервера HTTP.
func OptLogger(lg zerolog.Logger) ServerOption {
	return func(srv *HTTP) {
		srv.WithLogger(lg)
	}
}

func OptStuf(stuf *stuf.Stuf) ServerOption {
	return func(srv *HTTP) {
		srv.stuf = stuf
	}
}

// HTTP структура сервера, содержащая адрес прослушивания, движок Echo и логгер.
type HTTP struct {
	listenAddr string
	engine     *echo.Echo
	stuf       *stuf.Stuf
	log        zerolog.Logger
}

// NewHTTP создает новый экземпляр сервера HTTP с заданными адресом прослушивания и опциями.
func NewHTTP(listenAddr string, opts ...ServerOption) *HTTP {
	s := &HTTP{
		listenAddr: listenAddr,
		engine:     echo.New(),
		log:        zerolog.Nop(),
	}

	// Применение переданных опций
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithLogger устанавливает логгер для сервера и добавляет промежуточный middleware для логирования.
func (s *HTTP) WithLogger(lg zerolog.Logger) *HTTP {
	s.log = lg.With().Str("from", "HTTP server").Logger()

	// middleware для логирования запросов
	s.engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// получаем или генерирует ID запроса
			reqID := req.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
				req.Header.Set("X-Request-ID", reqID)
			}

			// не забываем добавить X-Request-ID в заголовки
			res.Header().Set("X-Request-ID", reqID)
			// создаем новый логгер с дополненным контекстом (с новыми постоянными полями)
			lg := s.log.With().Str("request_id", reqID).Logger()

			// добавляем request ID и логгер в контекст запроса (для следующих компонентов приложения)
			ctx := log.UpdateContext(req.Context(), lg)
			ctx = context.WithValue(req.Context(), "request_id", reqID)
			// обновляет структуру HTTP запроса с обновленным контекстом
			c.SetRequest(req.WithContext(ctx))

			lg.Info().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Msg("HTTP request")

			return next(c)
		}
	})

	return s
}

// ping обрабатывает HTTP GET-запросы на /ping и возвращает "pong".
func (s *HTTP) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (s *HTTP) doSomething(c echo.Context) error {
	ctx := c.Request().Context()
	value := c.QueryParam("value")

	if value == "" {
		// Note: здесь логировать не нужно - обычно бесмысленно, но не запрещается
		return c.String(http.StatusBadRequest, "no value")
	}
	// создавать саблоггеры вначале функций - хорошая практика
	lg := log.FromContext(ctx).With().Str("value", value).Logger()

	lg.Info().Str("t", "123").Str("t", "456")

	err := s.stuf.DoSomethingWith(c.Request().Context(), c.QueryParam("value"))
	if err != nil {
		// Note: здесь логировать - хорошая практика: это конечная точка нешго control-flow, будет
		// виден весь контекст входящих данных, при которых эта ошибка на каком-то из
		// внутренних слоев проекта возникла
		lg.Err(err).Msg("stuf has failed to do something")

		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "ok")
}

func (s *HTTP) getSomeValue(c echo.Context) error {
	ctx := c.Request().Context()
	lg := log.FromContext(ctx)
	// Полезно бывает залогировать начало обработки запроса, чтобы потом можно было посмотреть
	// задержки выполнения вызовов различных функций внутри обработчика запроса
	lg.Trace().Msg("Trying to get some value from stuf")

	value, err := s.stuf.GetSomeValue(ctx)
	if err != nil {
		if err != nil {
			// Note: несмотря на то, что у нас нет каких-либо параметров для контекста лога,
			// все равно хорошо бы залогировать ошибку, т.к. это конечная точка нешего control-flow
			lg.Err(err).Msg("stuf can't return some value")

			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	lg.Debug().Str("value", value).Msg("Got some value from stuf")

	return c.String(http.StatusOK, value)
}

// ListenAndServe начинает прослушивание на указанном адресе и служит HTTP-запросам.
func (s *HTTP) ListenAndServe() error {
	s.engine.GET("/ping", s.ping)
	s.engine.POST("/do-something", s.doSomething)
	s.engine.GET("/get-some-value", s.getSomeValue)

	return s.engine.Start(s.listenAddr)
}
