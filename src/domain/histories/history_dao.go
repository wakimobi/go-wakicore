package histories

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryInsertHistory = "INSERT INTO histories(product_id, msisdn, adnet, subject, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id"
	queryGetHistory    = "SELECT id, product_id, msisdn, adnet, subject, created_at FROM histories WHERE product_id = $1 AND msisdn = $2"
)

func (h *History) Save() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryInsertHistory)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()
	var historyId int64
	insertResult := stmt.QueryRow(h.ProductID, h.Msisdn, h.Adnet, h.Subject)
	if getErr := insertResult.Scan(&historyId); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	h.ID = historyId
	return nil
}

func (h *History) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetHistory)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(h.ProductID, h.Msisdn)
	if getErr := result.Scan(&h.ID, &h.ProductID, &h.Msisdn, &h.Adnet, &h.Subject, &h.CreatedAt); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}
