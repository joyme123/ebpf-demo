/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"

	"github.com/joyme123/ebpf-demo/pkg/utils"
	"github.com/joyme123/ebpf-demo/pkg/xdp_pass_and_drop"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var opt options

type options struct {
	Prog string
	Intf string
}

// passdropCmd represents the passdrop command
var passdropCmd = &cobra.Command{
	Use:   "passdrop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("eBPF pass or drop")

		ctx, cancel := context.WithCancel(context.Background())
		utils.SetupSigHandlers(cancel)

		app, err := xdp_pass_and_drop.NewXdpPassAndDropApp()
		if err != nil {
			log.Panicf("new app error: %s", err)
		}
		err = app.Launch(ctx, opt.Intf, opt.Prog)
		if err != nil {
			log.Panicf("launch app error: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(passdropCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// passdropCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// passdropCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	passdropCmd.Flags().StringVar(&opt.Prog, "prog", "", "which program to load: xdp_pass, xdp_drop, xdp_aborted")
	passdropCmd.Flags().StringVar(&opt.Intf, "intf", "", "the network interface to load xdp program")
	cobra.MarkFlagRequired(passdropCmd.Flags(), "prog")
	cobra.MarkFlagRequired(passdropCmd.Flags(), "intf")
}
