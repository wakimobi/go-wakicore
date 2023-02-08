package subscriptions

import (
	"time"

	"github.com/wakimobi/go-wakicore/src/domain/products"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

type Subscription struct {
	ID               int64 `json:"id"`
	ProductID        int   `json:"product_id"`
	Product          *products.Product
	Msisdn           string     `json:"msisdn"`
	Adnet            string     `json:"adnet"`
	LatestSubject    string     `json:"latest_subject"`
	LatestStatus     string     `json:"latest_status"`
	Amount           float64    `json:"amount"`
	RenewalAt        *time.Time `json:"renewal_at"`
	PurgeAt          *time.Time `json:"purge_at"`
	UnsubAt          *time.Time `json:"unsub_at"`
	ChargeAt         *time.Time `json:"charge_at"`
	RetryAt          *time.Time `json:"retry_at"`
	SuccessFirstpush uint       `json:"success_firstpush"`
	SuccessRenewal   uint       `json:"success_renewal"`
	TotalSuccess     uint       `json:"total_success"`
	TotalFirstpush   float64    `json:"total_firstpush"`
	TotalRenewal     float64    `json:"total_renewal"`
	TotalAmount      float64    `json:"total_amount"`
	IsRetry          bool       `json:"is_retry"`
	IsTrial          bool       `json:"is_trial"`
	IsPurge          bool       `json:"is_purge"`
	IsSuspend        bool       `json:"is_suspend"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

func (s *Subscription) Validate() *errors.RestErr {
	if s.ProductID == 0 {
		return errors.NewBadRequestError("invalid product_id")
	}

	if s.Msisdn == "" {
		return errors.NewBadRequestError("invalid msisdn")
	}

	return nil
}
