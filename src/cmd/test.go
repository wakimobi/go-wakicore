package cmd

import (
	"log"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/providers/telco"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Service CLI",
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
		 * SETUP LOG
		 */
		logger := logger.NewLogger(cfg)

		log.Println(cfg)

		sub := &entity.Subscription{
			Msisdn: "6281299708781",
		}

		service := &entity.Service{
			ID: 1,
		}

		content := &entity.Content{
			Value: "Gabung kompetisi GOALY dan dapatkan reward Exclusive, untuk login klik: https://tsel.goaly.mobi/ PIN: 318477 Stop Ketik UNREG GOALY ke 99790. CS:021-52922391",
			Tid:   "0",
		}

		t := telco.NewTelco(cfg, logger, sub, service, content)

		t.Token()

	},
}
