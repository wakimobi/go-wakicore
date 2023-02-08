package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wakimobi/go-wakicore/src/datasource/rabbitmq/queue"
)

var publisherRenewalCmd = &cobra.Command{
	Use:   "publisher-renewal",
	Short: "Renewal CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Setup Channel
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RENEWALEXCHANGE, true, RMQ_RENEWALQUEUE)
	},
}

var publisherRetryCmd = &cobra.Command{
	Use:   "publisher-retry",
	Short: "Retry CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Setup Channel
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYEXCHANGE, true, RMQ_RETRYQUEUE)
	},
}

var hitMtCmd = &cobra.Command{
	Use:   "hit-mt",
	Short: "Hit MT CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// sub := subscriptions.Subscription{
		// 	ProductID: 1,
		// 	Msisdn:    "081299708787",
		// }

	},
}
