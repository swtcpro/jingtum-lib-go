/**
 *
 * 文件功能介绍
 *
 * @FileName: bits.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-11 10:44:32
 * @UpdateTime: 2018-07-11 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"regexp"
	"sort"
	"strconv"
)

type FieldWrapper struct {
	fields []string
	by     func(p, q string) bool
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

func SortByFieldName(fields []string) {
	sortField(fields, func(p, q string) bool {
		xMap := INVERSE_FIELDS_MAP[p]
		xTypeBits := xMap.Key
		xFieldBits := xMap.Value
		yMap := INVERSE_FIELDS_MAP[q]
		yTypeBits := yMap.Key
		yFieldBits := yMap.Value

		if xTypeBits != yTypeBits {
			ret := xTypeBits - yTypeBits
			if ret > 0 || ret == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := xFieldBits - yFieldBits
			if ret > 0 || ret == 0 {
				return true
			} else {
				return false
			}
		}
	})
}

//func getLedgerEntryType(structure interface{}) interface{} {
//	var output interface{}
//	switch v := structure.(type) {
//	case uint8:
//		switch v {
//		case 97:
//			output = "AccountRoot"
//		case 99:
//			output = "Contract"
//		case 100:
//			output = "DirectoryNode"
//		case 102:
//			output = "EnabledFeatures"
//		case 115:
//			output = "FeeSettings"
//		case 103:
//			output = "GeneratorMap"
//		case 104:
//			output = "LedgerHashes"
//		case 110:
//			output = "Nickname"
//		case 111:
//			output = "Offer"
//		case 114:
//			output = "SkywellState"
//		default:
//			panic(fmt.Sprintf("Invalid input type for ransaction result %v", v))
//		}
//	case string:
//		switch v {
//		case "AccountRoot":
//			output = 97
//		case "Contract":
//			output = 99
//		case "DirectoryNode":
//			output = 100
//		case "EnabledFeatures":
//			output = 102
//		case "FeeSettings":
//			output = 115
//		case "GeneratorMap":
//			output = 103
//		case "LedgerHashes":
//			output = 104
//		case "Nickname":
//			output = 110
//		case "Offer":
//			output = 111
//		case "SkywellState":
//			output = 114
//		default:
//			output = 0
//		}
//	default:
//		output = "UndefinedLedgerEntry"
//	}
//	return output
//}
//
//func getTransactionType(structure interface{}) interface{} {
//	var output interface{}
//	switch v := structure.(type) {
//	case uint8:
//		switch v {
//		case 0:
//			output = "Payment"
//		case 3:
//			output = "AccountSet"
//		case 5:
//			output = "SetRegularKey"
//		case 7:
//			output = "OfferCreate"
//		case 8:
//			output = "OfferCancel"
//		case 9:
//			output = "Contract"
//		case 10:
//			output = "RemoveContract"
//		case 20:
//			output = "TrustSet"
//		case 100:
//			output = "EnableFeature"
//		case 101:
//			output = "SetFee"
//		default:
//			panic(fmt.Sprintf("Invalid transaction type %v", v))
//		}
//	case string:
//		switch v {
//		case "Payment":
//			output = 0
//		case "AccountSet":
//			output = 3
//		case "SetRegularKey":
//			output = 5
//		case "OfferCreate":
//			output = 7
//		case "OfferCancel":
//			output = 8
//		case "Contract":
//			output = 9
//		case "RemoveContract":
//			output = 10
//		case "TrustSet":
//			output = 20
//		case "EnableFeature":
//			output = 100
//		case "SetFee":
//			output = 101
//		default:
//			panic(fmt.Sprintf("Invalid transaction type %v", v))
//		}
//		break
//	default:
//		panic(fmt.Sprintf("Invalid input type for transaction type %v", v))
//	}
//	return output
//}
