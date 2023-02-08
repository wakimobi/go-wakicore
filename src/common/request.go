package common

/**
 * {"adn":"99345","sms":"UNREG WECARE","trx_id":"4b52ebedb2c91e941601bf85a02c5bb9","telco":"HU","msisdn":"62895335990043"}
 */
type MORequest struct {
	Adn    string `form:"adn" query:"adn" json:"adn"`
	Sms    string `form:"sms" query:"sms" json:"sms"`
	TrxId  string `form:"trx_id" query:"trx_id" json:"trx_id"`
	Telco  string `form:"telco" query:"telco" json:"telco"`
	Msisdn string `form:"msisdn" query:"msisdn" json:"msisdn"`
}

/**
 * {"trx_id":"2023012300011147065918954","status":"FAIL","error_code":"5:8","telco":"HU","keyword":"WECARE","msisdn":"6289669756342","mt_type":"2"}
 */
type DRRequest struct {
	TrxId     string `form:"trx_id" query:"trx_id" json:"trx_id"`
	Status    string `form:"status" query:"status" json:"status"`
	ErrorCode string `form:"error_code" query:"error_code" json:"error_code"`
	Telco     string `form:"telco" query:"telco" json:"telco"`
	Keyword   string `form:"keyword" query:"keyword" json:"keyword"`
	Msisdn    string `form:"msisdn" query:"msisdn" json:"msisdn"`
	MtType    string `form:"mt_type" query:"mt_type" json:"mt_type"`
}
