package postback

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/sirupsen/logrus"
)

type Postback struct {
	cfg          *config.Secret
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
}

func NewPostback(
	cfg *config.Secret,
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
) *Postback {
	return &Postback{
		cfg:          cfg,
		logger:       logger,
		subscription: subscription,
		service:      service,
	}
}

func (p *Postback) Send() ([]byte, error) {
	l := p.logger.Init("pb", true)

	payload := url.Values{}
	payload.Add("partner", "passtisel")
	payload.Add("px", "")
	payload.Add("serv_id", p.service.Code)
	payload.Add("msisdn", p.subscription.Msisdn)
	payload.Add("trxid", "")

	req, err := http.NewRequest("GET", p.service.UrlPostback+"?"+payload.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlPostback + "?" + payload.Encode()}).Info("POSTBACK")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "response": string(body)}).Info("POSTBACK")

	return body, nil
}
