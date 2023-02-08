package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/histories"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func CreateHistory(his histories.History) (*histories.History, *errors.RestErr) {
	if err := his.Validate(); err != nil {
		return nil, err
	}
	if err := his.Save(); err != nil {
		return nil, err
	}
	return &his, nil
}

func GetHistory(productId int, msisdn string) (*histories.History, *errors.RestErr) {
	result := histories.History{
		ProductID: productId,
		Msisdn:    msisdn,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
