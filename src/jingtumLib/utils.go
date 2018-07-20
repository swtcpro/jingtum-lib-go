/**
 *
 * 文件功能介绍
 *
 * @FileName: utils.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-10 10:44:32
 * @UpdateTime: 2018-07-10 10:44:54
 * Copyright@2018 版权所有
 */

package jingtumLib

import (
     "strconv"
     "regexp"
     "encoding/hex"
     "errors"

     jtbLib "jingtumLib/jingtumBaseLib"
)

var (

	LEDGER_STATES = map[string]string{"current": "current", "closed": "closed", "validated": "validated"}
    CURRENCY_RE   = "^([a-zA-Z0-9]{3,6}|[A-F0-9]{40})$"
)

func number(s string) bool {
    
    return number("^(-?)(\\d*)(\\.\\d{0,6})?$", s)
}

func IsNumber(s string) bool {
    return number(s)
}

func IsNumberType(obj interface{}) bool {
    switch obj.(type) {
        case float64,float32,int,int8,int32,int64,byte,uint32,uint64:
            return true
        default:
            return false
    }
}

func NumberToString(obj interface{}) string {
    switch v := obj.(type) {
        case float32:
            return strconv.FormatFloat(float64(v), 'f', -1, 64)
        case float64:
            return strconv.FormatFloat(v, 'f', -1, 64)
        case int:
            return strconv.Itoa(v)
        case int8:
            return strconv.FormatInt(int64(v), 10)
        case int32:
            return strconv.FormatInt(int64(v), 10)
        case int64:
            return strconv.FormatInt(v, 10)
        case uint8:
            return strconv.FormatUint(uint64(v),10)
        case uint32:
            return strconv.FormatUint(uint64(v),10)
        case uint64:
            return strconv.FormatUint(uint64(v),10)
        default:
            return ""
    }
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
    if (amount.currency != JTConfig.Read("Config","currency") && !jtbLib.IsValidAddress(amount.issuer)) {
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

func IsValidCurrency(currency string) bool {
    if (currency == "") {
        return false;
    }

    return matchString(CURRENCY_RE, currency);
}

func isValidAddress(issuer string) bool {
    return jtbLib.IsValidAddress(issuer)
}