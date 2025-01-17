package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
)

// healthcheckCommand ヘルスチェックコマンド
func healthcheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "healthcheck",
		Short: "Run healthcheck",
		Run: func(cmd *cobra.Command, args []string) {
			logger := getCLILogger()
			resp, err := http.DefaultClient.Get(fmt.Sprintf("http://localhost:%d/api/ping", c.Port))
			if err != nil {
				logger.Fatal("HTTP Client Error", zap.Error(err))
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				logger.Fatal("Unexpected status", zap.Int("status", resp.StatusCode))
			}
		},
	}
}
