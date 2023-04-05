package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/domain/entity"
)

const (
	queryInsertTransaction = "INSERT INTO transactions(tx_id, service_id, msisdn, channel, adnet, pub_id, aff_sub, keyword, pin, amount, status, status_code, status_detail, subject, ip_address, payload, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)"
	queryDeleteTransaction = "DELETE FROM transactions WHERE service_id = $1 AND msisdn = $2 AND subject = $3 AND status = $4 AND DATE(created_at) = DATE($5)"
)

type TransactionRepository struct {
	db *sql.DB
}

type ITransactionRepository interface {
	Save(*entity.Transaction) error
	Delete(*entity.Transaction) error
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Save(t *entity.Transaction) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertTransaction)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.TxID, t.ServiceID, t.Msisdn, t.Channel, t.Adnet, t.PubID, t.AffSub, t.Keyword, t.PIN, t.Amount, t.Status, t.StatusCode, t.StatusDetail, t.Subject, t.IpAddress, t.Payload, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into transactions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d transactions created ", rows)
	return nil
}

func (r *TransactionRepository) Delete(t *entity.Transaction) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryDeleteTransaction)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, t.ServiceID, t.Msisdn, t.Subject, t.Status, time.Now())
	if err != nil {
		log.Printf("Error %s when remove row into transactions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d transactions deleted ", rows)
	return nil
}
