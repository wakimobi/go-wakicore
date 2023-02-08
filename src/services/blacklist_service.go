package services

import (
	"github.com/wakimobi/go-wakicore/src/domain/blacklists"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func CountBlacklist(msisdn string) (int, *errors.RestErr) {
	result := blacklists.Blacklist{
		Msisdn: msisdn,
	}
	count, err := result.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetBlacklist(msisdn string) (*blacklists.Blacklist, *errors.RestErr) {
	result := blacklists.Blacklist{
		Msisdn: msisdn,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
