/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"

	"github.com/joyme123/ebpf-demo/pkg/map_counter"
	"github.com/joyme123/ebpf-demo/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var statsOpt statsOptions

type statsOptions struct {
	Intf        string
	PrintPeriod int
}

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("eBPF network stats")

		ctx, cancel := context.WithCancel(context.Background())
		utils.SetupSigHandlers(cancel)

		app, err := map_counter.NewMapCounterApp()
		if err != nil {
			log.Panicf("new app error: %s", err)
		}

		err = app.Launch(ctx, statsOpt.Intf, statsOpt.PrintPeriod)
		if err != nil {
			log.Panicf("launch app error: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	statsCmd.Flags().StringVar(&statsOpt.Intf, "intf", "", "network interface to load xdp program")
	statsCmd.Flags().IntVar(&statsOpt.PrintPeriod, "peroid", 1, "the peroid to print network statistic data")
	cobra.MarkFlagRequired(statsCmd.Flags(), "intf")
}
