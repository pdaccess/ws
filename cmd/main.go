package main

import (
	"fmt"

	"github.com/pdaccess/ws/cmd/app"
	"github.com/spf13/cobra"
)

var Commit, BuildTime, BuildEnv string

var AppDescription = fmt.Sprintf(`WebService is a server application for managing inventory and user groups.
It provides APIs for creating, updating, deleting and listing inventory items and user groups.

Version: %s
Build Time: %s
Build Environment: %s`, Commit, BuildTime, BuildEnv)

var rootCmd = &cobra.Command{
	Use:   "ws",
	Short: "Main WebService application for managing inventory and user groups",
	Long:  AppDescription,
}

func main() {
	rootCmd.AddCommand(app.ServerCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("app root: %w", err))
	}
}
