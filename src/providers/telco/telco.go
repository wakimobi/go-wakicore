package telco

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/idprm/go-pass-tsel/src/config"
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/logger"
	"github.com/idprm/go-pass-tsel/src/utils/hash_utils"
	"github.com/sirupsen/logrus"
)

type Telco struct {
	cfg          *config.Secret
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	content      *entity.Content
}

func NewTelco(
	cfg *config.Secret,
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	content *entity.Content,
) *Telco {
	return &Telco{
		cfg:          cfg,
		logger:       logger,
		subscription: subscription,
		service:      service,
		content:      content,
	}
}

type ITelco interface {
	Token() ([]byte, error)
	WebOptInOTP() (string, error)
	WebOptInUSSD() (string, error)
	WebOptInCaptcha() (string, error)
	SMSbyParam() ([]byte, error)
}

func (t *Telco) Token() ([]byte, error) {
	l := t.logger.Init("mt", true)

	req, err := http.NewRequest("GET", t.cfg.Telco.UrlKey+"/scrt/1/generate.php", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cp_name", t.cfg.Telco.CpName)
	q.Add("pwd", t.cfg.Telco.Pwd)
	q.Add("programid", t.service.GetProgramId())
	q.Add("sid", t.service.GetProgramId())

	req.URL.RawQuery = q.Encode()

	now := time.Now()
	timeStamp := strconv.Itoa(int(now.Unix()))
	strData := t.cfg.Telco.Key + t.cfg.Telco.Secret + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Set("api_key", t.cfg.Telco.Key)
	req.Header.Set("x-signature", signature)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{"request": t.cfg.Telco.UrlKey + "/scrt/1/generate.php?" + q.Encode()}).Info("MT_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{"response": string(body)}).Info("MT_TOKEN")

	return body, nil
}

func (t *Telco) WebOptInOTP() (string, error) {
	l := t.logger.Init("mt", true)

	token, err := t.Token()
	if err != nil {
		return "", err
	}
	l.WithFields(logrus.Fields{"redirect": t.cfg.Telco.UrlAuth + "/transaksi/tauthwco?token=" + string(token)}).Info("MT_OPTIN")
	return t.cfg.Telco.UrlAuth + "/transaksi/tauthwco?token=" + string(token), nil
}

func (t *Telco) WebOptInUSSD() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return t.cfg.Telco.UrlAuth + "/transaksi/konfirmasi/ussd?token=" + string(token), nil
}

func (t *Telco) WebOptInCaptcha() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return t.cfg.Telco.UrlAuth + "/transaksi/captchawco?token=" + string(token), nil
}

func (t *Telco) SMSbyParam() ([]byte, error) {
	l := t.logger.Init("mt", true)

	req, err := http.NewRequest(http.MethodGet, t.cfg.Telco.UrlKey+"/scrt/cp/submitSM.jsp", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cpid", t.cfg.Telco.CpId)
	q.Add("sender", t.cfg.Telco.Sender)
	q.Add("sms", t.content.GetValue())
	q.Add("pwd", t.cfg.Telco.Pwd)
	q.Add("msisdn", t.subscription.GetMsisdn())
	q.Add("sid", t.service.GetSid())
	q.Add("tid", t.content.GetTid())

	req.URL.RawQuery = q.Encode()

	now := time.Now()
	timeStamp := strconv.Itoa(int(now.Unix()))
	strData := t.cfg.Telco.Key + t.cfg.Telco.Secret + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Set("api_key", t.cfg.Telco.Key)
	req.Header.Set("x-signature", signature)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{"msisdn": t.subscription.GetMsisdn(), "request": t.cfg.Telco.UrlKey + "/scrt/cp/submitSM.jsp?" + q.Encode()}).Info("MT_SMS")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{"msisdn": t.subscription.GetMsisdn(), "response": string(body)}).Info("MT_SMS")

	return body, nil
}
