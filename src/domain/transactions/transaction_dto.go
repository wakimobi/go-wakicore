package transactions

import (
	"time"

	"github.com/wakimobi/go-wakicore/src/domain/products"
)

type Transaction struct {
	ID            int64  `json:"id"`
	TransactionID string `json:"transaction_id"`
	ProductID     int    `json:"-"`
	Product       *products.Product
	Msisdn        string     `json:"msisdn"`
	Adnet         string     `json:"adnet"`
	Amount        float64    `json:"amount"`
	Subject       string     `json:"subject"`
	Status        string     `json:"status"`
	StatusDetail  string     `json:"status_detail"`
	Payload       string     `json:"payload"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

func (t *Transaction) SetAmount(tr *Transaction) (*Transaction, error) {
	t.Amount = tr.Amount
	return t, nil
}

func (t *Transaction) SetSubject(tr *Transaction) (*Transaction, error) {
	t.Subject = tr.Subject
	return t, nil
}
