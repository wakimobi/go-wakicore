package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/transactions"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func CreateTransaction(tr transactions.Transaction) (*transactions.Transaction, *errors.RestErr) {

	if err := tr.Save(); err != nil {
		return nil, err
	}
	return &tr, nil
}

func GetTransaction(productId int, msisdn string) (*transactions.Transaction, *errors.RestErr) {
	result := transactions.Transaction{ProductID: productId, Msisdn: msisdn}
	if err := result.Get(); err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateTransaction(isPartial bool, tr transactions.Transaction) (*transactions.Transaction, *errors.RestErr) {
	current, err := GetTransaction(tr.ProductID, tr.Msisdn)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if tr.Status != "" {
			current.Status = tr.Status
		}

		if tr.Subject != "" {
			current.Subject = tr.Subject
		}

	} else {
		current.Status = tr.Status
		current.Subject = tr.Subject
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}
