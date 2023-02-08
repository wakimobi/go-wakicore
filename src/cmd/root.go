package cmd

import (
	"github.com/spf13/cobra"
)

const (
	RMQ_AUTOACK            = false
	RMQ_EXCHANGETYPE       = "direct"
	RMQ_EXCHANGEDURABILITY = false
	RMQ_QUEUEDURABILITY    = true
	RMQ_DATATYPE           = "application/json"
	RMQ_MOEXCHANGE         = "E_MO_H3I_KB"
	RMQ_MOQUEUE            = "Q_MO_H3I_KB"
	RMQ_DREXCHANGE         = "E_DR_H3I_KB"
	RMQ_DRQUEUE            = "Q_DR_H3I_KB"
	RMQ_RENEWALEXCHANGE    = "E_RENEWAL_H3I_KB"
	RMQ_RENEWALQUEUE       = "Q_RENEWAL_H3I_KB"
	RMQ_RETRYEXCHANGE      = "E_RETRY_H3I_KB"
	RMQ_RETRYQUEUE         = "Q_RETRY_H3I_KB"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long:  `Cobra is a CLI library for Go that empowers applications.`,
	}
)

func init() {

	/**
	 * Server service
	 */
	rootCmd.AddCommand(serverCmd)

	/**
	 * Consumer service
	 */
	rootCmd.AddCommand(consumerMOCmd)
	rootCmd.AddCommand(consumerDRCmd)
	rootCmd.AddCommand(consumerRenewalCmd)
	rootCmd.AddCommand(consumerRetryCmd)

	/**
	 * Publisher service
	 */
	rootCmd.AddCommand(publisherRenewalCmd)
	rootCmd.AddCommand(publisherRetryCmd)

	/**
	 * Test service
	 */
	rootCmd.AddCommand(hitMtCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
