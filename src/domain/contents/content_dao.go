package contents

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryGetContent = "SELECT name, value FROM contents WHERE product_id = $1 AND name = $2 LIMIT 1"
)

func (c *Content) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetContent)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(c.ProductID, c.Name)
	if getErr := result.Scan(&c.Name, &c.Value); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}
