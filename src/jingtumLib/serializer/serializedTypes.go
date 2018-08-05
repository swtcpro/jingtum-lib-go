/**
 *
 * 各数据类型的序列化实现
 *
 * @FileName: serializedTypes.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-11 10:44:32
 * @UpdateTime: 2018-07-11 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"bytes"
	"container/list"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"

	"jingtumLib/constant"
	jtUtils "jingtumLib/utils"
)

var (
	currencyNameLen  = 3
	currencyNameLen2 = 6
	typeBoundary     = 0xff
	typeEnd          = 0x00
	typeAccount      = 0x01
	typeCurrency     = 0x10
	typeIssuer       = 0x20
)

// ISerializedType 是一个序列化接口。
type ISerializedType interface {
	Serialize(so *Serializer, val interface{}, noMarker bool)
	Parse(so *Serializer) interface{}
}

// PathComputed 结构体。
type PathComputed struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Value    string `json:"value"`
	Account  string `json:"account"`
	Type     int    `json:"type"`
	TypeHex  string `json:"type_hex"`
}

// PathData 结构体。
type PathData struct {
	PathsComputed [][]PathComputed
	Choice        interface{}
}

//SerializedInt16 int16
type SerializedInt16 struct {
}

//SerializedInt32 int32
type SerializedInt32 struct {
}

//SerializedInt64 int64
type SerializedInt64 struct {
}

//SerializedMemo memo
type SerializedMemo struct {
}

//SerializedArg args
type SerializedArg struct {
}

//SerializedHash128 hash 128
type SerializedHash128 struct {
}

//SerializedHash256 hash 256
type SerializedHash256 struct {
}

//SerializedAmount amount
type SerializedAmount struct {
}

//SerializedCurrency currency
type SerializedCurrency struct {
}

//SerializedObject object
type SerializedObject struct {
}

//SerializedArray array
type SerializedArray struct {
}

//SerializedHash160 hash 160
type SerializedHash160 struct {
}

//SerializedPathSet pathset
type SerializedPathSet struct {
}

//SerializedVector256 vector 256
type SerializedVector256 struct {
}

//SerializedVariableLength variable length
type SerializedVariableLength struct {
}

//SerializedAccount account
type SerializedAccount struct {
}

//Serialize int8
func (serInt8 SerializedInt8) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuit8, ok := val.(uint8); ok {
		so.Append(jtUtils.GetBytes(vuit8))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint8 || vint < 0 {
			so.err = fmt.Errorf("Value out of bounds %d", vint)
			return
		}

		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = fmt.Errorf("Serialize int8 type error %T, %v", val, val)
		return
	}
}

//Parse int8
func (serInt8 SerializedInt8) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize int16
func (serInt16 SerializedInt16) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuint16, ok := val.(uint16); ok {
		so.Append(jtUtils.GetBytes(vuint16))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint16 || vint < 0 {
			so.err = fmt.Errorf("Value out of bounds %d", vint)
			return
		}

		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = fmt.Errorf("Serialize int16 type error %T, %v", val, val)
		return
	}
}

//Parse int16
func (serInt16 SerializedInt16) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize int32
func (serInt32 SerializedInt32) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuint32, ok := val.(uint32); ok {
		so.Append(jtUtils.GetBytes(vuint32))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint32 || vint < 0 {
			so.err = fmt.Errorf("Value out of bounds %d", vint)
			return
		}
		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = fmt.Errorf("Serialize int16 type error %T, %v", val, val)
		return
	}
}

//Parse int32
func (serInt32 SerializedInt32) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize int64
func (serInt64 SerializedInt64) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if number, ok := val.(uint64); ok {
		so.Append(jtUtils.GetBytes(number))
		return
	}

	if str, ok := val.(string); ok {
		if !jtUtils.IsHexString(str) {
			so.err = fmt.Errorf("Invalid hex string %s", str)
			return
		}

		if len(str) > 16 {
			so.err = fmt.Errorf("Int64 is too large %s", str)
			return
		}

		b := bytes.NewBufferString("")

		for b.Len() < 16-len(str) {
			b.WriteString("0")
		}

		b.WriteString(str)

		SerializeHex(so, b.String(), true)
		return
	}

	so.err = fmt.Errorf("Invalid type for Int64 %T, %v", val, val)
}

//Parse int64
func (serInt64 SerializedInt64) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Parse memo
func (serMemo SerializedMemo) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize memo
func (serMemo SerializedMemo) Serialize(so *Serializer, val interface{}, noMarker bool) {
	memo, ok := val.(*MemoDataInfo)
	if !ok {
		so.err = fmt.Errorf("Serialize Memo type error %T", val)
		return
	}

	fileds := jtUtils.GetFieldNames(val)
	for i := 0; i < len(fileds); i++ {
		_, ok := constant.INVERSE_FIELDS_MAP[fileds[i]]
		if !ok {
			so.err = fmt.Errorf("JSON contains unknown field : %s", fileds[i])
			return
		}
	}

	jtUtils.SortByFieldName(fileds)

	isJSON := memo.MemoFormat == "json"

	for _, fn := range fileds {
		value := jtUtils.GetFieldValue(val, fn)
		switch fn {
		case "MemoType", "MemoFormat":
			value = jtUtils.StringToHex(value.(string))
		case "MemoData":
			if _, ok := value.(string); ok {
				value = jtUtils.StringToHex(value.(string))
				break
			}
			if isJSON {
				mjson, _ := json.Marshal(value)
				value = jtUtils.StringToHex(string(mjson))
				break
			}
			so.err = fmt.Errorf("MemoData can only be a JSON object with a valid json MemoFormat. %v. %T", value, value)
			return
		}

		Serialize(so, fn, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

//Parse arg
func (serArg SerializedArg) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize arg
func (serArg SerializedArg) Serialize(so *Serializer, val interface{}, noMarker bool) {
	fileds := jtUtils.GetFieldNames(val)
	for i := 0; i < len(fileds); i++ {
		kvp := constant.INVERSE_FIELDS_MAP[fileds[i]]
		if kvp == nil {
			so.err = fmt.Errorf("JSON contains unknown field %s", fileds[i])
			return
		}
	}
	jtUtils.SortByFieldName(fileds)

	for _, fn := range fileds {
		if fn == "Parameter" {
			break
		}
		value := jtUtils.GetFieldValue(val, fn)
		Serialize(so, fn, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

//Parse Hash128反序列化。
func (serHash128 SerializedHash128) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize Hash128序列化。
func (serHash128 SerializedHash128) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if v, ok := val.(string); ok && jtUtils.MatchString("^[0-9A-F]{0,32}$", v) && len(v) <= 32 {
		SerializeHex(so, v, true)
		return
	}
}

//Parse Hash256反序列化。
func (serHash256 SerializedHash256) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize Hash256序列化。
func (serHash256 SerializedHash256) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if v, ok := val.(string); ok && jtUtils.MatchString("^[0-9A-F]{0,32}$", v) && len(v) <= 64 {
		SerializeHex(so, v, true)
		return
	}
}

//Parse 金额反序列化。
func (serAmount SerializedAmount) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize 金额序列化。
func (serAmount SerializedAmount) Serialize(so *Serializer, val interface{}, noMarker bool) {
	tumAmount, err := fromJSON(val)
	if err != nil {
		so.err = err
		return
	}
	if !tumAmount.IsValid() {
		so.err = fmt.Errorf("Not a valid Amount object")
		return
	}

	if tumAmount.IsNative {
		valueHex := hex.EncodeToString(tumAmount.Value.Bytes())

		if len(valueHex) > 16 {
			so.err = fmt.Errorf("Amount value out of bounds")
		}
		b := bytes.NewBufferString("")
		for b.Len() < 16 {
			b.WriteString("0")
		}
		b.WriteString(valueHex)

		valueBytes, err := jtUtils.HexToBytes(b.String())
		if err != nil {
			so.err = fmt.Errorf("Hex to bytes error %s", b.String())
			return
		}

		valueBytes[0] &= 0x3f

		if tumAmount.IsNegative {
			valueBytes[0] |= 0x40
		}
		so.Append(valueBytes)
	} else {
		//For other non-native currency
		//1. Serialize the currency value with offset
		//Put offset
		var hi, lo int64 = 0, 0
		hi |= 1 << 31
		if !tumAmount.IsZeroM() {
			// Second bit: non-negative?
			if tumAmount.IsNegative {
				hi |= 1 << 30
			}
			// Next eight bits: offset/exponent
			hi |= ((int64(97) + int64(tumAmount.Offset)) & 0xff) << 22
			// Remaining 54 bits: mantissa
			hi |= (tumAmount.Value.Int64() >> 32) & 0x3fffff
			lo = tumAmount.Value.Int64() & 0xffffffff
		}

		// Convert from a bitArray to an array of bytes.
		arr := []int64{hi, lo}
		l := len(arr)
		var bl int64
		if l == 0 {
			bl = 0
		} else {
			x := arr[l-1]
			roundX := x / 0x10000000000
			if roundX == 0 {
				roundX = 32
			}
			bl = int64((l-1)*32) + int64(roundX)
		}

		var tmparray []byte
		var tmp int64
		for i := 0; int64(i) < int64(bl/8); i++ {
			if (i & 3) == 0 {
				tmp = arr[i/4]
			}

			tmparray = append(tmparray, byte(tmp>>24))
			tmp <<= 8
		}

		if len(tmparray) > 8 {
			so.err = fmt.Errorf("Invalid byte array length in AMOUNT value representation")
			return
		}

		so.Append(tmparray)
		tumBytes, err := tumAmount.TumToBytes()
		if err != nil {
			so.err = err
			return
		}
		so.Append(tumBytes)
		addrByte, err := jtUtils.DecodeAddress(tumAmount.Issuer)
		if err != nil {
			so.err = err
			return
		}
		so.Append(addrByte)
	}
}

//Parse currency
func (serCurrency SerializedCurrency) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize currency
func (serCurrency SerializedCurrency) Serialize(so *Serializer, val interface{}, noMarker bool) {
	currencty := val.(string)
	so.Append(serCurrency.fromJSONToBytes(currencty))
}

func (serCurrency SerializedCurrency) fromJSONToBytes(currencty string) []byte {
	var result []byte
	if currencty != "" {
		if jtUtils.IsHexString(currencty) && len(currencty) == 40 {
			var err error
			result, err = jtUtils.HexToBytes(currencty)

			if err != nil {
				panic("Invalid currencty.")
			}

		} else if jtUtils.IsValidCurrency(currencty) {
			if len(currencty) >= currencyNameLen && len(currencty) <= currencyNameLen2 {
				var end = 14
				var clen = len(currencty) - 1
				for x := clen; x >= 0; x-- {
					result[end-x] = byte(currencty[clen-x] & 0xff)
				}
			}
		} else {
			panic(fmt.Sprintf("Input tum code invalid %v", currencty))
		}
	} else {
		panic(fmt.Sprintf("Input tum code invalid %v", currencty))
	}

	return result
}

//Parse object
func (serObject SerializedObject) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize object
func (serObject SerializedObject) Serialize(so *Serializer, val interface{}, noMarker bool) {
	txData, ok := val.(map[string]interface{})
	if !ok {
		so.err = fmt.Errorf("Serialive object type must be map[string]interface{}. Actual type : %T. Value : %v", val, val)
		return
	}

	var fieldNames []string

	for k := range txData {
		_, ok := constant.INVERSE_FIELDS_MAP[k]
		if !ok {
			so.err = fmt.Errorf("Not fund field name %s", k)
			return
		}

		fieldNames = append(fieldNames, k)
	}

	jtUtils.SortByFieldName(fieldNames)

	for _, field := range fieldNames {
		value := txData[field]
		if value == nil {
			continue
		}

		Serialize(so, field, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

//Parse array
func (serArray SerializedArray) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize array
func (serArray SerializedArray) Serialize(so *Serializer, val interface{}, noMarker bool) {
	array, ok := val.(*list.List)
	if !ok {
		so.err = fmt.Errorf("Serialize array type error %T. Value : %v", val, val)
		return
	}

	for e := array.Front(); e != nil; e = e.Next() {
		fields := jtUtils.GetFieldNames(e.Value)

		if len(fields) != 1 {
			so.err = fmt.Errorf("Cannot serialize an array containing non-single-key objects")
			return
		}

		value := jtUtils.GetFieldValue(e.Value, fields[0])

		Serialize(so, fields[0], value)
	}

	STInt8.Serialize(so, 0xf1, false)
}

//Parse hash 160
func (serHash160 SerializedHash160) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize hash 160
func (serHash160 SerializedHash160) Serialize(so *Serializer, val interface{}, noMarker bool) {
	valStr := val.(string)
	SerializeHex(so, valStr, true)
}

//Parse path set
func (serPathSet SerializedPathSet) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize path set
func (serPathSet SerializedPathSet) Serialize(so *Serializer, val interface{}, noMarker bool) {
	path := val.([][]PathComputed)
	for i := 0; i < len(path); i++ {
		if i > 0 {
			STInt8.Serialize(so, typeBoundary, false)
		}

		pathes := path[i]
		for j, l2 := 0, len(pathes); j < l2; j++ {
			entry := pathes[j]
			typev := 0
			if entry.Account != "" {
				typev |= typeAccount
			}

			if entry.Currency != "" {
				typev |= typeCurrency
			}

			if entry.Issuer != "" {
				typev |= typeIssuer
			}

			STInt8.Serialize(so, typev, false)

			if entry.Account != "" {
				addrByte, err := jtUtils.DecodeAddress(entry.Account)
				if err != nil {
					so.err = err
					return
				}
				so.Append(addrByte)
			}

			if entry.Currency != "" {
				sc := new(SerializedCurrency)
				so.Append(sc.fromJSONToBytes(entry.Currency))
			}

			if entry.Issuer != "" {
				addrByte, err := jtUtils.DecodeAddress(entry.Issuer)
				if err != nil {
					so.err = err
					return
				}
				so.Append(addrByte)
			}
		}
	}

	STInt8.Serialize(so, typeEnd, false)
}

//Parse Vector 256 反序列化。
func (serVector256 SerializedVector256) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize Vector 256 序列化。
func (serVector256 SerializedVector256) Serialize(so *Serializer, val interface{}, noMarker bool) {
	array := val.([]string)
	SerializeVarint(so, uint(len(array)*32))

	for _, v := range array {
		STHash256.Serialize(so, v, false)
	}
}

//Parse variable length 反序列化。
func (serVL SerializedVariableLength) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize variable length 序列化。
func (serVL SerializedVariableLength) Serialize(so *Serializer, val interface{}, noMarker bool) {
	vl, ok := val.(string)
	if !ok {
		so.err = fmt.Errorf("Serialized Variable type error %T", val)
		return
	}
	SerializeHex(so, vl, false)
}

//Parse 账号反序列化。
func (serAccount SerializedAccount) Parse(so *Serializer) interface{} {
	return fmt.Errorf("Not implemented error")
}

//Serialize 账号序列化。
func (serAccount SerializedAccount) Serialize(so *Serializer, val interface{}, noMarker bool) {
	account, ok := val.(string)
	if !ok {
		so.err = fmt.Errorf("Account type error %T", val)
		return
	}
	addrByte, err := jtUtils.DecodeAddress(account)
	if err != nil {
		so.err = err
		return
	}
	SerializeVarint(so, uint(len(addrByte)))
	so.Append(addrByte)
}
