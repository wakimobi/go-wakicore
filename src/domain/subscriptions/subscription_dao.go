package subscriptions

import (
	"fmt"

	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryInsertSubscription = "INSERT INTO subscriptions(product_id, msisdn, adnet, latest_subject, latest_status, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id"
	queryGetSubscription    = "SELECT id, adnet, latest_subject, latest_status, renewal_at, purge_at, unsub_at, charge_at, retry_at, success_firstpush, success_renewal, total_success, total_firstpush, total_renewal, total_amount, is_retry, is_purge, is_suspend, is_active FROM subscriptions WHERE product_id = $1 AND msisdn = $2 LIMIT 1"
	queryCountSubscription  = "SELECT COUNT(*) FROM subscriptions WHERE product_id = $1 AND msisdn = $2 AND is_active = true"
	queryUpdateSubscription = "UPDATE subscriptions SET adnet = $1, latest_subject = $2, latest_status = $3, renewal_at = $4, purge_at = $5, unsub_at = $6, charge_at = $7, retry_at = $8, success_firstpush = success_firstpush + $9, success_renewal = success_renewal + $10, total_success = total_success + $11, total_firstpush = total_firstpush + $12, total_renewal = total_renewal + $13, total_amount = total_amount + $14, is_retry = $15, is_purge = $16, is_suspend = $17, is_active = $18, updated_at = NOW() WHERE product_id = $19 AND msisdn = $20"
	queryFindSubByRenewal   = "SELECT product_id, msisdn FROM subscriptions WHERE DATE(renewal_at) <= DATE(NOW()) AND is_suspend = false AND is_active = true"
)

func (s *Subscription) Save() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryInsertSubscription)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()
	var subscriptionId int64
	insertResult := stmt.QueryRow(s.ProductID, s.Msisdn, s.Adnet, s.LatestSubject, s.LatestStatus, s.IsActive)
	if getErr := insertResult.Scan(&subscriptionId); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	s.ID = subscriptionId
	return nil
}

func (s *Subscription) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetSubscription)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.ProductID, s.Msisdn)
	if getErr := result.Scan(&s.ID, &s.Adnet, &s.LatestSubject, &s.LatestStatus, &s.RenewalAt, &s.PurgeAt, &s.UnsubAt, &s.ChargeAt, &s.RetryAt, &s.SuccessFirstpush, &s.SuccessRenewal, &s.TotalSuccess, &s.TotalFirstpush, &s.TotalRenewal, &s.TotalAmount, &s.IsRetry, &s.IsPurge, &s.IsSuspend, &s.IsActive); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}

func (s *Subscription) Count() (int, *errors.RestErr) {
	var count int
	stmt, err := db.Client.Prepare(queryCountSubscription)
	if err != nil {
		return 0, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.ProductID, s.Msisdn)
	if getErr := result.Scan(&count); getErr != nil {
		return 0, pgsql_utils.ParseErrorSkip(getErr)
	}
	return count, nil
}

func (s *Subscription) Update() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryUpdateSubscription)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(s.Adnet, s.LatestSubject, s.LatestStatus, s.RenewalAt, s.PurgeAt, s.UnsubAt, s.ChargeAt, s.RetryAt, s.SuccessFirstpush, s.SuccessRenewal, s.TotalSuccess, s.TotalFirstpush, s.TotalRenewal, s.TotalAmount, s.IsRetry, s.IsPurge, s.IsActive)
	if updateErr != nil {
		return pgsql_utils.ParseError(updateErr)
	}

	return nil
}

func (s *Subscription) FindByStatus(status string) ([]Subscription, *errors.RestErr) {
	stmt, err := db.Client.Prepare(queryFindSubByRenewal)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	results := make([]Subscription, 0)
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ProductID, &s.Msisdn); err != nil {
			return nil, pgsql_utils.ParseError(err)
		}
		results = append(results, s)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no subscriptions matching status %s", status))
	}

	return results, nil
}
