package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wakimobi/go-wakicore/src/app"
	"github.com/wakimobi/go-wakicore/src/datasource/rabbitmq/queue"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		/**
		 * SETUP CHANNEL
		 */
		queue.Rabbit.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_MOEXCHANGE, true, RMQ_MOQUEUE)

		port := os.Getenv("H3I_KB_PORT")

		app := app.StartApplication()
		app.Run(":" + port)

	},
}
