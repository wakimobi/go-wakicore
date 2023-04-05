package handler

import (
	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/services"
)

type DRHandler struct {
	cfg                 *config.Secret
	logger              *logger.Logger
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	req                 *entity.ReqDRParams
}

func NewDRHandler(
	cfg *config.Secret,
	logger *logger.Logger,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	req *entity.ReqDRParams) *DRHandler {

	return &DRHandler{
		cfg:                 cfg,
		logger:              logger,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		req:                 req,
	}
}
