package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wakimobi/go-wakicore/src/common"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/logger_utils"
)

func DRHandler(c *gin.Context) {
	// {"trx_id":"2023012300011147065918954","status":"FAIL","error_code":"5:8","telco":"HU","keyword":"WECARE","msisdn":"6289669756342","mt_type":"2"}

	loggerDr := logger_utils.MakeLogger("dr", true)

	var req common.MORequest
	if err := c.ShouldBindQuery(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid query param")
		c.XML(restErr.Status, restErr)
		return
	}

	loggerDr.WithFields(logrus.Fields{})
	c.XML(200, common.ResponseXML{Status: "OK"})

}
