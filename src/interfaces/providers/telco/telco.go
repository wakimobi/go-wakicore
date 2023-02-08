package telco

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/wakimobi/go-wakicore/src/interfaces/providers"
)

type TelcoCfg struct {
	URL,
	User,
	Pass,
	Prefix string
}

type Telco struct {
	cfg   TelcoCfg
	reqID string
}

func NewTelco(cfg TelcoCfg) *Telco {
	return &Telco{
		cfg: cfg,
	}
}

type tokenBody struct {
	Username string `json:"username"`
	Password string `json:"passwd"`
}

type mobileTerminatingBody struct {
	Token     string `json:"token"`
	Msisdn    string `json:"msisdn"`
	Sms       string `json:"sms"`
	ServiceId string `json:"serviceid"`
	Telco     string `json:"telco"`
	MtType    string `json:"mt_type"`
}

type methodResponse struct {
	StatusCode string `json:"status_code"`
	Data       string `json:"data"`
}

func (t *Telco) Token() (string, error) {
	handleErr := func(err error) error {
		return fmt.Errorf("charge > %w", err)
	}

	headers := map[string]string{
		"ContentType": "application/json",
	}

	jsonStr, _ := json.Marshal(&tokenBody{
		Username: t.cfg.User,
		Password: t.cfg.Pass,
	})

	// the body to pass
	body := bytes.NewBuffer(jsonStr)
	return "", nil
}

func (t *Telco) MobileTerminating(token string, p providers.ChargeParams) error {

	handleErr := func(err error) error {
		return fmt.Errorf("charge > %w", err)
	}
	if !strings.HasPrefix(p.Msisdn, t.cfg.Prefix) {
		return handleErr(fmt.Errorf("invalid prefix in msisdn: %q", p.Msisdn))
	}
	ctx := context.TODO()
	if err := t.runRequest(ctx, p.Msisdn, t.reqID, p.Amount); err != nil {
		return handleErr(err)
	}
	return nil

}

func (t *Telco) runRequest(ctx context.Context, msisdn string, reqID string, price int) error {
	handleErr := func(err error) error {
		return fmt.Errorf("running request: msisdn: %s, reqID: %s, price %d > %w", msisdn, reqID, price, err)
	}

	u := t.cfg.URL

	// token
	token, err := t.tf.FetchToken()
	if err != nil {
		return nil
	}

	jsonStr, _ := json.Marshal(&mobileTerminatingBody{
		Token:     token,
		Msisdn:    msisdn,
		Sms:       "",
		ServiceId: "",
		Telco:     "",
		MtType:    "",
	})

	fmt.Println("telemor telco request:", u)

	// the body to pass
	body := bytes.NewBuffer(jsonStr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return handleErr(err)
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return handleErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return handleErr(fmt.Errorf("unexpected response status: %q", resp.Status))
	}

	//8.response reading
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return handleErr(err)
	}
	defer resp.Body.Close()
}

func (m methodResponse) getMTResponseCode() (string, error) {
	return m.StatusCode, fmt.Errorf("balance not found in %#v", m)
}

func getMtType(val string) string {
	switch val {
	case "1":
		return "firstpush"
	case "2":
		return "dailypush"
	case "3":
		return "retry_dailypush"
	default:
		return "undefined"
	}
}

func getStatus(val string) string {
	switch val {
	case "1":
		return "Charging success"
	case "3321":
		return "Not Enough Credit"
	case "51":
		return "subs quota finished"
	case "52":
		return "subs doesn't have subscription"
	case "6":
		return "Quarantine"
	default:
		return "undefined"
	}
}
