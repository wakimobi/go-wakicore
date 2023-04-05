package entity

import "time"

type Transaction struct {
	ID           int64  `json:"id"`
	TxID         string `json:"tx_id"`
	ServiceID    int    `json:"service_id"`
	Service      *Service
	Msisdn       string    `json:"msisdn"`
	Channel      string    `json:"channel"`
	Adnet        string    `json:"adnet"`
	PubID        string    `json:"pub_id"`
	AffSub       string    `json:"aff_sub"`
	Keyword      string    `json:"keyword"`
	PIN          int       `json:"pin"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	StatusCode   string    `json:"status_code"`
	StatusDetail string    `json:"status_detail"`
	Subject      string    `json:"subject"`
	IpAddress    string    `json:"ip_address"`
	Payload      string    `json:"payload"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (t *Transaction) SetAmount(amount float64) {
	t.Amount = amount
}

func (t *Transaction) SetStatus(status string) {
	t.Status = status
}

func (t *Transaction) SetStatusCode(statusCode string) {
	t.StatusCode = statusCode
}

func (t *Transaction) SetStatusDetail(statusDetail string) {
	t.StatusDetail = statusDetail
}

func (t *Transaction) SetSubject(subject string) {
	t.Subject = subject
}
