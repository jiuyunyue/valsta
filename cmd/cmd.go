package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "valsta",
		Short: "validator information collection for cosmos chains",
		Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "create database and table",
		Run: func(cmd *cobra.Command, args []string) {
			Init()
		},
	}
	cleanCmd = &cobra.Command{
		Use:     "cleanDatabase",
		Aliases: []string{"clean"},
		Short:   "drop tables and database",
		Run: func(cmd *cobra.Command, args []string) {
			CleanDatabase()
		},
	}
	queryCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "query all validator info",
		Run: func(cmd *cobra.Command, args []string) {
			infos, err := GetValInfos()
			if err != nil {
				cmd.Println(err.Error())
				_ = cmd.Help()
				return
			}
			for _, info := range infos {
				cmd.Printf(info.String())
			}
			cmd.Printf("total:%d \n", len(infos))

		},
	}
	startCmd = &cobra.Command{
		Use:   "start [start height] [end height]",
		Short: "Start valsta",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			startHeight, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				cmd.Println(err.Error())
				_ = cmd.Help()
				return
			}
			endHeight, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				cmd.Println(err.Error())
				_ = cmd.Help()
				return
			}
			if startHeight > endHeight{
				cmd.Println("startHeight cannot bigger than endHeight")
				return
			}
			if startHeight<=0{
				cmd.Println("startHeight error")
				return
			}
			sta, err := ValSta(startHeight, endHeight)
			if err != nil {
				cmd.Println(err.Error())
				_ = cmd.Help()
			}
			cmd.Printf("insert %d validators success \n", len(sta))
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&GrpcUrl, "grpc", "g", "localhost:9090", "-g <url>")
	startCmd.Flags().StringVarP(&RpcUrl, "rpc", "r", "http://localhost:26657", "-r <url>")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(queryCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
