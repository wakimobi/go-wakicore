package cmd

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
	"github.com/idprm/go-pass-tsel/src/handler"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/services"
)

func moProcessor(cfg *config.Secret, db *sql.DB, logger *logger.Logger, wg *sync.WaitGroup, message []byte) {
	/**
	 * -. Check Valid Prefix
	 * -. Filter REG / UNREG
	 * -. Check Blacklist
	 * -. Check Active Sub
	 * -. Save Sub
	 * -. MT API
	 * -. Update Sub
	 * -. Save Transaction
	 * -. Save History
	 */
	campaignRepo := repository.NewCampaignRepository(db)
	campaignService := services.NewCampaignService(campaignRepo)
	blacklistRepo := repository.NewBlacklistRepository(db)
	blacklistService := services.NewBlacklistService(blacklistRepo)
	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	historyRepo := repository.NewHistoryRepository(db)
	historyService := services.NewHistoryService(historyRepo)

	var req *entity.ReqMOParams
	json.Unmarshal([]byte(message), &req)
	reqMO := entity.NewReqMOParams(req.SMS, req.Adn, req.Msisdn, req.Channel, req.IpAddress)

	handlerMO := handler.NewMOHandler(cfg, logger, campaignService, blacklistService, serviceService, contentService, subscriptionService, transactionService, historyService, reqMO)

	// check valid prefix
	if reqMO.IsInValidPrefix() {
		// invalid prefix
		log.Println("invalid prefix msisdn")
	} else if handlerMO.CheckBlacklist() {
		// blacklist msisdn
		log.Println("blacklist msisdn")
	} else {
		// check service by category
		if handlerMO.CheckService() {
			// filter REG
			if reqMO.IsREG() {
				if handlerMO.CheckService() {
					// active sub
					if !handlerMO.CheckActiveSub() {
						// Firstpush MT API
						handlerMO.Firstpush()
					}
				}
			}
			if reqMO.IsUNREG() {
				if handlerMO.CheckService() {
					// active sub
					if handlerMO.CheckActiveSub() {
						// unsub
						handlerMO.Unsub()
					}
				}
			}
		} else {
			// error keyword service
		}

	}
	wg.Done()
}

func drProcessor(cfg *config.Secret, db *sql.DB, logger *logger.Logger, wg *sync.WaitGroup, message []byte) {

	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	var req *entity.ReqDRParams
	json.Unmarshal(message, &req)

	reqDR := entity.NewReqDRParams(req.SMS, req.TrxId, req.Msisdn, req.IpAddress)

	handler.NewDRHandler(cfg, logger, serviceService, contentService, subscriptionService, transactionService, reqDR)

	wg.Done()
}

func renewalProcessor(cfg *config.Secret, db *sql.DB, logger *logger.Logger, wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	handlerRenewal := handler.NewRenewalHandler(cfg, logger, sub, serviceService, contentService, subscriptionService, transactionService)

	// Dailypush MT API
	handlerRenewal.Dailypush()

	wg.Done()
}

func retryProcessor(cfg *config.Secret, db *sql.DB, logger *logger.Logger, wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	handlerRetry := handler.NewRetryHandler(cfg, logger, sub, serviceService, contentService, subscriptionService, transactionService)

	if sub.IsCreatedAtToday() {
		handlerRetry.Firstpush()
	} else {
		handlerRetry.Dailypush()
	}

	wg.Done()
}
