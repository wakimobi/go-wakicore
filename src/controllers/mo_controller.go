package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wakimobi/go-wakicore/src/common"
	"github.com/wakimobi/go-wakicore/src/domain/subscriptions"
	"github.com/wakimobi/go-wakicore/src/services"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
)

func MOHandler(c *gin.Context) {
	var req common.MORequest
	if err := c.ShouldBindQuery(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid query param")
		c.XML(restErr.Status, restErr)
		return
	}

	/**
	 * check product
	 */
	product, err := services.GetProduct("")
	if err != nil {
		restErr := errors.NewInternalServerError("product not found")
		c.XML(restErr.Status, restErr)
		return
	}

	/**
	 * check blacklist
	 */
	blacklist, err := services.CountBlacklist(req.Msisdn)
	if err != nil {
		restErr := errors.NewInternalServerError("msisdn has blacklist")
		c.XML(restErr.Status, restErr)
		return
	}

	if blacklist > 0 {
		restErr := errors.NewForbiddenError("msisdn has blacklist")
		c.XML(restErr.Status, restErr)
		return
	}

	/**
	 * check active sub
	 */
	sub, err := services.CountSubscription(product.ID, req.Msisdn)
	if err != nil {
		restErr := errors.NewInternalServerError("check active sub")
		c.XML(restErr.Status, restErr)
		return
	}

	if sub > 0 {
		restErr := errors.NewForbiddenError("msisdn already sub")
		c.XML(restErr.Status, restErr)
		return
	}

	c.XML(200, common.ResponseXML{Status: "OK"})
}

func InsertHandler(c *gin.Context) {
	var req subscriptions.Subscription
	if err := c.ShouldBindJSON(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid body param")
		c.XML(restErr.Status, restErr)
		return
	}

	_, err := services.CreateSubscription(req)
	if err != nil {
		restErr := errors.NewInternalServerError("create subscription")
		c.XML(restErr.Status, restErr)
		return
	}

	c.XML(http.StatusCreated, common.ResponseXML{Status: "OK"})
}
