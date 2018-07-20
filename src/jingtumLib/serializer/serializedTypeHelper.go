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
    "fmt"
)

type KeyValuePair struct {
    Key     uint8
    Value   uint8
}

var (
    STInt8 SerializedInt8 = new(SerializedInt8)
    STInt16 SerializedInt16 = new(SerializedInt16)
    STInt32 SerializedInt32 = new(SerializedInt32)
    STInt64 SerializedInt64 = new(SerializedInt64)
    STMemo SerializedMemo = new(SerializedMemo)
    STArg SerializedArg = new(SerializedArg)
    STHash256 SerializedHash256 = new(SerializedHash256)
    STObject SerializedObject = new(SerializedObject)

    TYPES_MAP = map[uint8]ISerializedType{1:new(SerializedInt16),2:STInt32,3:new(SerializedInt64),4:new(SerializedHash128),5:STHash256,6:new(SerializedAmount),7:new(SerializedVariableLength),8:new(SerializedAccount),14:STObject,15:new(SerializedArray),16:STInt8,17:new(SerializedHash160),18:new(SerializedPathSet),19:new(SerializedVector256)}

	INVERSE_FIELDS_MAP = map[string]*KeyValuePair{
	    "LedgerEntryType":&KeyValuePair{1,1},
		"TransactionType":&KeyValuePair{1,2},
		"Flags":&KeyValuePair{2, 2},
		"SourceTag":&KeyValuePair{2, 3},
		"Sequence":&KeyValuePair{2, 4},
		"PreviousTxnLgrSeq":&KeyValuePair{2, 5},
		"LedgerSequence":&KeyValuePair{2, 6},
		"CloseTime":&KeyValuePair{2, 7},
		"ParentCloseTime":&KeyValuePair{2, 8},
		"SigningTime":&KeyValuePair{2, 9},
		"Expiration":&KeyValuePair{2, 10},
		"TransferRate":&KeyValuePair{2, 11},
		"WalletSize":&KeyValuePair{2, 12},
		"OwnerCount":&KeyValuePair{2, 13},
		"DestinationTag":&KeyValuePair{2, 14},
		"Timestamp":&KeyValuePair{2, 15},
		"HighQualityIn":&KeyValuePair{2, 16},
		"HighQualityOut":&KeyValuePair{2, 17},
		"LowQualityIn":&KeyValuePair{2, 18},
		"LowQualityOut":&KeyValuePair{2, 19},
		"QualityIn":&KeyValuePair{2, 20},
		"QualityOut":&KeyValuePair{2, 21},
		"StampEscrow":&KeyValuePair{2, 22},
		"BondAmount":&KeyValuePair{2, 23},
		"LoadFee":&KeyValuePair{2, 24},
		"OfferSequence":&KeyValuePair{2, 25},
		"FirstLedgerSequence":&KeyValuePair{2, 26},
		"LastLedgerSequence":&KeyValuePair{2, 27},
		"TransactionIndex":&KeyValuePair{2, 28},
		"OperationLimit":&KeyValuePair{2, 29},
		"ReferenceFeeUnits":&KeyValuePair{2, 30},
		"ReserveBase":&KeyValuePair{2, 31},
		"ReserveIncrement":&KeyValuePair{2, 32},
		"SetFlag":&KeyValuePair{2, 33},
		"ClearFlag":&KeyValuePair{2, 34},
		"RelationType":&KeyValuePair{2, 35},
		"Method":&KeyValuePair{2, 36},
		"Contracttype":&KeyValuePair{2, 39},
		"IndexNext":&KeyValuePair{3, 1},
		"IndexPrevious":&KeyValuePair{3, 2},
		"BookNode":&KeyValuePair{3, 3},
		"OwnerNode":&KeyValuePair{3, 4},
		"BaseFee":&KeyValuePair{3, 5},
		"ExchangeRate":&KeyValuePair{3, 6},
		"LowNode":&KeyValuePair{3, 7},
		"HighNode":&KeyValuePair{3, 8},
		"EmailHash":&KeyValuePair{4, 1},
		"LedgerHash":&KeyValuePair{5, 1},
		"ParentHash":&KeyValuePair{5, 2},
		"TransactionHash":&KeyValuePair{5, 3},
		"AccountHash":&KeyValuePair{5, 4},
		"PreviousTxnID":&KeyValuePair{5, 5},
		"LedgerIndex":&KeyValuePair{5, 6},
		"WalletLocator":&KeyValuePair{5, 7},
		"RootIndex":&KeyValuePair{5, 8},
		"AccountTxnID":&KeyValuePair{5, 9},
		"BookDirectory":&KeyValuePair{5, 16},
		"InvoiceID":&KeyValuePair{5, 17},
		"Nickname":&KeyValuePair{5, 18},
		"Amendment":&KeyValuePair{5, 19},
		"TicketID":&KeyValuePair{5, 20},
		"Amount":&KeyValuePair{6, 1},
		"Balance":&KeyValuePair{6, 2},
		"LimitAmount":&KeyValuePair{6, 3},
		"TakerPays":&KeyValuePair{6, 4},
		"TakerGets":&KeyValuePair{6, 5},
		"LowLimit":&KeyValuePair{6, 6},
		"HighLimit":&KeyValuePair{6, 7},
		"Fee":&KeyValuePair{6, 8},
		"SendMax":&KeyValuePair{6, 9},
		"MinimumOffer":&KeyValuePair{6, 16},
		"JingtumEscrow":&KeyValuePair{6, 17},
		"DeliveredAmount":&KeyValuePair{6, 18},
		"PublicKey":&KeyValuePair{7, 1},
		"MessageKey":&KeyValuePair{7, 2},
		"SigningPubKey":&KeyValuePair{7, 3},
		"TxnSignature":&KeyValuePair{7, 4},
		"Generator":&KeyValuePair{7, 5},
		"Signature":&KeyValuePair{7, 6},
		"Domain":&KeyValuePair{7, 7},
		"FundCode":&KeyValuePair{7, 8},
		"RemoveCode":&KeyValuePair{7, 9},
		"ExpireCode":&KeyValuePair{7, 10},
		"CreateCode":&KeyValuePair{7, 11},
		"MemoType":&KeyValuePair{7, 12},
		"MemoData":&KeyValuePair{7, 13},
		"MemoFormat":&KeyValuePair{7, 14},
		"Payload":&KeyValuePair{7, 15},
		"ContractMethod":&KeyValuePair{7, 17},
		"Parameter":&KeyValuePair{7, 18},
		"Account":&KeyValuePair{8, 1},
		"Owner":&KeyValuePair{8, 2},
		"Destination":&KeyValuePair{8, 3},
		"Issuer":&KeyValuePair{8, 4},
		"Target":&KeyValuePair{8, 7},
		"RegularKey":&KeyValuePair{8, 8},
		"undefined":&KeyValuePair{15, 1},
		"TransactionMetaData":&KeyValuePair{14, 2},
		"CreatedNode":&KeyValuePair{14, 3},
		"DeletedNode":&KeyValuePair{14, 4},
		"ModifiedNode":&KeyValuePair{14, 5},
		"PreviousFields":&KeyValuePair{14, 6},
		"FinalFields":&KeyValuePair{14, 7},
		"NewFields":&KeyValuePair{14, 8},
		"TemplateEntry":&KeyValuePair{14, 9},
		"Memo":&KeyValuePair{14, 10},
		"Arg":&KeyValuePair{14, 11},
		"SigningAccounts":&KeyValuePair{15, 2},
		"TxnSignatures":&KeyValuePair{15, 3},
		"Signatures":&KeyValuePair{15, 4},
		"Template":&KeyValuePair{15, 5},
		"Necessary":&KeyValuePair{15, 6},
		"Sufficient":&KeyValuePair{15, 7},
		"AffectedNodes":&KeyValuePair{15, 8},
		"Memos":&KeyValuePair{15, 9},
		"Args":&KeyValuePair{15, 10},
		"CloseResolution":&KeyValuePair{16, 1},
		"TemplateEntryType":&KeyValuePair{16, 2},
		"TransactionResult":&KeyValuePair{16, 3},
		"TakerPaysCurrency":&KeyValuePair{17, 1},
		"TakerPaysIssuer":&KeyValuePair{17, 2},
		"TakerGetsCurrency":&KeyValuePair{17, 3},
		"TakerGetsIssuer":&KeyValuePair{17, 4},
		"Paths":&KeyValuePair{18, 1},
		"Indexes":&KeyValuePair{19, 1},
		"Hashes":&KeyValuePair{19, 2},
		"Amendments":&KeyValuePair{19, 3}}
)

func SerializeHex(so *Serializer, string val, noLength bool) {
    bytes, err := HexToBytes(val)

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
        panic(fmt.Sprintf("Variable integers are unsigned %d", val))
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
        panic(fmt.Sprintf("Variable integer overflow %d", val))
    }
}

func getLedgerEntryType(structure interface{}) interface{} {
    var output interface{}
    switch v := structure.(type) {
        case uint8 :
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
            break;
        default:
            panic(fmt.Sprintf("Invalid input type for transaction type %v", v))
    }
    return output
}

func Serialize(so *Serializer, fieldName string, value interface{}) {
    fieldCoordinates := INVERSE_FIELDS_MAP[fieldName]
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

    STInt8.Serialize(so,byte(tagByte), false)

    if typeBits >= 16 {
        STInt8.Serialize(so,byte(typeBits),false)
    }

    if fieldBits >= 16 {
        STInt8.Serialize(so,byte(fieldBits),false)
    }

    var serializedType ISerializedType

    if _,ok := value.(MemoDataInfo); ok && fieldName == "Memo" {
        serializedType = STMemo
    } else {
        serializedType = TYPES_Map[typeBits]
    }
    
    serializedType.Serialize(so, value, false)
}

func main() {
    fmt.Println("Hello World",INVERSE_FIELDS_MAP["INVERSE_FIELDS_MAP"])
}
