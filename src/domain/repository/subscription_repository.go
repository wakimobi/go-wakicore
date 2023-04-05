package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/domain/entity"
)

const (
	queryInsertSubscription      = "INSERT INTO subscriptions(category, service_id, msisdn, channel, adnet, pub_id, aff_sub, latest_trxid, latest_keyword, latest_subject, latest_status, latest_pin, success, ip_address, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)"
	queryUpdateSubSuccess        = "UPDATE subscriptions SET latest_trxid = $1, latest_subject = $2, latest_status = $3, latest_pin = $4, amount = amount + $5, renewal_at = $6, charge_at = $7, success = success + $8, is_retry = $9, charging_count = charging_count + $10, charging_count_all = charging_count_all + $11, total_firstpush = total_firstpush + $12, total_renewal = total_renewal + $13, total_amount_firstpush = total_amount_firstpush + $14, total_amount_renewal = total_amount_renewal + $15, updated_at = NOW() WHERE service_id = $16 AND msisdn = $17"
	queryUpdateSubFailed         = "UPDATE subscriptions SET latest_trxid = $1, latest_subject = $2, latest_status = $3, renewal_at = $4, retry_at = $5, is_retry = $6, updated_at = NOW() WHERE service_id = $7 AND msisdn = $8"
	queryUpdateSubLatest         = "UPDATE subscriptions SET latest_trxid = $1, latest_keyword = $2, latest_subject = $3, latest_status = $4, updated_at = NOW() WHERE service_id = $5 AND msisdn = $6"
	queryUpdateSubEnable         = "UPDATE subscriptions SET channel = $1, adnet = $2, pub_id = $3, aff_sub = $4, latest_trxid = $5, latest_keyword = $6, latest_subject = $7, latest_status = $8, ip_address = $9, is_retry = $10, is_active = $11, updated_at = NOW() WHERE service_id = $12 AND msisdn = $13"
	queryUpdateSubDisable        = "UPDATE subscriptions SET channel = $1, latest_trxid = $2, latest_keyword = $3, latest_subject = $4, latest_status = $5, unsub_at = $6, ip_address = $7, is_retry = $8, is_active = $9, charging_count = $10, updated_at = NOW() WHERE service_id = $11 AND msisdn = $12"
	queryCountSubscription       = "SELECT COUNT(*) as count FROM subscriptions WHERE service_id = $1 AND msisdn = $2"
	queryCountActiveSubscription = "SELECT COUNT(*) as count FROM subscriptions WHERE service_id = $1 AND msisdn = $2 AND is_active = true"
	querySelectSubscription      = "SELECT id, service_id, msisdn, channel, adnet, pub_id, aff_sub, latest_subject, latest_status, amount, trial_at, renewal_at, unsub_at, charge_at, retry_at, success, ip_address, total_firstpush, total_renewal, total_amount_firstpush, total_amount_renewal, is_trial, is_retry, is_active FROM subscriptions WHERE service_id = $1 AND msisdn = $2"
	querySelectPopulateRenewal   = "SELECT id, msisdn, service_id, channel, latest_keyword, latest_pin, ip_address, created_at FROM subscriptions WHERE renewal_at IS NOT NULL AND DATE(renewal_at) <= DATE(NOW()) AND is_active = true ORDER BY success DESC"
	querySelectPopulateRetry     = "SELECT id, msisdn, service_id, channel, latest_keyword, latest_pin, ip_address, created_at FROM subscriptions WHERE renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '1 day') AND is_retry = true AND is_active = true ORDER BY success DESC"
	querySelectPopulateReminder  = "SELECT id, msisdn, service_id, channel, latest_keyword, latest_pin, ip_address, created_at FROM subscriptions WHERE renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '2 day') AND is_retry = false AND is_active = true ORDER BY success DESC"
	querySelectArpu              = "SELECT c.name, a.adnet, COUNT(a.id) as subs, COUNT(b.id) as subs_active, SUM(a.amount) as revenue FROM subscriptions a LEFT JOIN subscriptions b ON b.service_id = a.service_id AND b.msisdn = a.msisdn AND b.is_active = true LEFT JOIN services c ON c.id = a.service_id WHERE DATE(a.created_at) BETWEEN DATE($1) AND DATE($2) AND DATE(a.renewal_at) >= DATE($3) GROUP BY a.adnet, a.service_id, c.name ORDER BY SUM(a.total_amount_renewal + a.total_amount_firstpush) DESC"
	// queryUpdateSubscription      = "UPDATE subscriptions SET channel = $1, adnet = $2, pub_id = $3, aff_sub = $4, latest_trxid = $5, latest_keyword = $6, latest_subject = $7, latest_status = $8, amount = amount + $9, trial_at = $10, renewal_at = $11, unsub_at = $12, charge_at = $13, retry_at = $14, success = success + $15, ip_address = $16, is_trial = $17, is_retry = $18, is_active = $19, total_firstpush = total_firstpush + $20, total_renewal = total_renewal + $21, total_amount_firstpush = total_amount_firstpush + $22, total_amount_renewal = total_amount_renewal + $23, updated_at = NOW() WHERE service_id = $24 AND msisdn = $25"
)

type SubscriptionRepository struct {
	db *sql.DB
}

type ISubscriptionRepository interface {
	Save(*entity.Subscription) error
	UpdateSuccess(*entity.Subscription) error
	UpdateFailed(*entity.Subscription) error
	UpdateLatest(*entity.Subscription) error
	UpdateEnable(*entity.Subscription) error
	UpdateDisable(*entity.Subscription) error
	Count(int, string) (int, error)
	CountActive(int, string) (int, error)
	Get(int, string) (*entity.Subscription, error)
	Renewal() (*[]entity.Subscription, error)
	Retry() (*[]entity.Subscription, error)
	Reminder() (*[]entity.Subscription, error)
	AveragePerUser(string, string, string) (*[]entity.AveragePerUser, error)
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Save(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertSubscription)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Category, s.ServiceID, s.Msisdn, s.Channel, s.Adnet, s.PubID, s.AffSub, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.LatestPIN, s.Success, s.IpAddress, s.IsActive, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions created ", rows)
	return nil
}

func (r *SubscriptionRepository) UpdateSuccess(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubSuccess)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestSubject, s.LatestStatus, s.LatestPIN, s.Amount, s.RenewalAt, s.ChargeAt, s.Success, s.IsRetry, s.ChargingCount, s.ChargingcountAll, s.TotalFirstpush, s.TotalRenewal, s.TotalAmountFirstpush, s.TotalAmountRenewal, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateFailed(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubFailed)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestSubject, s.LatestStatus, s.RenewalAt, s.RetryAt, s.IsRetry, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateLatest(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubLatest)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateEnable(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubEnable)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Channel, s.Adnet, s.PubID, s.AffSub, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.IpAddress, s.IsRetry, s.IsActive, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateDisable(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubDisable)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Channel, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.UnsubAt, s.IpAddress, s.IsRetry, s.IsActive, s.ChargingCount, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) Count(serviceId int, msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountSubscription, serviceId, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) CountActive(serviceId int, msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountActiveSubscription, serviceId, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) Get(serviceId int, msisdn string) (*entity.Subscription, error) {
	var s entity.Subscription
	err := r.db.QueryRow(querySelectSubscription, serviceId, msisdn).Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.PubID, &s.AffSub, &s.LatestSubject, &s.LatestStatus, &s.Amount, &s.TrialAt, &s.RenewalAt, &s.UnsubAt, &s.ChargeAt, &s.RetryAt, &s.Success, &s.IpAddress, &s.TotalFirstpush, &s.TotalRenewal, &s.TotalAmountFirstpush, &s.TotalAmountRenewal, &s.IsTrial, &s.IsRetry, &s.IsActive)
	if err != nil {
		return &s, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) Renewal() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRenewal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.Msisdn, &s.ServiceID, &s.Channel, s.LatestKeyword, s.LatestPIN, &s.IpAddress, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) Retry() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRetry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.Msisdn, &s.ServiceID, &s.Channel, s.LatestPIN, &s.IpAddress, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) Reminder() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateReminder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.Msisdn, &s.ServiceID, &s.Channel, s.LatestPIN, &s.IpAddress, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) AveragePerUser(start, end, renewal string) (*[]entity.AveragePerUser, error) {
	rows, err := r.db.Query(querySelectArpu, start, end, renewal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.AveragePerUser

	for rows.Next() {
		var s entity.AveragePerUser
		if err := rows.Scan(&s.Name, &s.Adnet, &s.Subs, &s.SubsActive, &s.Revenue); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}
