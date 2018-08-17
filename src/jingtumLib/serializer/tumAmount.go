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

	"github.com/yangxuebo-138/decimal"

	"jingtumLib/constant"
	jtUtils "jingtumLib/utils"
)

//TumAmount 金额结构体。
type TumAmount struct {
	Currency   string
	Issuer     string
	Value      *big.Int
	IsNative   bool
	IsNegative bool
	IsZero     bool
	IsPositive bool
	Offset     int
}

var (
	biXnsMax = big.NewInt(9000000000000000000)
	biXnsMin = big.NewInt(-9000000000000000000)

	//ConfigCurrencty 配置的货币
	ConfigCurrencty string
)

//NewTumAmount 初始化金额结构体。
func NewTumAmount() *TumAmount {
	tumAmount := new(TumAmount)
	tumAmount.IsNative = true
	return tumAmount
}

//IsPositiveM IsPositiveM
func (amount *TumAmount) IsPositiveM() bool {
	return !amount.IsZeroM() && !amount.IsNegative
}

//IsZeroM 零值判断。
func (amount *TumAmount) IsZeroM() bool {
	if amount.Value == nil {
		amount.Value = big.NewInt(0)
	}

	return (amount.Value.Cmp(big.NewInt(0)) == 0)
}

func fromJSON(inJSON interface{}) (*TumAmount, error) {
	tumAmount := new(TumAmount)
	err := tumAmount.parseJSON(inJSON)
	if err != nil {
		return nil, err
	}
	return tumAmount, nil
}

func (amount *TumAmount) parseJSON(inJSON interface{}) error {
	if jtUtils.IsNumberType(inJSON) {
		err := amount.parseSwtValue(jtUtils.NumberToString(inJSON))
		if err != nil {
			return err
		}
	} else if v, ok := inJSON.(string); ok {
		err := amount.parseSwtValue(v)
		if err != nil {
			return err
		}
	} else if jsonAmount, ok := inJSON.(constant.Amount); ok {
		if jsonAmount.Currency == constant.CFGCurrency {
			err := amount.parseSwtValue(jsonAmount.Value)
			if err != nil {
				return err
			}
		} else {
			amount.Currency = jsonAmount.Currency
			amount.IsNative = false
			if jsonAmount.Issuer == "" || !jtUtils.IsValidAddress(jsonAmount.Issuer) {
				return fmt.Errorf("Input Amount has invalid issuer info %s", jsonAmount.Issuer)
			}
			amount.Issuer = jsonAmount.Issuer
			value, err := strconv.ParseFloat(jsonAmount.Value, 64)
			if err != nil {
				return fmt.Errorf("Input JSON swt value invalid %s", jsonAmount.Value)
			}

			valueStr := fmt.Sprintf("%.16e", value)
			powStr := valueStr[strings.LastIndex(valueStr, "e")+1:]
			vpow, pintErr := strconv.ParseInt(powStr, 10, 64)
			if pintErr != nil {
				return fmt.Errorf("pow parse int value invalid %s", powStr)
			}
			offset := 15 - vpow
			value = value * math.Pow10(int(15-vpow))

			bIntV, ok := big.NewInt(0).SetString(fmt.Sprintf("%0.0f", value), 10)
			if !ok {
				return fmt.Errorf("Input JSON swt value invalid %v", value)
			}
			amount.Value = bIntV
			amount.Offset = int(-offset)
			return nil
		}
	}

	return nil
}

//TumToBytes 金额转字节
func (amount *TumAmount) TumToBytes() ([]byte, error) {
	var currencyData = make([]byte, 20)
	if len(amount.Currency) >= currencyNameLen && len(amount.Currency) <= currencyNameLen2 {
		currencyCode := amount.Currency //区分大小写
		end := 14
		cclen := len(currencyCode) - 1
		for j := cclen; j >= 0; j-- {
			currencyData[end-j] = (currencyCode[cclen-j] & 0xff)
		}
	} else if len(amount.Currency) == 40 {
		bIntV, ok := big.NewInt(0).SetString(amount.Currency, 16)
		if !ok {
			return nil, fmt.Errorf("Invalid currency code %s", amount.Currency)
		}
		bytes := bIntV.Bytes()
		return bytes, nil
	} else {
		return nil, fmt.Errorf("Incorrect currency code length")
	}

	return currencyData, nil
}

func (amount *TumAmount) parseSwtValue(jsonStr string) error {
	if !jtUtils.MatchString("^(-?)(\\d*)(\\.\\d{0,6})?$", jsonStr) {
		amount.Value = nil
		return nil
	}

	value, err := decimal.NewFromString(jsonStr)

	if err != nil {
		return fmt.Errorf("Input JSON swt value invalid %s", jsonStr)
	}

	amount.IsNative = true
	amount.Offset = 0
	vflt64, ok := value.Float64()
	if ok {
		amount.IsNegative = vflt64 < 0
	}

	if amount.IsNegative {
		value = decimal.NewFromFloat(-vflt64)
	}

	value = value.Mul(decimal.NewFromFloat(float64(1000000)))
	//bigIntStr := fmt.Sprintf("%0.0f", value)
	bIntV, ok := big.NewInt(0).SetString(value.String(), 10)
	if !ok {
		return fmt.Errorf("Input JSON swt value invalid %v", value)
	}

	amount.Value = bIntV

	if amount.Value.Cmp(biXnsMax) > 0 {
		amount.Value = nil
	}

	return nil
}

//IsValid 金额合法验证
func (amount *TumAmount) IsValid() bool {
	return amount.Value != nil
}
