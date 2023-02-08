package pgsql_utils

import (
	"log"
	"strings"

	pq "github.com/lib/pq"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*pq.Error)
	log.Println(sqlErr)
	log.Println(ok)

	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("No record matching given id")
		}
		return errors.NewInternalServerError("Error parsing database response")
	}

	switch sqlErr.Code.Name() {
	case "02001":
		return errors.NewBadRequestError("Invalid data")
	}

	return errors.NewInternalServerError("Error processing request")
}

func ParseErrorSkip(err error) *errors.RestErr {
	sqlErr, _ := err.(*pq.Error)

	switch sqlErr.Code.Name() {
	case "02001":
		return errors.NewBadRequestError("Invalid data")
	}

	return errors.NewInternalServerError("Error processing request")
}
