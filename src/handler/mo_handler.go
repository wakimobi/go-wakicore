package handler

import (
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/providers/portal"
	"github.com/idprm/go-pass-tsel/src/providers/postback"
	"github.com/idprm/go-pass-tsel/src/providers/telco"
	"github.com/idprm/go-pass-tsel/src/services"
	"github.com/idprm/go-pass-tsel/src/utils/pin_utils"
	"github.com/idprm/go-pass-tsel/src/utils/response_utils"
	"github.com/idprm/go-pass-tsel/src/utils/uuid_utils"
)

type MOHandler struct {
	cfg                 *config.Secret
	logger              *logger.Logger
	campaignService     services.ICampaignService
	blacklistService    services.IBlacklistService
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	historyService      services.IHistoryService
	req                 *entity.ReqMOParams
}

func NewMOHandler(
	cfg *config.Secret,
	logger *logger.Logger,
	campaignService services.ICampaignService,
	blacklistService services.IBlacklistService,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	historyService services.IHistoryService,
	req *entity.ReqMOParams,
) *MOHandler {

	return &MOHandler{
		cfg:                 cfg,
		logger:              logger,
		campaignService:     campaignService,
		blacklistService:    blacklistService,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		historyService:      historyService,
		req:                 req,
	}
}

func (h *MOHandler) getService() (*entity.Service, error) {
	keyword := h.req.GetSubKeyword()
	return h.serviceService.GetServiceByCode(keyword)
}

func (h *MOHandler) CheckBlacklist() bool {
	return h.blacklistService.GetBlacklist(h.req.Msisdn)
}

func (h *MOHandler) CheckActiveSub() bool {
	service, _ := h.getService()
	return h.subscriptionService.GetActiveSubscription(service.GetID(), h.req.GetMsisdn())
}

func (h *MOHandler) checkSub() bool {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	return h.subscriptionService.GetSubscription(service.GetID(), h.req.GetMsisdn())
}

func (h *MOHandler) CheckService() bool {
	keyword := h.req.GetSubKeyword()
	return h.serviceService.CheckService(keyword)
}

func (h *MOHandler) getContentFirstpush(serviceId, pin int) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_FIRSTPUSH, pin)
}

func (h *MOHandler) Firstpush() {
	pin := pin_utils.Generate(100000, 999999)
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	content, _ := h.getContentFirstpush(service.GetID(), pin)
	if err != nil {
		log.Println(err)
	}
	channel := response_utils.ParseChannel(h.req.SMS)
	trxId := uuid_utils.GenerateTrxId()

	subscription := &entity.Subscription{
		ServiceID:     service.GetID(),
		Category:      service.GetCategory(),
		Msisdn:        h.req.GetMsisdn(),
		LatestKeyword: h.req.GetKeyword(),
		LatestSubject: SUBJECT_FIRSTPUSH,
		Channel:       channel,
		Adnet:         "",
		PubID:         "",
		AffSub:        "",
		IsActive:      true,
		IpAddress:     h.req.GetIpAddress(),
	}

	if h.checkSub() {
		h.subscriptionService.UpdateEnable(subscription)

	} else {
		h.subscriptionService.SaveSubscription(subscription)
	}

	smsMT := telco.NewTelco(h.cfg, h.logger, subscription, service, content)
	resp, err := smsMT.SMSbyParam()
	if err != nil {
		log.Println(err)
	}

	if response_utils.ParseStatus(string(resp)) {
		subSuccess := &entity.Subscription{
			ServiceID:            service.GetID(),
			Msisdn:               h.req.GetMsisdn(),
			LatestTrxId:          trxId,
			LatestSubject:        SUBJECT_FIRSTPUSH,
			LatestStatus:         STATUS_SUCCESS,
			LatestPIN:            pin,
			Amount:               service.GetPrice(),
			RenewalAt:            time.Now().AddDate(0, 0, service.GetRenewalDay()),
			ChargeAt:             time.Now(),
			Success:              1,
			IsRetry:              false,
			TotalFirstpush:       1,
			TotalAmountFirstpush: service.GetPrice(),
		}

		h.subscriptionService.UpdateSuccess(subSuccess)

		transSuccess := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetID(),
			Msisdn:       h.req.GetMsisdn(),
			Channel:      channel,
			Adnet:        h.req.GetAdnet(),
			Keyword:      h.req.GetKeyword(),
			Amount:       service.GetPrice(),
			PIN:          pin,
			Status:       STATUS_SUCCESS,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      string(resp),
			IpAddress:    h.req.GetIpAddress(),
		}

		h.transactionService.SaveTransaction(transSuccess)

		historySuccess := &entity.History{
			ServiceID: service.GetID(),
			Msisdn:    h.req.GetMsisdn(),
			Channel:   channel,
			Adnet:     "",
			Keyword:   h.req.GetKeyword(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_SUCCESS,
			IpAddress: h.req.GetIpAddress(),
		}

		h.historyService.SaveHistory(historySuccess)
		h.portalSubscription(subSuccess, service, pin)
	} else {

		subFailed := &entity.Subscription{
			ServiceID:     service.GetID(),
			Msisdn:        h.req.GetMsisdn(),
			LatestTrxId:   trxId,
			LatestSubject: SUBJECT_FIRSTPUSH,
			LatestStatus:  STATUS_FAILED,
			RenewalAt:     time.Now().AddDate(0, 0, 1),
			RetryAt:       time.Now(),
			IsRetry:       true,
		}

		h.subscriptionService.UpdateFailed(subFailed)

		transFailed := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetID(),
			Msisdn:       h.req.GetMsisdn(),
			Channel:      channel,
			Adnet:        h.req.GetAdnet(),
			Keyword:      h.req.GetKeyword(),
			Status:       STATUS_FAILED,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      string(resp),
			IpAddress:    h.req.GetIpAddress(),
		}
		h.transactionService.SaveTransaction(transFailed)

		historyFailed := &entity.History{
			ServiceID: service.GetID(),
			Msisdn:    h.req.GetMsisdn(),
			Channel:   channel,
			Adnet:     "",
			Keyword:   h.req.GetKeyword(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_FAILED,
			IpAddress: h.req.GetIpAddress(),
		}
		h.historyService.SaveHistory(historyFailed)
	}

	h.postback(subscription, service)
}

func (h *MOHandler) Unsub() {

	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	channel := response_utils.ParseChannel(h.req.SMS)
	trxId := uuid_utils.GenerateTrxId()

	subscription := &entity.Subscription{
		ServiceID:     service.GetID(),
		Msisdn:        h.req.GetMsisdn(),
		Channel:       channel,
		LatestTrxId:   trxId,
		LatestKeyword: h.req.GetKeyword(),
		LatestSubject: SUBJECT_UNSUB,
		LatestStatus:  STATUS_SUCCESS,
		UnsubAt:       time.Now(),
		IpAddress:     h.req.GetIpAddress(),
		ChargingCount: 0,
		IsRetry:       false,
		IsActive:      false,
	}
	h.subscriptionService.UpdateDisable(subscription)

	transaction := &entity.Transaction{
		TxID:         trxId,
		ServiceID:    service.GetID(),
		Msisdn:       h.req.GetMsisdn(),
		Channel:      channel,
		Keyword:      h.req.GetKeyword(),
		Status:       STATUS_SUCCESS,
		StatusCode:   "-",
		StatusDetail: "-",
		Subject:      SUBJECT_UNSUB,
		Payload:      "-",
		IpAddress:    h.req.GetIpAddress(),
	}
	h.transactionService.SaveTransaction(transaction)

	history := &entity.History{
		ServiceID: service.GetID(),
		Msisdn:    h.req.GetMsisdn(),
		Channel:   channel,
		Adnet:     "",
		Keyword:   h.req.GetKeyword(),
		Subject:   SUBJECT_UNSUB,
		Status:    STATUS_SUCCESS,
		IpAddress: h.req.GetIpAddress(),
	}
	h.historyService.SaveHistory(history)

	h.portalUnsubscription(subscription, service)
}

func (h *MOHandler) portalSubscription(sub *entity.Subscription, service *entity.Service, pin int) {
	p := portal.NewPortal(h.cfg, h.logger, sub, service, pin)
	p.Subscription()
}

func (h *MOHandler) portalUnsubscription(sub *entity.Subscription, service *entity.Service) {
	p := portal.NewPortal(h.cfg, h.logger, sub, service, 0)
	p.Unsubscription()
}

func (h *MOHandler) postback(sub *entity.Subscription, service *entity.Service) {
	p := postback.NewPostback(h.cfg, h.logger, sub, service)
	p.Send()
}
