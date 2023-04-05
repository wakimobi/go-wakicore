package entity

import (
	"strings"
)

const (
	MO_REG        = "REG"
	MO_UNREG      = "UNREG"
	VALID_PREFIX  = "628"
	SERVICE_GOALY = "GOALY"
)

type ReqMOParams struct {
	SMS       string `validate:"required" query:"sms" json:"sms"`
	Adn       string `query:"adn" json:"adn"`
	Msisdn    string `validate:"required" query:"msisdn" json:"msisdn"`
	Channel   string `query:"channel" json:"channel"`
	TrxId     string `query:"trx_id" json:"trx_id"`
	Number    string `query:"http_segment_number" json:"http_segment_number"`
	Count     string `query:"http_segment_count" json:"http_segment_count"`
	IpAddress string `query:"ip_address" json:"ip_address"`
}

type ReqMOBody struct {
	MessageID struct {
		Sms struct {
			Retry struct {
				Count       string `json:"count" xml:"count"`
				Max         string `json:"max" xml:"max"`
				Destination struct {
					Address struct {
						Unknown struct {
							Cnpi string `json:"cnpi" xml:"cnpi"`
						} `json:"unknown" xml:"unknown"`
					} `json:"address" xml:"address"`
				} `json:"destination" xml:"destination"`
				Source struct {
					Address struct {
						Number struct {
							Type string `json:"type" xml:"type"`
						} `json:"number" xml:"number"`
					} `json:"address" xml:"address"`
				} `json:"source" xml:"source"`
				Ud struct {
					Type string `json:"type" xml:"type"`
				} `json:"ud" xml:"ud"`
				Param struct {
					Name  string `json:"name" xml:"name"`
					Value string `json:"value" xml:"value"`
				} `json:"param" xml:"param"`
			} `json:"retry"`
		} `json:"sms" xml:"sms"`
	} `json:"message" xml:"message"`
}

type ReqDRParams struct {
	SMS       string `validate:"required" query:"sms" json:"sms"`
	CpId      string `query:"cpid" json:"cpid"`
	Pwd       string `query:"pwd" json:"pwd"`
	Msisdn    string `validate:"required" query:"msisdn" json:"msisdn"`
	TrxId     string `validate:"required" query:"trx_id" json:"trx_id"`
	Sid       string `query:"sid" json:"sid"`
	Sender    string `query:"sender" json:"sender"`
	Tid       string `query:"tid" json:"tid"`
	IpAddress string `query:"ip_address" json:"ip_address"`
}

type ReqMTParams struct {
	SMS    string `url:"sms,omitempty" query:"sms"`
	CpId   string `url:"cpid,omitempty" query:"cpid"`
	Pwd    string `url:"pwd,omitempty" query:"pwd"`
	Msisdn string `url:"msisdn,omitempty" query:"msisdn"`
	TrxId  string `url:"trx_id,omitempty" query:"trx_id"`
	Sid    string `url:"sid,omitempty" query:"sid"`
	Sender string `url:"sender,omitempty" query:"sender"`
	Tid    string `url:"tid,omitempty" query:"tid"`
}

type ReqMTBody struct {
	Message struct {
		Sms struct {
			Type        string `xml:"type,attr"`
			Destination struct {
				Address struct {
					Number string `xml:"number"`
				} `xml:"address"`
			} `xml:"destination"`
			Source struct {
				Address struct {
					Number string `xml:"number"`
				} `xml:"address"`
			} `xml:"source"`
			Ud    string           `xml:"ud"`
			Param []ReqMTBodyParam `xml:"param"`
		} `xml:"sms"`
	} `xml:"message"`
}

type ReqMTBodyParam struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type ReqOptInParam struct {
	Service string `json:"service"`
	Adnet   string `json:"adnet"`
	PubId   string `json:"pub_id"`
	AffSub  string `json:"aff_sub"`
}

type ResponseMO struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Message    string `json:"message" xml:"message"`
}

type ResponseDR struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Message    string `json:"message" xml:"message"`
}

type AveragePerUser struct {
	Name       string `json:"name"`
	Adnet      string `json:"adnet"`
	Subs       string `json:"subs"`
	SubsActive string `json:"subs_active"`
	Revenue    string `json:"revenue"`
}

type ErrorResponse struct {
	FailedField string `json:"failed_field" xml:"failed_field"`
	Tag         string `json:"tag" xml:"tag"`
	Value       string `json:"value" xml:"value"`
}

func NewReqMOParams(sms, adn, msisdn, channel, ip_address string) *ReqMOParams {
	return &ReqMOParams{
		SMS:       sms,
		Adn:       adn,
		Msisdn:    msisdn,
		Channel:   channel,
		IpAddress: ip_address,
	}
}

func NewReqDRParams(sms, TrxId, msisdn, ip_address string) *ReqDRParams {
	return &ReqDRParams{
		SMS:       sms,
		TrxId:     TrxId,
		Msisdn:    msisdn,
		IpAddress: ip_address,
	}
}

func (s *ReqMOParams) GetSMS() string {
	return s.SMS
}

func (s *ReqMOParams) SetSMS(sms string) {
	s.SMS = strings.ToUpper(sms)
}

func (s *ReqMOParams) GetAdn() string {
	return s.Adn
}

func (s *ReqMOParams) GetMsisdn() string {
	return s.Msisdn
}

func (s *ReqMOParams) GetChannel() string {
	return s.Channel
}

func (s *ReqMOParams) GetIpAddress() string {
	return s.IpAddress
}

func (s *ReqMOParams) IsInValidPrefix() bool {
	return !strings.HasPrefix(s.Msisdn, VALID_PREFIX)
}

func (s *ReqMOParams) IsREG() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_REG && (strings.Contains(message, MO_REG+" "+SERVICE_GOALY)) {
		return true
	}
	return false
}

func (s *ReqMOParams) IsUNREG() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_UNREG && (strings.Contains(message, MO_UNREG+" "+SERVICE_GOALY)) {
		return true
	}
	return false
}

func (s *ReqMOParams) GetKeyword() string {
	return strings.ToUpper(s.SMS)
}

func (s *ReqMOParams) GetSubKeyword() string {
	message := strings.ToUpper(s.SMS)
	i := strings.Index(message, " ")
	if i > -1 {
		keyword := s.SMS[i+1:]
		if strings.Contains(strings.ToUpper(keyword), SERVICE_GOALY) {
			i := len(SERVICE_GOALY)
			return strings.ToUpper(keyword[:i])
		}
		return keyword
	}
	return message
}

func (s *ReqMOParams) GetAdnet() string {
	message := strings.ToUpper(s.SMS)
	keyword := MO_REG + " " + SERVICE_GOALY
	if strings.Contains(message, keyword) {
		i := len(keyword)
		return message[i:]
	}
	return ""
}

func (s *ReqDRParams) GetTrxId() string {
	return s.TrxId
}

func (s *ReqDRParams) GetMsisdn() string {
	return s.Msisdn
}
