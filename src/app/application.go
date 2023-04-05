package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/wiliehidayat87/rmqp"
)

func StartApplication(cfg *config.Secret, db *sql.DB, logger *logger.Logger, rmpq rmqp.AMQP) *fiber.App {
	return mapUrls(cfg, db, logger, rmpq)
}
