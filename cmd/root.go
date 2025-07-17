/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	ui "syspulse/internal/services/UI"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "syspulse",
	Short: "Monitor and analyze system metrics in real-time",
	Long: `SysPulse is a powerful CLI tool for monitoring and analyzing system metrics in real-time.
It provides insights into CPU, memory, disk, and network usage, helping users to diagnose performance issues,
track resource consumption, and optimize system operations. SysPulse supports customizable output formats,
alerting, and historical data analysis for comprehensive system monitoring.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("help") {
			return
		}

		ui.StartUI()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
