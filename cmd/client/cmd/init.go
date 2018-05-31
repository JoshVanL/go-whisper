package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/joshvanl/go-whisper/pkg/client"
)

const FlagLogLevel = "log-level"
const FlagServerAddr = "server-address"
const FlagConfigDir = "config"

var RootCmd = &cobra.Command{
	Use:   "client",
	Short: "An end to end encrypted messaging app written in Go.",
	Run: func(cmd *cobra.Command, args []string) {

		log := LogLevel(cmd)

		addr, err := cmd.PersistentFlags().GetString(FlagServerAddr)
		if err != nil {
			log.Fatalf("failed to resolve server address: %v", err)
		}

		dir, err := cmd.PersistentFlags().GetString(FlagConfigDir)
		if err != nil {
			log.Fatalf("failed to resolve configAdirectory flag: %v", err)
		}

		if dir == "." {
			dir, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get working directory: %v", err)
			}
		} else {
			dir, err = homedir.Expand(dir)
			if err != nil {
				log.Fatalf("failed to expand go-whipser config directory: %v", err)
			}
		}

		_, err = client.New(addr, dir, log)
		if err != nil {
			log.Fatalf("error creating client: %v", err)
		}

		stopCh := make(chan struct{})
		<-stopCh

		//if err := c.Connect(); err != nil {
		//	log.Fatalf("error running client: %v", err)
		//}
	},
}

func init() {
	RootCmd.PersistentFlags().IntP(FlagLogLevel, "l", 1, "Set the log level of output. 0-Fatal 1-Info 2-Debug")
	RootCmd.PersistentFlags().StringP(FlagServerAddr, "s", "127.0.0.1:6667", "Set the address of the server")
	RootCmd.PersistentFlags().StringP(FlagConfigDir, "c", "~/.go-whisper", "Directory of go-whipser directory")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func LogLevel(cmd *cobra.Command) *logrus.Entry {
	logger := logrus.New()

	i, err := cmd.PersistentFlags().GetInt(FlagLogLevel)
	if err != nil {
		logrus.Fatalf("failed to get log level of flag: %s", err)
	}
	if i < 0 || i > 2 {
		logrus.Fatalf("not a valid log level")
	}
	switch i {
	case 0:
		logger.Level = logrus.FatalLevel
	case 1:
		logger.Level = logrus.InfoLevel
	case 2:
		logger.Level = logrus.DebugLevel
	}

	return logrus.NewEntry(logger)
}
