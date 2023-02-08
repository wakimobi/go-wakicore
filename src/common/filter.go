package common

import "strings"

const (
	REG       = "REG"
	UNREG     = "UNREG"
	PRODUCT_1 = "WECARE"
)

func FilterReg(message string) bool {
	index := strings.Split(strings.ToUpper(message), " ")
	if index[0] == REG &&
		(strings.Contains(strings.ToUpper(message), REG+" "+PRODUCT_1)) {
		return true
	}
	return false
}

func FilterUnreg(message string) bool {
	index := strings.Split(strings.ToUpper(message), " ")
	if index[0] == UNREG &&
		(strings.Contains(strings.ToUpper(message), UNREG+" "+PRODUCT_1)) {
		return true
	}
	return false
}

func GetSubKeyword(message string) string {

	if strings.Contains(message, REG+" "+PRODUCT_1) {
		i := len(REG + " " + PRODUCT_1)
		return message[i:]
	}

	return ""
}
