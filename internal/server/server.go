package server

import (
	"context"
	"fmt"
	"net/url"
	"os"

	// "github.com/bhoriuchi/embedded-nats-jetstream/pkg/addr"
	"github.com/bhoriuchi/embedded-nats-jetstream/pkg/logging"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
)

const (
	DefaultNATSServerPort  = 4222
	DefaultNATSClusterPort = 6222
	DefaultNATSWebPort     = 8222
	DefaultNATSBindAddress = "0.0.0.0"
	DefaultNATSClusterName = "enats"
)

// Server a server object
type Server struct {
	config   *Config
	hostname string
	wd       string
	nats     *natsserver.Server
	routes   []*url.URL
}

// NewServer creates a new server
func NewServer(config *Config) (*Server, error) {
	// host info
	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msgf("failed to get hostname")
		return nil, err
	}

	llog := log.With().Str("host", hostname).Logger()

	wd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msgf("failed to get working directory")
		return nil, err
	}

	if len(config.Routes) < 1 {
		log.Error().Msgf("no routes provided")
		return nil, fmt.Errorf("no routes provided")
	}

	// parse routes
	routes, err := parseNatsURLs(config.Routes)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse routes")
		return nil, err
	}

	// get bind address
	/*
		saddr, err := addr.Extract("")
		if err != nil {
			saddr = DefaultNATSBindAddress
		}
	*/

	// create jetstream server
	srv, err := natsserver.NewServer(&natsserver.Options{
		Host:            DefaultNATSBindAddress,
		Port:            DefaultNATSServerPort,
		Logtime:         true,
		Trace:           true,
		Debug:           true,
		JetStream:       true,
		ServerName:      hostname,
		StoreDir:        wd,
		Routes:          routes,
		NoSystemAccount: true,
		Username:        "enats",
		Password:        "enats",
		SystemAccount:   natsserver.DEFAULT_SYSTEM_ACCOUNT,
		Accounts: []*natsserver.Account{
			natsserver.NewAccount(natsserver.DEFAULT_SYSTEM_ACCOUNT),
		},
		Cluster: natsserver.ClusterOpts{
			Advertise:      "",
			Host:           DefaultNATSBindAddress,
			Name:           DefaultNATSClusterName,
			Port:           DefaultNATSClusterPort,
			ConnectRetries: 25,
		},
		Users: []*natsserver.User{
			{
				Username: "enats",
				Password: "enats",
				Account:  natsserver.NewAccount(natsserver.DEFAULT_SYSTEM_ACCOUNT),
			},
		},
	})

	if err != nil {
		log.Error().Err(err).Msgf("failed to create new nats server")
		return nil, err
	}

	natsLog := logging.NewNATSLogger(llog)
	srv.SetLogger(natsLog, true, true)

	s := &Server{
		config:   config,
		hostname: hostname,
		wd:       wd,
		routes:   routes,
		nats:     srv,
	}

	return s, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	// Start things up. Block here until done.
	if err := natsserver.Run(s.nats); err != nil {
		natsserver.PrintAndDie(err.Error())
	}

	// Adjust MAXPROCS if running under linux/cgroups quotas.
	undo, err := maxprocs.Set(maxprocs.Logger(s.nats.Debugf))
	if err != nil {
		s.nats.Warnf("Failed to set GOMAXPROCS: %v", err)
	} else {
		defer undo()
		// Reset these from the snapshots from init for monitor.go
		natsserver.SnapshotMonitorInfo()
	}

	s.nats.WaitForShutdown()

	return nil
}
