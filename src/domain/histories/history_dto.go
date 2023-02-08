package histories

import (
	"time"

	"github.com/wakimobi/go-wakicore/src/domain/products"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

type History struct {
	ID        int64 `json:"id"`
	ProductID int   `json:"product_id"`
	Product   *products.Product
	Msisdn    string    `json:"msisdn"`
	Subject   string    `json:"subject"`
	Adnet     string    `json:"adnet"`
	PubID     string    `json:"pub_id"`
	AffSub    string    `json:"aff_sub"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *History) Validate() *errors.RestErr {
	if h.ProductID == 0 {
		return errors.NewBadRequestError("invalid product_id")
	}

	if h.Msisdn == "" {
		return errors.NewBadRequestError("invalid msisdn")
	}

	return nil
}
