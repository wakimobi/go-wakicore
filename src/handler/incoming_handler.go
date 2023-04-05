package handler

import (
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/providers/telco"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

type IncomingHandler struct {
	cfg                 *config.Secret
	logger              *logger.Logger
	message             rmqp.AMQP
	serviceService      services.IServiceService
	subscriptionService services.ISubscriptionService
}

func NewIncomingHandler(cfg *config.Secret, logger *logger.Logger, message rmqp.AMQP, serviceService services.IServiceService, subscriptionService services.ISubscriptionService) *IncomingHandler {
	return &IncomingHandler{
		cfg:                 cfg,
		logger:              logger,
		message:             message,
		serviceService:      serviceService,
		subscriptionService: subscriptionService,
	}
}

const (
	RMQ_DATATYPE      = "application/json"
	RMQ_MOEXCHANGE    = "E_MO"
	RMQ_MOQUEUE       = "Q_MO"
	RMQ_DREXCHANGE    = "E_DR"
	RMQ_DRQUEUE       = "Q_DR"
	MT_FIRSTPUSH      = "FIRSTPUSH"
	MT_RENEWAL        = "RENEWAL"
	STATUS_SUCCESS    = "SUCCESS"
	STATUS_FAILED     = "FAILED"
	SUBJECT_FIRSTPUSH = "FIRSTPUSH"
	SUBJECT_RENEWAL   = "RENEWAL"
	SUBJECT_UNSUB     = "UNSUB"
)

var validate = validator.New()

func ValidateStruct(data interface{}) []*entity.ErrorResponse {
	var errors []*entity.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (h *IncomingHandler) GoalySubPage(c *fiber.Ctx) error {
	return c.Render("goaly/sub", fiber.Map{
		"host": h.cfg.App.Host,
	})
}

func (h *IncomingHandler) GoalyUnsubPage(c *fiber.Ctx) error {
	return c.Render("goaly/unsub", fiber.Map{
		"host": h.cfg.App.Host,
	})
}

func (h *IncomingHandler) GoalySuccessPage(c *fiber.Ctx) error {
	return c.Render("goaly/page-success", fiber.Map{
		"host": h.cfg.App.Host,
	})
}

func (h *IncomingHandler) GoalyCancelPage(c *fiber.Ctx) error {
	return c.Render("goaly/page-fail", fiber.Map{
		"host": h.cfg.App.Host,
	})
}

func (h *IncomingHandler) GoalyTermPage(c *fiber.Ctx) error {
	return c.Render("goaly/term", fiber.Map{
		"host": h.cfg.App.Host,
	})
}

func (h *IncomingHandler) GoalyHEPage(c *fiber.Ctx) error {
	return c.JSON(c.GetReqHeaders())
}

func (h *IncomingHandler) GoalyOptIn(c *fiber.Ctx) error {
	req := new(entity.ReqOptInParam)

	err := c.BodyParser(req)
	if err != nil {
		log.Println(err)
	}
	var sub *entity.Subscription
	var content *entity.Content

	if !h.serviceService.CheckService(req.Service) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	service, err := h.serviceService.GetServiceByCode(req.Service)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}
	telco := telco.NewTelco(h.cfg, h.logger, sub, service, content)

	redirect, err := telco.WebOptInOTP()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    false,
		"message":  "Success",
		"redirect": redirect,
	})
}

func (h *IncomingHandler) GoalyCallback(c *fiber.Ctx) error {
	return c.Redirect(h.cfg.Portal.Url)
}

func (h *IncomingHandler) MessageOriginated(c *fiber.Ctx) error {
	l := h.logger.Init("mo", true)

	/**
	 * Query Parser
	 */
	req := new(entity.ReqMOParams)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			entity.ResponseMO{
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.IpAddress = c.Get("Cf-Connecting-Ip")
	} else {
		req.IpAddress = c.Get("X-Forwarded-For")
	}
	json, _ := json.Marshal(req)

	h.message.IntegratePublish(RMQ_MOEXCHANGE, RMQ_MOQUEUE, RMQ_DATATYPE, "", string(json))

	l.WithFields(logrus.Fields{"request": req}).Info("MO")

	return c.Status(fiber.StatusOK).JSON(entity.ResponseMO{
		StatusCode: fiber.StatusOK,
		Message:    "Successful",
	})
}

func (h *IncomingHandler) DeliveryReport(c *fiber.Ctx) error {
	l := h.logger.Init("dr", true)

	/**
	 * Query Parser
	 */
	req := new(entity.ReqMOParams)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ResponseDR{
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.IpAddress = c.Get("Cf-Connecting-Ip")
	} else {
		req.IpAddress = c.Get("X-Forwarded-For")
	}

	json, _ := json.Marshal(req)

	h.message.IntegratePublish(RMQ_DREXCHANGE, RMQ_DRQUEUE, RMQ_DATATYPE, "", string(json))

	l.WithFields(logrus.Fields{"request": req}).Info("DR")

	return c.Status(fiber.StatusOK).JSON(entity.ResponseMO{
		StatusCode: fiber.StatusOK,
		Message:    "Successful",
	})
}

func (h *IncomingHandler) AveragePerUser(c *fiber.Ctx) error {
	subs := h.subscriptionService.AveragePerUser()
	return c.Status(fiber.StatusOK).JSON(subs)
}
