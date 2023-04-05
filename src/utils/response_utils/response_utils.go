package response_utils

import "strings"

func ParseStatusCode(code string) string {
	switch code {
	case "0:1":
		return "Default error code"
	case "0:2":
		return "MT rejected due to storage partition is full"
	case "1":
		return "Success"
	case "2":
		return "Authentication failed (binding failed)"
	case "3:":
		return "Charging failed"
	case "3:101":
		return "Charging timeout"
	case "3:105":
		return "Invalid MSISDN (recipient)"
	case "3:3:21":
		return "Not enough credit"
	case "4:1":
		return "Invalid shortcode (sender)"
	case "4:2:":
		return "Mandatory parameter is missing"
	case "4:3":
		return "MT rejected due to long message restriction"
	case "4:4:1":
		return "Multiple tariff is not allowed, but “tid” parameter is provided by CP"
	case "4:4:2":
		return "The provided “tid” by CP is not allowed"
	case "5:997":
		return "Invalid trx_id"
	case "5:1":
		return "MT rejected due to subscription quota is finished"
	case "5:2":
		return "MT rejected due to subscriber doesn't have this subscription"
	case "5:3":
		return "MT rejected due to subscription is disabled"
	case "5:4":
		return "Throttling error"
	case "6":
		return "MT rejected due to quarantine"
	case "7":
		return "Error XML"
	default:
		return ""
	}
}

func ParseStatus(code string) bool {
	return strings.HasPrefix(code, "1")
}

func ParseChannel(sms string) string {
	if strings.Contains(strings.ToUpper(sms), "TOKEN=") {
		return "WAP"
	}
	return "SMS"
}
