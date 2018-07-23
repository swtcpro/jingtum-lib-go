/**
 *
 * 文件功能介绍
 *
 * @FileName: tumAmount.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-16 10:44:32
 * @UpdateTime: 2018-07-16 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	jtbLib "jingtumLib/jingtumBaseLib"
)

type Amount struct {
	Currency string
	Issuer   string
	Value    string
}

type TumAmount struct {
	Currency   string
	Issuer     string
	Value      *big.Int
	IsValid    bool
	IsNative   bool
	IsNegative bool
	IsZero     bool
	IsPositive bool
	Offset     int
}

var (
	bi_xns_max      = big.NewInt(9000000000000000000)
	bi_xns_min      = big.NewInt(-9000000000000000000)
	ConfigCurrencty string
)

func NewTumAmount() *TumAmount {
	tumAmount := new(TumAmount)
	tumAmount.IsNative = true
	tumAmount.IsValid = tumAmount.Value != nil
	return tumAmount
}

func (amount *TumAmount) IsPositiveM() bool {
	return !amount.IsZeroM() && !amount.IsNegative
}

func (amount *TumAmount) IsZeroM() bool {
	if amount.Value == nil {
		amount.Value = big.NewInt(0)
	}

	return (amount.Value.Cmp(big.NewInt(0)) == 0)
}

func FromJsonTumA(json interface{}) *TumAmount {
	amount := new(TumAmount)
	amount.parseJson(json)
	return amount
}

func (amount *TumAmount) parseJson(json interface{}) {
	if IsNumberType(json) {
		amount.parseSwtValue(NumberToString(json))
	} else if v, ok := json.(string); ok {
		amount.parseSwtValue(v)
	} else if jsonAmount, ok := json.(Amount); ok {
		if jsonAmount.Currency == ConfigCurrencty {
			amount.parseSwtValue(jsonAmount.Value)
		} else {
			amount.Currency = jsonAmount.Currency
			amount.IsNative = false
			if jsonAmount.Issuer == "" || !jtbLib.IsValidAddress(jsonAmount.Issuer) {
				panic(fmt.Sprintf("Input Amount has invalid issuer info %v", jsonAmount.Issuer))
			} else {
				amount.Issuer = jsonAmount.Issuer
				value, err := strconv.ParseFloat(jsonAmount.Value, 64)
				if err != nil {
					panic(fmt.Sprintf("Input JSON swt value invalid %v", jsonAmount.Value))
				}

				valueStr := fmt.Sprintf("%.16e", value)
				powStr := valueStr[strings.LastIndex(valueStr, "e")+1:]
				vpow, pintErr := strconv.ParseInt(powStr, 10, 64)
				if pintErr != nil {
					panic(fmt.Sprintf("pow parse int value invalid %v", powStr))
				}
				offset := 15 - vpow
				value = value * math.Pow10(int(15-vpow))

				bIntV, ok := big.NewInt(0).SetString(fmt.Sprintf("%0.0f", value), 10)
				if !ok {
					panic(fmt.Sprintf("Input JSON swt value invalid %v", value))
				}
				//factor := fmt.Sprintf("%0.0f",math.Pow10(int(15-pow)))
				amount.Value = bIntV
				amount.Offset = int(-offset)
			}
		}
	}
}

func (amount *TumAmount) TumToBytes() []byte {
	var currencyData []byte
	if len([]rune(amount.Currency)) >= CURRENCY_NAME_LEN && len([]rune(amount.Currency)) <= CURRENCY_NAME_LEN2 {
		currencyCode := amount.Currency //区分大小写
		end := 14
		cclen := len([]rune(currencyCode))
		for j := cclen; j >= 0; j-- {
			currencyData[end-j] = (currencyCode[cclen-j] & 0xff)
		}
	} else if len([]rune(amount.Currency)) == 40 {
		bIntV, ok := big.NewInt(0).SetString(amount.Currency, 16)
		if !ok {
			panic(fmt.Sprintf("Invalid currency code %v", amount.Currency))
		}
		bytes := bIntV.Bytes()
		return bytes
	} else {
		panic("Incorrect currency code length")
	}

	return currencyData
}

func (amount *TumAmount) parseSwtValue(jsonStr string) {
	if !MatchString("^(-?)(\\d*)(\\.\\d{0,6})?$", jsonStr) {
		amount.Value = nil
		return
	}

	value, err := strconv.ParseFloat(jsonStr, 64)
	if err != nil {
		panic(fmt.Sprintf("Input JSON swt value invalid %v", jsonStr))
	}

	amount.IsNative = true
	amount.Offset = 0
	amount.IsNegative = value < 0
	if amount.IsNegative {
		value = -value
	}

	value = value * 1000000
	bigIntStr := fmt.Sprintf("%0.0f", value)
	bIntV, ok := big.NewInt(0).SetString(bigIntStr, 10)
	if !ok {
		panic(fmt.Sprintf("Input JSON swt value invalid %v", bigIntStr))
	}

	amount.Value = bIntV

	if amount.Value.Cmp(bi_xns_max) > 0 {
		amount.Value = nil
	}
}
