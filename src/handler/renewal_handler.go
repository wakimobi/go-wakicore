package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/providers/portal"
	"github.com/idprm/go-pass-tsel/src/providers/telco"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/idprm/go-pass-tsel/src/utils/response_utils"
	"github.com/idprm/go-pass-tsel/src/utils/uuid_utils"
)

type RenewalHandler struct {
	cfg                 *config.Secret
	logger              *logger.Logger
	sub                 *entity.Subscription
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
}

func NewRenewalHandler(
	cfg *config.Secret,
	logger *logger.Logger,
	sub *entity.Subscription,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
) *RenewalHandler {

	return &RenewalHandler{
		cfg:                 cfg,
		logger:              logger,
		sub:                 sub,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
	}
}

func (h *RenewalHandler) getContentRenewal(serviceId, pin int) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_RENEWAL, pin)
}

func (h *RenewalHandler) Dailypush() {
	service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())
	content, _ := h.getContentRenewal(h.sub.GetServiceId(), h.sub.GetLatestPIN())
	smsMT := telco.NewTelco(h.cfg, h.logger, h.sub, service, content)
	resp, err := smsMT.SMSbyParam()
	if err != nil {
		log.Println(err)
	}
	trxId := uuid_utils.GenerateTrxId()

	if response_utils.ParseStatus(string(resp)) {
		subSuccess := &entity.Subscription{
			ServiceID:          h.sub.GetServiceId(),
			Msisdn:             h.sub.GetMsisdn(),
			LatestTrxId:        trxId,
			LatestSubject:      SUBJECT_RENEWAL,
			LatestStatus:       STATUS_SUCCESS,
			LatestPIN:          h.sub.GetLatestPIN(),
			Amount:             service.GetPrice(),
			RenewalAt:          time.Now().AddDate(0, 0, service.GetRenewalDay()),
			ChargeAt:           time.Now(),
			Success:            1,
			IsRetry:            false,
			TotalRenewal:       1,
			TotalAmountRenewal: service.GetPrice(),
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

		h.transactionService.SaveTransaction(transSuccess)

		h.portalRenewal(h.sub, service, h.sub.GetLatestPIN())

	} else {
		subFailed := &entity.Subscription{
			ServiceID:     h.sub.GetServiceId(),
			Msisdn:        h.sub.GetMsisdn(),
			LatestTrxId:   trxId,
			LatestSubject: SUBJECT_RENEWAL,
			LatestStatus:  STATUS_FAILED,
			RenewalAt:     time.Now().AddDate(0, 0, 1),
			RetryAt:       time.Now(),
			IsRetry:       true,
		}

		h.subscriptionService.UpdateFailed(subFailed)

		transFailed := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    h.sub.GetServiceId(),
			Msisdn:       h.sub.GetMsisdn(),
			Channel:      h.sub.GetChannel(),
			Keyword:      h.sub.GetLatestKeyword(),
			Status:       STATUS_FAILED,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_RENEWAL,
			Payload:      string(resp),
			IpAddress:    h.sub.GetIpAddress(),
		}

		h.transactionService.SaveTransaction(transFailed)
	}

}

func (h *RenewalHandler) portalRenewal(sub *entity.Subscription, service *entity.Service, pin int) {
	p := portal.NewPortal(h.cfg, h.logger, sub, service, pin)
	notif, err := p.Renewal()
	if err != nil {
		h.logger.Writer(err)
	}
	/**
	 *  Parsing Response Notif Renewal
	 */
	type resRenewal struct {
		Message string `json:"msg"`
	}
	var responseRenewal resRenewal
	json.Unmarshal(notif, &responseRenewal)

	if responseRenewal.Message == "User doest not exits" {
		p.SubscriptionRetry()
	}
}
