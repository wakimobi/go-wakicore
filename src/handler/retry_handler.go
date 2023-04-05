package handler

import (
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/providers/portal"
	"github.com/idprm/go-pass-tsel/src/providers/telco"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/idprm/go-pass-tsel/src/utils/pin_utils"
	"github.com/idprm/go-pass-tsel/src/utils/response_utils"
	"github.com/idprm/go-pass-tsel/src/utils/uuid_utils"
)

type RetryHandler struct {
	cfg                 *config.Secret
	logger              *logger.Logger
	sub                 *entity.Subscription
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
}

func NewRetryHandler(
	cfg *config.Secret,
	logger *logger.Logger,
	sub *entity.Subscription,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
) *RetryHandler {

	return &RetryHandler{
		cfg:                 cfg,
		logger:              logger,
		sub:                 sub,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
	}
}

func (h *RetryHandler) getContentFirstpush(serviceId, pin int) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_FIRSTPUSH, pin)
}

func (h *RetryHandler) getContentRenewal(serviceId, pin int) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_RENEWAL, pin)
}

func (h *RetryHandler) Firstpush() {
	pin := pin_utils.Generate(100000, 999999)
	service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())
	content, _ := h.getContentFirstpush(h.sub.GetServiceId(), pin)
	smsMT := telco.NewTelco(h.cfg, h.logger, h.sub, service, content)
	resp, err := smsMT.SMSbyParam()
	if err != nil {
		log.Println(err)
	}
	trxId := uuid_utils.GenerateTrxId()

	if response_utils.ParseStatus(string(resp)) {

		subSuccess := &entity.Subscription{
			ServiceID:            h.sub.GetServiceId(),
			Msisdn:               h.sub.GetMsisdn(),
			LatestTrxId:          trxId,
			LatestStatus:         STATUS_SUCCESS,
			LatestSubject:        SUBJECT_FIRSTPUSH,
			Amount:               service.GetPrice(),
			RenewalAt:            time.Now().AddDate(0, 0, service.GetRenewalDay()),
			ChargeAt:             time.Now(),
			Success:              1,
			TotalFirstpush:       1,
			TotalAmountFirstpush: service.GetPrice(),
			IsRetry:              false,
		}

		h.subscriptionService.UpdateSuccess(subSuccess)

		transSuccess := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    h.sub.GetServiceId(),
			Msisdn:       h.sub.GetMsisdn(),
			Channel:      h.sub.GetChannel(),
			Keyword:      h.sub.GetLatestKeyword(),
			Amount:       service.GetPrice(),
			Status:       STATUS_SUCCESS,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      string(resp),
			IpAddress:    h.sub.GetIpAddress(),
		}

		h.transactionService.UpdateTransaction(transSuccess)

		h.portalSubscription(h.sub, service, pin)
	}
}

func (h *RetryHandler) Dailypush() {
	service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())
	content, _ := h.getContentRenewal(h.sub.GetServiceId(), h.sub.GetLatestPIN())
	smsMT := telco.NewTelco(h.cfg, h.logger, h.sub, service, content)
	resp, err := smsMT.SMSbyParam()
	if err != nil {
		log.Println(err)
	}
	trxId := uuid_utils.GenerateTrxId()

	if response_utils.ParseStatus(string(resp)) {
		service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())

		subSuccess := &entity.Subscription{
			ServiceID:          h.sub.GetServiceId(),
			Msisdn:             h.sub.GetMsisdn(),
			LatestTrxId:        trxId,
			LatestStatus:       STATUS_SUCCESS,
			LatestSubject:      SUBJECT_RENEWAL,
			LatestPIN:          h.sub.GetLatestPIN(),
			Amount:             service.GetPrice(),
			RenewalAt:          time.Now().AddDate(0, 0, service.GetRenewalDay()),
			ChargeAt:           time.Now(),
			Success:            1,
			TotalRenewal:       1,
			TotalAmountRenewal: service.GetPrice(),
			IsRetry:            false,
		}

		h.subscriptionService.UpdateSuccess(subSuccess)

		transSuccess := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    h.sub.GetServiceId(),
			Msisdn:       h.sub.GetMsisdn(),
			Channel:      h.sub.GetChannel(),
			Keyword:      h.sub.GetLatestKeyword(),
			Amount:       service.GetPrice(),
			PIN:          h.sub.GetLatestPIN(),
			Status:       STATUS_SUCCESS,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_RENEWAL,
			Payload:      string(resp),
			IpAddress:    h.sub.GetIpAddress(),
		}

		h.transactionService.UpdateTransaction(transSuccess)

		h.portalRenewal(h.sub, service, h.sub.GetLatestPIN())
	}
}

func (h *RetryHandler) portalSubscription(sub *entity.Subscription, service *entity.Service, pin int) {
	p := portal.NewPortal(h.cfg, h.logger, sub, service, pin)
	p.Subscription()
}

func (h *RetryHandler) portalRenewal(sub *entity.Subscription, service *entity.Service, pin int) {
	p := portal.NewPortal(h.cfg, h.logger, sub, service, pin)
	p.Renewal()
}
