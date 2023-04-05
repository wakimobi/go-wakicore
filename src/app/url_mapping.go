package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	log_access "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
	"github.com/idprm/go-pass-tsel/src/handler"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/wiliehidayat87/rmqp"
)

func mapUrls(cfg *config.Secret, db *sql.DB, logger *logger.Logger, rmpq rmqp.AMQP) *fiber.App {
	engine := html.New("./src/presenter/views", ".html")

	/**
	 * Init Fiber
	 */
	router := fiber.New(fiber.Config{
		Views: engine,
	})

	/**
	 * Access log on browser
	 */
	router.Use("/logs", filesystem.New(filesystem.Config{
		Root:         http.Dir(cfg.Log.Path),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	/**
	 * Write access logger
	 */
	file, err := os.OpenFile(cfg.Log.Path+"/access_log/log-"+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	router.Use(requestid.New())
	router.Use(log_access.New(log_access.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   cfg.App.TimeZone,
		Output:     file,
	}))

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	router.Static("/static", path+"/public")

	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	incomingHandler := handler.NewIncomingHandler(cfg, logger, rmpq, serviceService, subscriptionService)

	/**
	 * Routes Landing Page SUB & UNSUB
	 */
	router.Get("goaly", incomingHandler.GoalySubPage)
	router.Get("goaly/unsub", incomingHandler.GoalyUnsubPage)
	router.Get("term", incomingHandler.GoalyTermPage)
	router.Get("goaly/he", incomingHandler.GoalyHEPage)
	router.Get("success", incomingHandler.GoalySuccessPage)
	router.Get("cbtsel", incomingHandler.GoalyCallback)
	router.Get("cancel", incomingHandler.GoalyCancelPage)
	router.Post("goaly", incomingHandler.GoalyOptIn)

	/**
	 * Routes MO & DR
	 */
	router.Get("notif/mo", incomingHandler.MessageOriginated)
	router.Get("notif/dr", incomingHandler.DeliveryReport)
	router.Get("mo", incomingHandler.MessageOriginated)
	router.Get("dr", incomingHandler.DeliveryReport)

	/**
	 * Routes Report
	 */
	router.Post("arpu", incomingHandler.AveragePerUser)

	return router
}
