package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/adnets"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func GetAdnet(name string) (*adnets.Adnet, *errors.RestErr) {
	result := adnets.Adnet{
		Name: name,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
