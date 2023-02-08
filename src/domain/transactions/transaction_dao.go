package transactions

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryInsertTransaction = "INSERT INTO transactions(transaction_id, product_id, msisdn, adnet, amount, subject, status, status_detail, payload, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id"
	queryUpdateTransaction = "UPDATE transactions SET adnet = $1, amount = $2, subject = $3, status = $4, status_detail = $5, payload = $6, updated_at = NOW() WHERE product_id = $7 AND msisdn = $8"
	queryGetTransaction    = "SELECT id, transaction_id, product_id, msisdn, adnet, amount, subject, status, status_detail, payload FROM transactions WHERE product_id = $1 AND msisdn = $2 LIMIT 1"
)

func (t *Transaction) Save() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryInsertTransaction)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()
	var transactionId int64
	insertResult := stmt.QueryRow(t.TransactionID, t.ProductID, t.Msisdn, t.Adnet, t.Amount, t.Subject, t.Status, t.StatusDetail, t.Payload)
	if getErr := insertResult.Scan(&transactionId); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	t.ID = transactionId
	return nil
}

func (t *Transaction) Update() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryUpdateTransaction)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()
	_, updateErr := stmt.Exec(t.Adnet, t.Amount, t.Subject, t.Status, t.StatusDetail, t.Payload, t.ProductID, t.Msisdn)
	if updateErr != nil {
		return pgsql_utils.ParseError(updateErr)
	}

	return nil
}

func (t *Transaction) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetTransaction)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(t.ProductID, t.Msisdn)
	if getErr := result.Scan(&t.ID, &t.TransactionID, &t.ProductID, &t.Msisdn, &t.Adnet, &t.Amount, &t.Subject, &t.Status, &t.StatusDetail, &t.Payload); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}
