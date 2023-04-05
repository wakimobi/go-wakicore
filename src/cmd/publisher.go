package cmd

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/datasource/pgsql/db"
	"github.com/idprm/go-pass-tsel/src/datasource/rabbitmq"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var publisherRenewalCmd = &cobra.Command{
	Use:   "publisher_renewal",
	Short: "Renewal CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * LOAD CONFIG
		 */
		cfg, err := config.LoadSecret("secret.yaml")
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP PGSQL
		 */
		db := db.InitDB(cfg)

		/**
		 * SETUP RMQ
		 */
		queue := rabbitmq.InitQueue(cfg)

		/**
		 * SETUP CHANNEL
		 */
		queue.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RENEWALEXCHANGE, true, RMQ_RENEWALQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked("RENEWAL", timeNow) {

				scheduleService.UpdateSchedule(false, "RENEWAL")

				go func() {
					populateRenewal(db, queue)
				}()
			}

			if scheduleService.GetLocked("RENEWAL", timeNow) {
				scheduleService.UpdateSchedule(true, "RENEWAL")
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherRetryCmd = &cobra.Command{
	Use:   "publisher_retry",
	Short: "Retry CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * LOAD CONFIG
		 */
		cfg, err := config.LoadSecret("secret.yaml")
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP PGSQL
		 */
		db := db.InitDB(cfg)

		/**
		 * SETUP RMQ
		 */
		queue := rabbitmq.InitQueue(cfg)

		/**
		 * SETUP CHANNEL
		 */
		queue.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYEXCHANGE, true, RMQ_RETRYQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked("RETRY", timeNow) {

				scheduleService.UpdateSchedule(false, "RETRY")

				go func() {
					populateRetry(db, queue)
				}()
			}

			if scheduleService.GetLocked("RETRY", timeNow) {
				scheduleService.UpdateSchedule(true, "RETRY")
			}

			time.Sleep(timeDuration * time.Minute)

		}

	},
}

var publisherReminderCmd = &cobra.Command{
	Use:   "publisher_reminder",
	Short: "Reminder CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * LOAD CONFIG
		 */
		cfg, err := config.LoadSecret("secret.yaml")
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP PGSQL
		 */
		db := db.InitDB(cfg)

		/**
		 * SETUP RMQ
		 */
		queue := rabbitmq.InitQueue(cfg)

		/**
		 * SETUP CHANNEL
		 */
		queue.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_REMINDEREXCHANGE, true, RMQ_REMINDERQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked("REMINDER", timeNow) {

				scheduleService.UpdateSchedule(false, "REMINDER")

				go func() {
					populateReminder(db, queue)
				}()
			}

			if scheduleService.GetLocked("REMINDER", timeNow) {
				scheduleService.UpdateSchedule(true, "REMINDER")
			}

			time.Sleep(timeDuration * time.Minute)

		}

	},
}

func populateRenewal(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RenewalSubscription()

	for _, s := range *subs {

		var sub *entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)

		queue.IntegratePublish(RMQ_RENEWALEXCHANGE, RMQ_RENEWALQUEUE, RMQ_DATATYPE, "", string(json))
	}
}

func populateRetry(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetrySubscription()

	for _, s := range *subs {
		var sub *entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)

		queue.IntegratePublish(RMQ_RETRYEXCHANGE, RMQ_RETRYQUEUE, RMQ_DATATYPE, "", string(json))
	}
}

func populateReminder(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.ReminderSubscription()

	for _, s := range *subs {
		var sub *entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)

		queue.IntegratePublish(RMQ_REMINDEREXCHANGE, RMQ_REMINDERQUEUE, RMQ_DATATYPE, "", string(json))
	}
}
