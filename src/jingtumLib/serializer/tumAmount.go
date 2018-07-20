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
    "math/big"

    "jingtumLib"
    jtbLib "jingtumLib/jingtumBaseLib"
)

type TumAmount struct {
    Currency        string
    Issuer          string
    Value           *big.Int
    IsValid         bool
    IsNative        bool = true
    IsNegative      bool
    IsZero          bool
    IsPositive      bool
    Offset          int
}

var (
    bi_xns_max = big.NewInt(9000000000000000000)
    bi_xns_min = big.NewInt(-9000000000000000000)
    CURRENCY_NAME_LEN = 3
    CURRENCY_NAME_LEN2 = 6

func (amount *TumAmount) IsValid() bool {
    return amount.Value != nil
}

func (amount *TumAmount) IsNative() bool {
    return amount.IsNative
}

func (amount *TumAmount) IsNegative() bool {
    return amount.IsNegative
}

func (amount *TumAmount) IsPositive() bool {
    return !amount.IsZero() && !amount.IsNegative()
}

func (amount *TumAmount) IsZero() bool {
    if amount.Value == nil {
        amount.Value = big.NewInt(0)
    }

    return (amount.Value.Cmp(big.NewInt(0)) == 0)
}

func (amount *TumAmount) Issuer() string {
    return amount.Issuer
}

func FromJson(json interface{}) *TumAmount {
    amount := new(TumAmount)
    amount.parseJson(json)
    return amount
}

func (amount *TumAmount) parseJson(json interface{}) {
    if IsNumberType(json) {
        amount.parseSwtValue(NumberToString(json))
    } else if v,ok:= json.(string); ok {
        amount.parseSwtValue(json)
    } else if jsonAmount,ok:= json.(Amount); ok {
        if jsonAmount.Currency == JTConfig.Read("Config","currency") {
            amount.parseSwtValue(jsonAmount.Value)
        } else {
            amount.Currency = jsonAmount.Currency;
            amount.IsNative = false
            if jsonAmount.Issuer == "" || !jtbLib.IsValidAddress(jsonAmount.Issuer) {
                panic(fmt.Sprintf("Input Amount has invalid issuer info %v", jsonAmount.Issuer))
            } else {
                amount.Issuer = jsonAmount.Issuer
                value,err := strconv.ParseFloat(jsonAmount.Value, 64)
                if err != nil {
                    panic(fmt.Sprintf("Input JSON swt value invalid %v", jsonStr))
                }

                valueStr := fmt.Sprintf("%.16e", value)
                powStr := valueStr[strings.LastIndex(valueStr,"e")+1:]
                vpow,pintErr := strconv.ParseInt(powStr,10,64)
                if pintErr ！= nil {
                    panic(fmt.Sprintf("pow parse int value invalid %v", powStr))
                }
                offset := 15 - vpow
                value = value * math.Pow10(int(15-pow))

                bIntV,ok := big.NewInt(0).SetString(fmt.Sprintf("%0.0f",value), 10)
                if !ok {
                    panic(fmt.Sprintf("Input JSON swt value invalid %v", bigIntStr))
                }
                //factor := fmt.Sprintf("%0.0f",math.Pow10(int(15-pow)))
                amount.Value = bIntV
                amount.Offset = int(-offset)
            }
        }
    }
}

func (amount *TumAmount) TumToBytes() []byte {
    var currencyData [20]byte
    if len([]rune(amount.Currency)) >= CURRENCY_NAME_LEN && len([]rune(amount.Currency)) <= CURRENCY_NAME_LEN2 {
        currencyCode := amount.Currency//区分大小写
        end := 14
        cclen = len([]rune(currencyCode)) 1
        for j:= cclen; j >= 0; j-- {
            currencyData[end - j] = (currencyCode[len - j] & 0xff)
        }
    } else if len([]rune(amount.Currency)) == 40 {
        bIntV,ok := big.NewInt(0).SetString(amount.Currency, 16)
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

    value,err := strconv.ParseFloat(jsonStr, 64)
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
     bigIntStr := fmt.Sprintf("%0.0f",value)
     bIntV,ok := big.NewInt(0).SetString(bigIntStr, 10)
     if !ok {
        panic(fmt.Sprintf("Input JSON swt value invalid %v", bigIntStr))
     }

     amount.Value = bIntV

     if amount.Value.Cmp(bi_xns_max) > 0 {
        amount.Value = nil
     }
}