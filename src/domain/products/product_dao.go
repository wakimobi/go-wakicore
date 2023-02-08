package products

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryGetProduct  = "SELECT id, code, name, auth_user, auth_pass, price, renewal_day, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback FROM products WHERE id = $1 LIMIT 1"
	queryFindProduct = "SELECT id, code, name, auth_user, auth_pass, price, renewal_day, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback FROM products WHERE name = $2 LIMIT 1"
)

func (p *Product) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetProduct)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(p.ID)
	if getErr := result.Scan(&p.ID, &p.Code, &p.Name, &p.AuthPass, &p.AuthPass, &p.Price, &p.RenewalDay, &p.UrlNotifSub, &p.UrlNotifUnsub, &p.UrlNotifRenewal, &p.UrlPostback); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}

func (p *Product) Search() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryFindProduct)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(p.Name)
	if getErr := result.Scan(&p.ID, &p.Code, &p.Name, &p.AuthPass, &p.AuthPass, &p.Price, &p.RenewalDay, &p.UrlNotifSub, &p.UrlNotifUnsub, &p.UrlNotifRenewal, &p.UrlPostback); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}
