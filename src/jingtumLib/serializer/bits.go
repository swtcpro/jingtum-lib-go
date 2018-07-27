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

/*import (
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
}*/

/**
 *  获取对象字段存储的值
 */
//func GetFieldValue(obj interface{}, fieldName string) interface{} {
//	v := reflect.ValueOf(obj)
//
//	if v.Kind() == reflect.Struct {
//		return v.FieldByName(fieldName).Interface()
//	} else if v.Kind() == reflect.Ptr {
//		return v.Elem().FieldByName(fieldName).Interface()
//	} else {
//		return nil
//	}
//}

/**
 *  获取对象的字段名
 */
//func GetFieldNames(obj interface{}) []string {
//	t := reflect.TypeOf(obj)
//	var fields []string
//	if t.Kind() == reflect.Struct {
//		for i := 0; i < t.NumField(); i++ {
//			fields = append(fields, t.Field(i).Name)
//		}
//	} else if t.Kind() == reflect.Ptr {
//		t = t.Elem()
//		for i := 0; i < t.NumField(); i++ {
//			fields = append(fields, t.Field(i).Name)
//		}
//	}
//
//	return fields
//}
//
//func HexToBytes(hexStr string) ([]byte, error) {
//	return hex.DecodeString(hexStr)
//}
//
//func StringToHex(str string) string {
//	return hex.EncodeToString([]byte(str))
//}
//
//func SortByFieldName(fields []string) {
//	sortField(fields, func(p, q string) bool {
//		xMap := constant.INVERSE_FIELDS_MAP[p]
//		xTypeBits := xMap.Key
//		xFieldBits := xMap.Value
//		yMap := constant.INVERSE_FIELDS_MAP[q]
//		yTypeBits := yMap.Key
//		yFieldBits := yMap.Value
//
//		if xTypeBits != yTypeBits {
//			ret := xTypeBits - yTypeBits
//			if ret > 0 || ret == 0 {
//				return true
//			} else {
//				return false
//			}
//		} else {
//			ret := xFieldBits - yFieldBits
//			if ret > 0 || ret == 0 {
//				return true
//			} else {
//				return false
//			}
//		}
//	})
//}
