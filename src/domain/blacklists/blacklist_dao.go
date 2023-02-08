package blacklists

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryGetBlacklist = "SELECT id, msisdn, created_at FROM blacklists WHERE msisdn = $1 LIMIT 1"
	queryCountAdnet   = "SELECT COUNT(*) FROM blacklists WHERE msisdn = $1"
)

func (b *Blacklist) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetBlacklist)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(b.Msisdn)
	if getErr := result.Scan(&b.ID, &b.Msisdn, &b.CreatedAt); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}

func (b *Blacklist) Count() (int, *errors.RestErr) {
	var count int
	stmt, err := db.Client.Prepare(queryCountAdnet)
	if err != nil {
		return 0, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(b.Msisdn)
	if getErr := result.Scan(&count); getErr != nil {
		return 0, pgsql_utils.ParseErrorSkip(getErr)
	}
	return count, nil
}
