package adnets

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryGetAdnet = "SELECT id, name, value FROM adnets WHERE name = $1 LIMIT 1"
)

func (a *Adnet) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetAdnet)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(a.Name)
	if getErr := result.Scan(&a.ID, &a.Name, &a.Value); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}
