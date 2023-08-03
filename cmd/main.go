package main

import (
	"github.com/xorium/logging-example/internal/config"
	"github.com/xorium/logging-example/internal/server"
	"github.com/xorium/logging-example/internal/store"
	"github.com/xorium/logging-example/internal/stuf"
	"github.com/xorium/logging-example/pkg/log"
)

func main() {
	cfg, err := config.ParseFromEnv()
	if err != nil {
		panic(err)
	}

	lg := log.NewLogger(cfg.LogLevel)

	store := store.NewMemory()
	stuf := stuf.NewStuf(store)
	httpSrv := server.NewHTTP(cfg.HttpListenAddr, stuf, server.OptLogger(lg))

	panic(httpSrv.ListenAndServe())
}
