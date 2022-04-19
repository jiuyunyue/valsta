package cmd

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := Init()
			if err != nil {
				return err
			}
			return nil
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
		Short:   "query commands",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	queryValCmd = &cobra.Command{
		Use:     "val",
		Aliases: []string{"q"},
		Short:   "query all validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			infos, err := GetValInfos()
			if err != nil {
				return err
			}
			for _, info := range infos {
				cmd.Printf(info.String())
			}
			cmd.Printf("total:%d \n", len(infos))
			return nil
		},
	}
	queryVoterCmd = &cobra.Command{
		Use:     "voters",
		Aliases: []string{"q"},
		Short:   "query all validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			voters, err := GetVoterInfos()
			if err != nil {
				return err
			}
			for k, v := range voters {
				cmd.Printf("voter %v, %s", k, v)
			}
			cmd.Printf("total:%d \n", len(voters))
			return nil
		},
	}
	startCmd = &cobra.Command{
		Use:   "start [start height] [end height]",
		Short: "Start valsta",
		Args:  cobra.RangeArgs(2, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			startHeight, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			endHeight, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			if startHeight > endHeight {
				cmd.Println("startHeight cannot bigger than endHeight")
				return errors.New("startHeight cannot bigger than endHeight")
			}
			if startHeight <= 0 {
				return errors.New("startHeight error")
			}
			sta, err := ValSta(startHeight, endHeight)
			if err != nil {
				return err
			}
			cmd.Printf("insert %d validators success \n", len(sta))
			return nil
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&GrpcUrl, "grpc", "g", "localhost:9090", "-g <url>")
	startCmd.Flags().StringVarP(&RpcUrl, "rpc", "r", "http://localhost:26657", "-r <url>")
	queryVoterCmd.Flags().StringVarP(&GrpcUrl, "grpc", "g", "localhost:9090", "-g <url>")
	queryVoterCmd.Flags().StringVarP(&RpcUrl, "rpc", "r", "http://localhost:26657", "-r <url>")

	queryCmd.AddCommand(queryValCmd)
	queryCmd.AddCommand(queryVoterCmd)
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
