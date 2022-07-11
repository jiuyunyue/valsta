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
	queryCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "query commands",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	querySignTimesCmd = &cobra.Command{
		Use:   "signTimes [address]",
		Args:  cobra.MinimumNArgs(1),
		Short: "query sign times at 841500-1412246 ",
		RunE: func(cmd *cobra.Command, args []string) error {
			height, err := SignTimes(args[0])
			if err != nil {
				return err
			}
			if len(height) != 0 {
				cmd.Printf("total:%d \n start at %d \n", len(height), height[0])
			} else {
				cmd.Printf("can not find your address at height 841500-1412246")
			}
			return nil
		},
	}
	querySignHeightCmd = &cobra.Command{
		Use:   "signHeight [address]",
		Args:  cobra.MinimumNArgs(1),
		Short: "query sign height at 841500-1412246 ",
		RunE: func(cmd *cobra.Command, args []string) error {
			have, height, err := SignHeight(args[0])
			if err != nil {
				return err
			}

			if have {
				cmd.Printf("start at %d \n", height)
			} else {
				cmd.Printf("can not find your address at height 841500-1412246")
			}
			return nil
		},
	}
	queryVoterCmd = &cobra.Command{
		Use:     "voters",
		Aliases: []string{"q"},
		Short:   "query all voter info",
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
	querySignHeightCmd.Flags().StringVarP(&GrpcUrl, "grpc", "g", "localhost:9090", "-g <url>")
	querySignHeightCmd.Flags().StringVarP(&RpcUrl, "rpc", "r", "http://localhost:26657", "-r <url>")
	querySignTimesCmd.Flags().StringVarP(&GrpcUrl, "grpc", "g", "localhost:9090", "-g <url>")
	querySignTimesCmd.Flags().StringVarP(&RpcUrl, "rpc", "r", "http://localhost:26657", "-r <url>")

	queryCmd.AddCommand(querySignTimesCmd)
	queryCmd.AddCommand(querySignHeightCmd)
	queryCmd.AddCommand(queryVoterCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(queryCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
