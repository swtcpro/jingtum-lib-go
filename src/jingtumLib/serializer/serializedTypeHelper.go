/**
 *
 * 文件功能介绍
 *
 * @FileName: serializedTypeHelper.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-11 10:44:32
 * @UpdateTime: 2018-07-11 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"errors"
	"fmt"

	"jingtumLib/constant"
	"jingtumLib/utils"
)

var (
	STInt8    = new(SerializedInt8)
	STInt16   = new(SerializedInt16)
	STInt32   = new(SerializedInt32)
	STInt64   = new(SerializedInt64)
	STMemo    = new(SerializedMemo)
	STArg     = new(SerializedArg)
	STHash256 = new(SerializedHash256)
	STObject  = new(SerializedObject)

	TYPES_MAP = map[uint8]ISerializedType{1: new(SerializedInt16), 2: STInt32, 3: new(SerializedInt64), 4: new(SerializedHash128), 5: STHash256, 6: new(SerializedAmount), 7: new(SerializedVariableLength), 8: new(SerializedAccount), 14: STObject, 15: new(SerializedArray), 16: STInt8, 17: new(SerializedHash160), 18: new(SerializedPathSet), 19: new(SerializedVector256)}
)

func SerializeHex(so *Serializer, val string, noLength bool) {
	bytes, err := utils.HexToBytes(val)

	if err != nil {
		panic(fmt.Sprintf("Invalid hex string %v", val))
	}

	if len(bytes) == 0 {
		bytes = []byte{0}
	}

	if !noLength {
		SerializeVarint(so, uint(len(bytes)))
	}

	so.Append(bytes)
}

func SerializeVarint(so *Serializer, val uint) {
	if val < 0 {
		so.err = errors.New(fmt.Sprintf("Variable integers are unsigned %d", val))
		return
	}

	if val <= 192 {
		so.Append([]byte{byte(val)})
	} else if val <= 12480 {
		val -= 193
		so.Append([]byte{byte(193 + (val >> 8)), byte(val & 0xff)})
	} else if val <= 918744 {
		val -= 12481
		so.Append([]byte{byte(241 + (val >> 16)), byte(val >> 8 & 0xff), byte(val & 0xff)})
	} else {
		so.err = errors.New(fmt.Sprintf("Variable integer overflow %d", val))
	}
}

func getLedgerEntryType(structure interface{}) interface{} {
	var output interface{}
	switch v := structure.(type) {
	case uint8:
		switch v {
		case 97:
			output = "AccountRoot"
		case 99:
			output = "Contract"
		case 100:
			output = "DirectoryNode"
		case 102:
			output = "EnabledFeatures"
		case 115:
			output = "FeeSettings"
		case 103:
			output = "GeneratorMap"
		case 104:
			output = "LedgerHashes"
		case 110:
			output = "Nickname"
		case 111:
			output = "Offer"
		case 114:
			output = "SkywellState"
		default:
			panic(fmt.Sprintf("Invalid input type for ransaction result %v", v))
		}
	case string:
		switch v {
		case "AccountRoot":
			output = 97
		case "Contract":
			output = 99
		case "DirectoryNode":
			output = 100
		case "EnabledFeatures":
			output = 102
		case "FeeSettings":
			output = 115
		case "GeneratorMap":
			output = 103
		case "LedgerHashes":
			output = 104
		case "Nickname":
			output = 110
		case "Offer":
			output = 111
		case "SkywellState":
			output = 114
		default:
			output = 0
		}
	default:
		output = "UndefinedLedgerEntry"
	}
	return output
}

func getTransactionType(structure interface{}) interface{} {
	var output interface{}
	switch v := structure.(type) {
	case uint8:
		switch v {
		case 0:
			output = "Payment"
		case 3:
			output = "AccountSet"
		case 5:
			output = "SetRegularKey"
		case 7:
			output = "OfferCreate"
		case 8:
			output = "OfferCancel"
		case 9:
			output = "Contract"
		case 10:
			output = "RemoveContract"
		case 20:
			output = "TrustSet"
		case 100:
			output = "EnableFeature"
		case 101:
			output = "SetFee"
		default:
			panic(fmt.Sprintf("Invalid transaction type %v", v))
		}
	case string:
		switch v {
		case "Payment":
			output = 0
		case "AccountSet":
			output = 3
		case "SetRegularKey":
			output = 5
		case "OfferCreate":
			output = 7
		case "OfferCancel":
			output = 8
		case "Contract":
			output = 9
		case "RemoveContract":
			output = 10
		case "TrustSet":
			output = 20
		case "EnableFeature":
			output = 100
		case "SetFee":
			output = 101
		default:
			panic(fmt.Sprintf("Invalid transaction type %v", v))
		}
		break
	default:
		panic(fmt.Sprintf("Invalid input type for transaction type %v", v))
	}
	return output
}

func Serialize(so *Serializer, fieldName string, value interface{}) {
	fieldCoordinates, ok := constant.INVERSE_FIELDS_MAP[fieldName]
	if !ok {
		so.err = errors.New(fmt.Sprintf("Not fund field name %s.", fieldName))
		return
	}

	typeBits := fieldCoordinates.Key
	fieldBits := fieldCoordinates.Value
	var temp uint8
	var temp2 uint8
	if typeBits < 16 {
		temp = typeBits << 4
	} else {
		temp = 0
	}

	if fieldBits < 16 {
		temp2 = fieldBits
	} else {
		temp2 = 0
	}
	tagByte := temp | temp2

	if v, ok := value.(string); ok {
		if fieldName == "LedgerEntryType" {
			value = getLedgerEntryType(v)
		} else if fieldName == "TransactionResult" {
			value = getTransactionType(v)
		}
	}

	STInt8.Serialize(so, byte(tagByte), false)

	if typeBits >= 16 {
		STInt8.Serialize(so, byte(typeBits), false)
	}

	if fieldBits >= 16 {
		STInt8.Serialize(so, byte(fieldBits), false)
	}

	var serializedType ISerializedType

	if _, ok := value.(*MemoInfo); ok && fieldName == "Memo" {
		serializedType = STMemo
	} else {
		serializedType = TYPES_MAP[typeBits]
	}

	serializedType.Serialize(so, value, false)
}
