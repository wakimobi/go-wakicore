package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/wakimobi/go-wakicore/src/datasource/rabbitmq/queue"
)

var consumerMOCmd = &cobra.Command{
	Use:   "consumer-mo",
	Short: "Consumer MO Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		/**
		 * SETUP CHANNEL
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_MOEXCHANGE, true, RMQ_MOQUEUE)

		/**
		 * QUEUE HANDLER
		 */
		messagesData := queue.Rabbit.Subscribe(1, false, RMQ_MOQUEUE, RMQ_MOEXCHANGE, RMQ_MOQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				moProcessor(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerDRCmd = &cobra.Command{
	Use:   "consumer-dr",
	Short: "Consumer DR Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Setup Channel
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_DREXCHANGE, true, RMQ_DRQUEUE)

		/**
		 * Queue Hanlder
		 */
		messagesData := queue.Rabbit.Subscribe(1, false, RMQ_DRQUEUE, RMQ_DREXCHANGE, RMQ_DRQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				drProcessor(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRenewalCmd = &cobra.Command{
	Use:   "consumer-renewal",
	Short: "Consumer Renewal Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Setup Channel
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RENEWALEXCHANGE, true, RMQ_RENEWALQUEUE)

		/**
		 * Queue Handler
		 */
		messagesData := queue.Rabbit.Subscribe(1, false, RMQ_RENEWALQUEUE, RMQ_RENEWALEXCHANGE, RMQ_RENEWALQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				renewalProcessor(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryCmd = &cobra.Command{
	Use:   "consumer-retry",
	Short: "Consumer Retry Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Setup Channel
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYEXCHANGE, true, RMQ_RETRYQUEUE)

		/**
		 * Queue Handler
		 */
		messagesData := queue.Rabbit.Subscribe(1, false, RMQ_RETRYQUEUE, RMQ_RETRYEXCHANGE, RMQ_RETRYQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				retryProcessor(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}
