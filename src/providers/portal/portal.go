package portal

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/sirupsen/logrus"
)

type Portal struct {
	cfg          *config.Secret
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	pin          int
}

func NewPortal(
	cfg *config.Secret,
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	pin int,
) *Portal {
	return &Portal{
		cfg:          cfg,
		logger:       logger,
		subscription: subscription,
		service:      service,
		pin:          pin,
	}
}

func (p *Portal) Subscription() ([]byte, error) {
	l := p.logger.Init("notif", true)

	payload := url.Values{}
	payload.Add("msisdn", p.subscription.Msisdn)
	payload.Add("event", "reg")
	payload.Add("password", strconv.Itoa(p.pin))
	payload.Add("package", "1")

	req, err := http.NewRequest("GET", p.service.UrlNotifSub+"?"+payload.Encode(), nil)
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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlNotifSub + "?" + payload.Encode()}).Info("SUBSCRIPTION")

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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "response": string(body)}).Info("SUBSCRIPTION")

	return body, nil
}

func (p *Portal) Unsubscription() ([]byte, error) {
	l := p.logger.Init("notif", true)

	payload := url.Values{}
	payload.Add("msisdn", p.subscription.Msisdn)
	payload.Add("event", "unreg")

	req, err := http.NewRequest("GET", p.service.UrlNotifUnsub+"?"+payload.Encode(), nil)
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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlNotifUnsub + "?" + payload.Encode()}).Info("UNSUBSCRIPTION")

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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "response": string(body)}).Info("UNSUBSCRIPTION")

	return body, nil
}

func (p *Portal) Renewal() ([]byte, error) {
	l := p.logger.Init("notif", true)

	payload := url.Values{}
	payload.Add("msisdn", p.subscription.Msisdn)
	payload.Add("event", "renewal")
	payload.Add("password", strconv.Itoa(p.pin))
	payload.Add("package", "1")

	req, err := http.NewRequest("GET", p.service.UrlNotifRenewal+"?"+payload.Encode(), nil)

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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlNotifRenewal + "?" + payload.Encode()}).Info("RENEWAL")

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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "response": string(body)}).Info("RENEWAL")

	return body, nil
}

func (p *Portal) SubscriptionRetry() ([]byte, error) {
	l := p.logger.Init("notif", true)

	payload := url.Values{}
	payload.Add("msisdn", p.subscription.GetMsisdn())
	payload.Add("event", "reg")
	payload.Add("password", strconv.Itoa(p.subscription.GetLatestPIN()))
	payload.Add("package", "1")

	req, err := http.NewRequest("GET", p.service.UrlNotifSub+"?"+payload.Encode(), nil)
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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlNotifSub + "?" + payload.Encode()}).Info("SUBSCRIPTION")

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
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "response": string(body)}).Info("SUBSCRIPTION")

	return body, nil
}

func (p *Portal) Callback() string {
	callbackUrl := p.service.UrlPortal + "?msisdn=" + p.subscription.Msisdn
	return callbackUrl
}
