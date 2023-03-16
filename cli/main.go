package main

import (
	"github.com/bhoriuchi/embedded-nats-jetstream/internal/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.
		With().
		Caller().
		Logger().
		Level(zerolog.DebugLevel)

	cmd.NewRootCommand().Execute()
}
