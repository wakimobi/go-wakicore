package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/products"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func GetProduct(code string) (*products.Product, *errors.RestErr) {
	result := products.Product{
		Code: code,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
