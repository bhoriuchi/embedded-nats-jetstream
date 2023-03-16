package cmd

import (
	"context"
	"strings"

	"github.com/bhoriuchi/embedded-nats-jetstream/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func NewRootCommand() *cobra.Command {
	cobra.OnInitialize(initServerConfig)
	cfg := &server.Config{}

	root := &cobra.Command{
		Use:          "enats",
		Version:      "0.0.1-beta1",
		Short:        "Example embedded NATS jetstream",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(cfg); err != nil {
				return err
			}

			srv, err := server.NewServer(cfg)
			if err != nil {
				return err
			}

			return srv.Start(context.Background())
		},
	}
	// define config file flag
	root.Flags().StringVarP(&configFile, "config", "c", "", "Embedded NATS configuration file")

	// additional flags
	root.PersistentFlags().StringSliceVar(&cfg.Routes, "routes", []string{}, "Routes")
	viper.BindPFlag("routes", root.PersistentFlags().Lookup("routes"))

	return root
}

// initializes the server config
func initServerConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".enatsrc")
	}

	viper.SetEnvPrefix("ENATS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
