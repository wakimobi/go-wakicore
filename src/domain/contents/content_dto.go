package contents

import "github.com/wakimobi/go-wakicore/src/domain/products"

type Content struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
	Product   *products.Product
	Name      string `json:"name"`
	Value     string `json:"value"`
}
