/**
 *
 * 文件功能介绍
 *
 * @FileName: utils.go
 * @Auther : 13851485286
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumLib

import (
     "strconv"
     "regexp"
     "encoding/hex"
)

var (

	LEDGER_STATES = map[string]string{"current": "current", "closed": "closed", "validated": "validated"}
    CURRENCY_RE   = "^([a-zA-Z0-9]{3,6}|[A-F0-9]{40})$"
)

func number(s string) bool {
    
    return MatchString("^[0-9]+\\.?[0-9]*$", s)
}

func matchString(patter string, str string) bool {
    match, _ := regexp.MatchString(patter, str)

    return match
}

func hexToString(hexStr string) (string, error) {
    bytes, err := hex.DecodeString(hexStr)
    return string(bytes), err
}

func stringToHex(str string) string {
    return hex.EncodeToString([]byte(str))
}

func isValidAmount(amount Amount) bool {

    if ((amount.Value == "")) {
        return false;
    }

    // check amount currency
    if ((amount.currency == "") || !isValidCurrency(amount.currency)) {
        return false;
    }

    // native currency issuer is empty
    if (amount.currency == JTConfig.Read("Config","currency") && amount.issuer != "") {
        return false;
    }

    // non native currency issuer is not allowed to be empty
    if (amount.currency != JTConfig.Read("Config","currency") && !isValidAddress(amount.issuer)) {
        return false;
    }

    return true;
}

func isValidCurrency(currency string) bool {
    if (currency == "") {
        return false;
    }

    return matchString(CURRENCY_RE, currency);
}

func isValidAddress(issuer string) bool {
    return false
}