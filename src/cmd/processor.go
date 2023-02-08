package cmd

import (
	"encoding/json"
	"sync"

	"github.com/wakimobi/go-wakicore/src/common"
	"github.com/wakimobi/go-wakicore/src/domain/subscriptions"
	"github.com/wakimobi/go-wakicore/src/domain/transactions"
	"github.com/wakimobi/go-wakicore/src/services"
)

func moProcessor(wg *sync.WaitGroup, message []byte) {
	// parsing string json
	var req common.MORequest
	json.Unmarshal(message, &req)

	// get product
	product, _ := services.GetProduct(req.Sms)

	sub, _ := services.CreateSubscription(
		subscriptions.Subscription{
			ProductID:     product.ID,
			Msisdn:        req.Msisdn,
			Adnet:         req.Adn,
			LatestSubject: "SUCCESS",
			LatestStatus:  "INPUT_MSIDN",
		})

	services.CreateTransaction(transactions.Transaction{
		ProductID: product.ID,
		Msisdn:    req.Msisdn,
		Adnet:     req.Adn,
	})

	sub.LatestSubject = "FAILED"
	sub.Update()

	wg.Done()
}

func drProcessor(wg *sync.WaitGroup, message []byte) {
	// parsing string json
	var req common.DRRequest
	json.Unmarshal(message, &req)

	wg.Done()
}

func renewalProcessor(wg *sync.WaitGroup, message []byte) {
	wg.Done()
}

func retryProcessor(wg *sync.WaitGroup, message []byte) {
	wg.Done()
}
