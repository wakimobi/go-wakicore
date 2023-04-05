package cmd

import "github.com/spf13/cobra"

const (
	RMQ_AUTOACK            = false
	RMQ_EXCHANGETYPE       = "direct"
	RMQ_EXCHANGEDURABILITY = false
	RMQ_QUEUEDURABILITY    = true
	RMQ_DATATYPE           = "application/json"
	RMQ_MOEXCHANGE         = "E_MO"
	RMQ_MOQUEUE            = "Q_MO"
	RMQ_DREXCHANGE         = "E_DR"
	RMQ_DRQUEUE            = "Q_DR"
	RMQ_REMINDEREXCHANGE   = "E_REMINDER"
	RMQ_REMINDERQUEUE      = "Q_REMINDER"
	RMQ_RENEWALEXCHANGE    = "E_RENEWAL"
	RMQ_RENEWALQUEUE       = "Q_RENEWAL"
	RMQ_RETRYEXCHANGE      = "E_RETRY"
	RMQ_RETRYQUEUE         = "Q_RETRY"
	MT_FIRSTPUSH           = "MT_FIRSTPUSH"
	ACT_RENEWAL            = "RENEWAL"
	ACT_RETRY              = "RETRY"
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
	 * WEBSERVER SERVICE
	 */
	rootCmd.AddCommand(serverCmd)

	/**
	 * RABBITMQ SERVICE
	 */
	rootCmd.AddCommand(consumerMOCmd)
	rootCmd.AddCommand(consumerDRCmd)
	rootCmd.AddCommand(consumerRenewalCmd)
	rootCmd.AddCommand(consumerRetryCmd)

	rootCmd.AddCommand(publisherReminderCmd)
	rootCmd.AddCommand(publisherRenewalCmd)
	rootCmd.AddCommand(publisherRetryCmd)

	rootCmd.AddCommand(testCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
