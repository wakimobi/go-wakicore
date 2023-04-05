package entity

import (
	"time"
)

type Subscription struct {
	ID                   int64 `json:"id"`
	ServiceID            int   `json:"service_id"`
	Service              *Service
	Category             string    `json:"category"`
	Msisdn               string    `json:"msisdn"`
	Channel              string    `json:"channel"`
	Adnet                string    `json:"adnet"`
	PubID                string    `json:"pub_id"`
	AffSub               string    `json:"aff_sub"`
	LatestTrxId          string    `json:"trx_id"`
	LatestKeyword        string    `json:"latest_keyword"`
	LatestSubject        string    `json:"latest_subject"`
	LatestStatus         string    `json:"latest_status"`
	LatestPIN            int       `json:"latest_pin"`
	Amount               float64   `json:"amount"`
	TrialAt              time.Time `json:"trial_at"`
	RenewalAt            time.Time `json:"renewal_at"`
	UnsubAt              time.Time `json:"unsub_at"`
	ChargeAt             time.Time `json:"charge_at"`
	RetryAt              time.Time `json:"retry_at"`
	Success              uint      `json:"success"`
	IpAddress            string    `json:"ip_address"`
	TotalFirstpush       uint      `json:"total_firstpush"`
	TotalRenewal         uint      `json:"total_renewal"`
	TotalAmountFirstpush float64   `json:"total_amount_firstpush"`
	TotalAmountRenewal   float64   `json:"total_amount_renewal"`
	ChargingCount        uint      `json:"charging_count"`
	ChargingcountAll     uint      `json:"charging_count_all"`
	IsTrial              bool      `json:"is_trial"`
	IsRetry              bool      `json:"is_retry"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (s *Subscription) GetServiceId() int {
	return s.ServiceID
}

func (s *Subscription) GetCategory() string {
	return s.Category
}

func (s *Subscription) GetMsisdn() string {
	return s.Msisdn
}

func (s *Subscription) GetChannel() string {
	return s.Channel
}

func (s *Subscription) GetLatestTrxId() string {
	return s.LatestTrxId
}

func (s *Subscription) GetLatestKeyword() string {
	return s.LatestKeyword
}

func (s *Subscription) GetLatestSubject() string {
	return s.LatestSubject
}

func (s *Subscription) GetLatestStatus() string {
	return s.LatestStatus
}

func (s *Subscription) GetLatestPIN() int {
	return s.LatestPIN
}

func (s *Subscription) GetIpAddress() string {
	return s.IpAddress
}

func (s *Subscription) SetIsActive(active bool) {
	s.IsActive = active
}

func (s *Subscription) SetIsRetry(retry bool) {
	s.IsRetry = retry
}

func (s *Subscription) SetIsTrial(trial bool) {
	s.IsTrial = trial
}

func (s *Subscription) SetRenewalAt(renewalAt time.Time) {
	s.RenewalAt = renewalAt
}

func (s *Subscription) SetRetryAt(retryAt time.Time) {
	s.RetryAt = retryAt
}

func (s *Subscription) SetChargeAt(chargeAt time.Time) {
	s.ChargeAt = chargeAt
}

func (s *Subscription) SetUnsubAt(unsubAt time.Time) {
	s.UnsubAt = unsubAt
}

func (s *Subscription) SetLatestSubject(latestSubject string) {
	s.LatestSubject = latestSubject
}

func (s *Subscription) SetLatestStatus(latestStatus string) {
	s.LatestStatus = latestStatus
}

func (s *Subscription) SetChannel(channel string) {
	s.Channel = channel
}

func (s *Subscription) SetAdnet(adnet string) {
	s.Adnet = adnet
}

func (s *Subscription) SetPubID(pubId string) {
	s.PubID = pubId
}

func (s *Subscription) SetAffSub(affsub string) {
	s.AffSub = affsub
}

func (s *Subscription) GetChargingCount() uint {
	return s.ChargingCount
}

func (s *Subscription) GetChargingcountAll() uint {
	return s.ChargingcountAll
}

func (s *Subscription) IsCreatedAtToday() bool {
	return s.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02")
}
