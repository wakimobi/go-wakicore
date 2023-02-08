package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/subscriptions"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func CreateSubscription(sub subscriptions.Subscription) (*subscriptions.Subscription, *errors.RestErr) {
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if err := sub.Save(); err != nil {
		return nil, err
	}
	return &sub, nil
}

func GetSubscription(productId int, msisdn string) (*subscriptions.Subscription, *errors.RestErr) {
	result := subscriptions.Subscription{
		ProductID: productId,
		Msisdn:    msisdn,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}

func CountSubscription(productId int, msisdn string) (int, *errors.RestErr) {
	result := subscriptions.Subscription{
		ProductID: productId,
		Msisdn:    msisdn,
	}
	count, err := result.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateSubscription(isPartial bool, sub subscriptions.Subscription) (*subscriptions.Subscription, *errors.RestErr) {
	current, err := GetSubscription(sub.ProductID, sub.Msisdn)
	if err != nil {
		return nil, err
	}

	if isPartial {

		if sub.Adnet != "" {
			current.Adnet = sub.Adnet
		}

		if sub.LatestSubject != "" {
			current.LatestSubject = sub.LatestSubject
		}

		if sub.LatestStatus != "" {
			current.LatestStatus = sub.LatestStatus
		}

		if !sub.RenewalAt.IsZero() {
			current.RenewalAt = sub.RenewalAt
		}

		if !sub.PurgeAt.IsZero() {
			current.PurgeAt = sub.PurgeAt
		}

		if !sub.UnsubAt.IsZero() {
			current.UnsubAt = sub.UnsubAt
		}

		if !sub.ChargeAt.IsZero() {
			current.ChargeAt = sub.ChargeAt
		}

		if !sub.RetryAt.IsZero() {
			current.RetryAt = sub.RetryAt
		}

		if sub.SuccessFirstpush != 0 {
			current.SuccessFirstpush = sub.SuccessFirstpush
		}

		if sub.SuccessRenewal != 0 {
			current.SuccessRenewal = sub.SuccessRenewal
		}

		if sub.TotalSuccess != 0 {
			current.TotalSuccess = sub.TotalSuccess
		}

		if sub.TotalFirstpush != 0 {
			current.TotalFirstpush = sub.TotalFirstpush
		}

		if sub.TotalRenewal != 0 {
			current.TotalRenewal = sub.TotalRenewal
		}

		if sub.TotalAmount != 0 {
			current.TotalAmount = sub.TotalAmount
		}

		if sub.IsRetry || !sub.IsRetry {
			current.IsRetry = sub.IsRetry
		}

		if sub.IsTrial || !sub.IsTrial {
			current.IsTrial = sub.IsTrial
		}

		if sub.IsPurge || !sub.IsPurge {
			current.IsPurge = sub.IsPurge
		}

		if sub.IsSuspend {
			current.IsSuspend = sub.IsSuspend
		}

		if sub.IsActive || !sub.IsActive {
			current.IsActive = sub.IsActive
		}

	} else {
		current.RenewalAt = sub.RenewalAt
		current.LatestSubject = sub.LatestSubject
		current.LatestStatus = sub.LatestStatus
		current.IsActive = sub.IsActive
	}

	return current, nil
}

func SearchSubscription(status string) ([]subscriptions.Subscription, *errors.RestErr) {
	dao := &subscriptions.Subscription{}
	return dao.FindByStatus(status)
}

func ChargeSubscription(productId int, msisdn string) *errors.RestErr {
	result := subscriptions.Subscription{
		ProductID: productId,
		Msisdn:    msisdn,
	}
	if err := result.Get(); err != nil {
		return err
	}

	return nil
}
