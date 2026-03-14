package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeout time.Duration
	url     string
)

var HealthCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check service health",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("unhealthy: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("healthy")
			return nil
		}
		return fmt.Errorf("unhealthy: status %d", resp.StatusCode)
	},
}
