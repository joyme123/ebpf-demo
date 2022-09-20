/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyme123/cobra-cli/pkg/xconnect"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// xconnectCmd represents the xconnect command
var xconnectCmd = &cobra.Command{
	Use:   "xconnect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("eBPF x-connect")

		ctx, cancel := context.WithCancel(context.Background())
		setupSigHandlers(cancel)

		conf, err := cmd.Flags().GetString("conf")
		if err != nil {
			log.Errorf("flag conf is not provided, err: %s\n", err)
			return
		}

		cfg, err := newFromFile(conf)
		if err != nil {
			log.Errorf("parse configurations from file failed, err: %s\n", err)
			return
		}

		updateCh := make(chan map[string]string, 1)
		go configWatcher(conf, updateCh)

		app, err := xconnect.NewXconnectApp(cfg.Links)
		if err != nil {
			log.Errorf("Loading eBPF: %s", err)
			return
		}

		app.Launch(ctx, updateCh)

		return
	},
}

func init() {
	rootCmd.AddCommand(xconnectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xconnectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	xconnectCmd.Flags().StringP("conf", "c", "", "configuration file")
	cobra.MarkFlagRequired(xconnectCmd.Flags(), "conf")
}

func setupSigHandlers(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Infof("Received syscall: %+v", sig)
		cancel()
	}()
}
