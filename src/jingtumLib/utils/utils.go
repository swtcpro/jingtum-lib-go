/**
 * 通用工具类.
 *
 * @FileName: utils.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package utils

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	jtConst "jingtumLib/constant"
	jtEncode "jingtumLib/encoding"

	"github.com/yangxuebo-138/decimal"
)

/**
 * concat an item and a buffer
 * @param {integer} item1, should be an integer
 * @param {buffer} buf2, a buffer
 * @returns {buffer} new Buffer
 */
func bufCat0(item1 uint8, buf2 []byte) []byte {
	var buf []byte
	buf = append(buf, item1)
	buf = append(buf, buf2...)
	return buf
}

/**
 * concat one buffer and another
 * @param {buffer} item1, should be an integer
 * @param {buffer} buf2, a buffer
 * @returns {buffer} new Buffer
 */
func bufCat1(buf1 []byte, buf2 []byte) []byte {
	var buf []byte
	buf = append(buf, buf1...)
	buf = append(buf, buf2...)
	return buf
}

/**
 * encode use jingtum base58 encoding
 * including version + data + checksum
 * @param {integer} version
 * @param {buffer} bytes
 * @returns {string}
 * @private
 */
func EncodeB58(version uint8, bytes []byte) string {
	buffer := bufCat0(version, bytes)
	checksum := Sha256Util(Sha256Util(buffer))[0:4]
	ret := bufCat1(buffer, checksum)
	encodedString := jtEncode.Base58Encode(ret, jtEncode.JingTumAlphabet)
	return encodedString
}

func DecodeB58(version uint8, input string) (decodedBytes []byte, err error) {
	decodedBytes, err = jtEncode.Base58Decode(input, jtEncode.JingTumAlphabet)
	if err != nil || decodedBytes[0] != version || len(decodedBytes) < 5 {
		err = errors.New("invalid input size")
		return
	}

	computed := Sha256Util(Sha256Util(decodedBytes[0 : len(decodedBytes)-4]))[0:4]
	checksum := decodedBytes[len(decodedBytes)-4:]

	for i := 0; i != 4; i++ {
		if computed[i] != checksum[i] {
			err = errors.New("invalid checksum")
			return
		}
	}

	decodedBytes = decodedBytes[1 : len(decodedBytes)-4]

	return
}

func BytesToBigInt(b []byte) *big.Int {
	b_buf := bytes.NewBuffer(b)
	var x big.Int
	binary.Read(b_buf, binary.BigEndian, &x)
	return &x
}

func IsValidCurrency(currency string) bool {
	ok, _ := regexp.MatchString(jtConst.REGEX_CURRENCY, currency)
	return ok
}

//DecodeAddress 地址解码。
func DecodeAddress(address string) ([]byte, error) {
	decodedBytes, err := DecodeB58(jtConst.ACCOUNT_PREFIX, address)
	if err != nil {
		return nil, fmt.Errorf("Issuer invalid issuer info %s", address)
	}

	return decodedBytes, nil
}

func IsValidAddress(address string) bool {
	if address == "" {
		return false
	}

	_, err := DecodeB58(jtConst.ACCOUNT_PREFIX, address)
	if err != nil {
		return false
	}
	return true
}

type FieldWrapper struct {
	fields []string
	by     func(p, q string) bool
}

type TestWrapper struct {
	Ids []string
	By  func(i, j string) bool
}

func (tw TestWrapper) Len() int {
	return len(tw.Ids)
}

func (tw TestWrapper) Swap(i, j int) {
	tw.Ids[i], tw.Ids[j] = tw.Ids[j], tw.Ids[i]
}

func (tw TestWrapper) Less(i, j int) bool {
	return tw.By(tw.Ids[i], tw.Ids[j])
}

type SortBy func(p, q string) bool

func sortField(fields []string, by SortBy) {
	sort.Sort(FieldWrapper{fields, by})
}

func (fw FieldWrapper) Len() int {
	return len(fw.fields)
}
func (fw FieldWrapper) Swap(i, j int) {
	fw.fields[i], fw.fields[j] = fw.fields[j], fw.fields[i]
}
func (fw FieldWrapper) Less(i, j int) bool {
	return fw.by(fw.fields[i], fw.fields[j])
}

func GetBytes(value interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, value)
	return bytesBuffer.Bytes()
}

func MatchString(patter string, str string) bool {
	match, _ := regexp.MatchString(patter, str)

	return match
}

func IsNumberType(obj interface{}) bool {
	switch obj.(type) {
	case float64, float32, int, int8, int32, int64, byte, uint32, uint64:
		return true
	default:
		return false
	}
}

func IsNumberString(s string) bool {

	return MatchString("^(-?)(\\d*)(\\.\\d{0,6})?$", s)
}

func NumberToString(obj interface{}) string {
	switch v := obj.(type) {
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
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
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(uint64(v), 10)
	default:
		return ""
	}
}

func IsHexString(str string) bool {
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", str)
	return match
}

/**
 * 根据货币类型转换成相应的金额对象。如果是SWT则返回基本数据类型
 */
func ToAmount(amount jtConst.Amount) (interface{}, error) {
	if amount.Value == "" {
		return nil, jtConst.ERR_EMPTY_PARAM
	}

	value, err := strconv.ParseFloat(amount.Value, 64)

	// value, err := decimal.NewFromString(amount.Value)

	if err != nil {
		return nil, err
	}

	// vf64, ok := value.Float64()

	// if !ok {
	// 	return nil, errors.New(fmt.Sprintf("Parse float %s error.", amount.Value))
	// }

	if value > 100000000000 {
		return nil, jtConst.ERR_PAYMENT_OUT_OF_AMOUNT
	}

	if amount.Currency == jtConst.CFGCurrency {
		// mul, err := decimal.NewFromString("1000000")
		// if err != nil {
		// 	return nil, errors.New(fmt.Sprintf("Parse float %s error.", "1000000"))
		// }

		retf64, ok := decimal.NewFromFloat(value).Mul(decimal.NewFromFloat(1000000)).Float64()

		if !ok {
			return nil, fmt.Errorf("Parse float %s error", decimal.NewFromFloat(value).Mul(decimal.NewFromFloat(1000000)).String())
		}
		return retf64, nil
	}

	return amount, nil
}

/**
 *  获取对象字段存储的值
 */
func GetFieldValue(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)

	if v.Kind() == reflect.Struct {
		return v.FieldByName(fieldName).Interface()
	} else if v.Kind() == reflect.Ptr {
		return v.Elem().FieldByName(fieldName).Interface()
	} else {
		return nil
	}
}

/**
 *  获取对象的字段名
 */
func GetFieldNames(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	var fields []string
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, t.Field(i).Name)
		}
	} else if t.Kind() == reflect.Ptr {
		t = t.Elem()
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, t.Field(i).Name)
		}
	}

	return fields
}

func HexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

func StringToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

func ByteToHexString(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func HexToString(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	return string(bytes), err
}

func IsValidAmount(amount *jtConst.Amount) bool {
	if nil == amount {
		return false
	}

	if amount.Value == "" {
		return false
	}

	// check amount currency
	if (amount.Currency == "") || !IsValidCurrency(amount.Currency) {
		return false
	}

	// native currency issuer is empty
	if amount.Currency == jtConst.CFGCurrency && amount.Issuer != "" {
		return false
	}

	// non native currency issuer is not allowed to be empty
	if amount.Currency != jtConst.CFGCurrency && !IsValidAddress(amount.Issuer) {
		return false
	}

	return true
}

//SortByFieldName SortByFieldName
func SortByFieldName(fields []string) {
	sortField(fields, func(p, q string) bool {
		xMap := jtConst.InverseFieldsMap[p]
		xTypeBits := xMap.Key
		xFieldBits := xMap.Value
		yMap := jtConst.InverseFieldsMap[q]
		yTypeBits := yMap.Key
		yFieldBits := yMap.Value
		if xTypeBits != yTypeBits {
			return (xTypeBits - yTypeBits) < 0
		} else {
			return (xFieldBits - yFieldBits) < 0
		}
	})
}

//DeepCopy 对象深度copy
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	} else if valueList, ok := value.(*list.List); ok {
		vl := list.New()
		for e := valueList.Front(); e != nil; e = e.Next() {
			vl.PushBack(DeepCopy(e.Value))
		}

		return vl
	}

	return value
}
