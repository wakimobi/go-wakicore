package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/contents"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func GetContent(productId int, name string) (*contents.Content, *errors.RestErr) {
	result := contents.Content{
		ProductID: productId,
		Name:      name,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
