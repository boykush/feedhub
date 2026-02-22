package main

import (
	"log"
	"os"
	"syscall"

	"github.com/boykush/feedhub/server/feed/cmd/server/provider"
	"github.com/caarlos0/env/v11"
	"github.com/samber/do/v2"
)

func main() {
	cfg, err := env.ParseAs[provider.Config]()
	if err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
	}

	injector := do.New()
	do.ProvideValue(injector, cfg)
	provider.Register(injector)

	do.MustInvoke[*provider.GRPCServer](injector)

	injector.ShutdownOnSignals(syscall.SIGTERM, os.Interrupt)
}
