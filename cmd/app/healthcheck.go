package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var HealthCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check service health",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		resp, err := client.Get("http://localhost:8080/health")
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
