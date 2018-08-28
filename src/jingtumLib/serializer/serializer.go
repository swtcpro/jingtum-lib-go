/**
 *
 * 文件功能介绍
 *
 * @FileName: serializer.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-09 10:44:32
 * @UpdateTime: 2018-07-09 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"encoding/hex"
	"fmt"
	"strconv"

	jtUtils "jingtumlib/utils"
)

//Serializer struct
type Serializer struct {
	Buffer []byte
	err    error
}

//SerializedInt8 int8
type SerializedInt8 struct {
}

var (
	required = 0
	optional = 1
	defaultv = 2

	//交易类型
	transactionTypeAccountSet     = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"EmailHash", optional}, {"WalletLocator", optional}, {"WalletSize", optional}, {"MessageKey", optional}, {"Domain", optional}, {"TransferRate", optional}}
	transactionTypeTrustSet       = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"LimitAmount", optional}, {"QualityIn", optional}, {"QualityOut", optional}}
	transactionTypeOfferCreate    = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"TakerPays", required}, {"TakerGets", required}, {"Expiration", optional}}
	transactionTypeOfferCancel    = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"OfferSequence", required}}
	transactionTypeSetRegularKey  = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"RegularKey", required}}
	transactionTypePayment        = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Destination", required}, {"Amount", required}, {"SendMax", optional}, {"Paths", defaultv}, {"InvoiceID", optional}, {"DestinationTag", optional}}
	transactionTypeContract       = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Expiration", required}, {"BondAmount", required}, {"StampEscrow", required}, {"JingtumEscrow", required}, {"CreateCode", optional}, {"FundCode", optional}, {"RemoveCode", optional}, {"ExpireCode", optional}}
	transactionTypeRemoveContract = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Target", required}}
	transactionTypeEnableFeature  = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Feature", required}}
	transactionTypeSetFee         = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Features", required}, {"BaseFee", required}, {"ReferenceFeeUnits", required}, {"ReserveBase", required}, {"ReserveIncrement", required}}
	transactionTypeConfigContract = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Method", required}, {"Payload", optional}, {"Destination", optional}, {"Amount", optional}, {"Contracttype", optional}, {"ContractMethod", optional}, {"Args", optional}}
	transactionTypeRelationSet = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Method", required}, {"Payload", optional}, {"Destination", optional}, {"Amount", optional}, {"Contracttype", optional}, {"ContractMethod", optional}, {"Args", optional}}
	transactionTypeRelationDel = [][]interface{}{{"TransactionType", required}, {"Flags", optional}, {"SourceTag", optional}, {"LastLedgerSequence", optional}, {"Account", required}, {"Sequence", optional}, {"Fee", required}, {"OperationLimit", optional}, {"SigningPubKey", optional}, {"TxnSignature", optional}, {"Method", required}, {"Payload", optional}, {"Destination", optional}, {"Amount", optional}, {"Contracttype", optional}, {"ContractMethod", optional}, {"Args", optional}}
	transactionTypes              = map[uint8][][]interface{}{3: transactionTypeAccountSet, 20: transactionTypeTrustSet, 7: transactionTypeOfferCreate, 8: transactionTypeOfferCancel, 5: transactionTypeSetRegularKey, 0: transactionTypePayment, 9: transactionTypeContract, 10: transactionTypeRemoveContract, 100: transactionTypeEnableFeature, 101: transactionTypeSetFee, 30: transactionTypeConfigContract, 21: transactionTypeRelationSet, 22:transactionTypeRelationDel}
	txTypeStrMapNumber            = map[string]uint8{"AccountSet": 3, "TrustSet": 20, "OfferCreate": 7, "OfferCancel": 8, "SetRegularKey": 5, "Payment": 0, "Contract": 9, "RemoveContract": 10, "EnableFeature": 100, "SetFee": 101, "ConfigContract": 30, "RelationSet": 21, "RelationDel": 22}

	ledgerEntryTypeAccountRoot     = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"Sequence", required}, {"PreviousTxnLgrSeq", required}, {"TransferRate", optional}, {"WalletSize", optional}, {"OwnerCount", required}, {"EmailHash", optional}, {"PreviousTxnID", required}, {"AccountTxnID", optional}, {"WalletLocator", optional}, {"Balance", required}, {"MessageKey", optional}, {"Domain", optional}, {"Account", required}, {"RegularKey", optional}}
	ledgerEntryTypeContract        = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"PreviousTxnLgrSeq", required}, {"Expiration", required}, {"BondAmount", required}, {"PreviousTxnID", required}, {"Balance", required}, {"FundCode", optional}, {"RemoveCode", optional}, {"ExpireCode", optional}, {"CreateCode", optional}, {"Account", required}, {"Owner", required}, {"Issuer", required}}
	ledgerEntryTypeDirectoryNode   = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"IndexNext", optional}, {"IndexPrevious", optional}, {"ExchangeRate", optional}, {"RootIndex", required}, {"Owner", optional}, {"TakerPaysCurrency", optional}, {"TakerPaysIssuer", optional}, {"TakerGetsCurrency", optional}, {"TakerGetsIssuer", optional}, {"Indexes", required}}
	ledgerEntryTypeEnabledFeatures = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"Features", required}}
	ledgerEntryTypeFeeSettings     = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"ReferenceFeeUnits", required}, {"ReserveBase", required}, {"ReserveIncrement", required}, {"BaseFee", required}, {"LedgerIndex", optional}}
	ledgerEntryTypeGeneratorMap    = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"Generator", required}}
	ledgerEntryTypeLedgerHashes    = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"LedgerEntryType", required}, {"Flags", required}, {"FirstLedgerSequence", optional}, {"LastLedgerSequence", optional}, {"LedgerIndex", optional}, {"Hashes", required}}
	ledgerEntryTypeNickName        = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"LedgerEntryType", required}, {"Flags", required}, {"LedgerIndex", optional}, {"MinimumOffer", optional}, {"Account", required}}
	ledgerEntryTypeOffer           = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"LedgerEntryType", required}, {"Flags", required}, {"Sequence", required}, {"PreviousTxnLgrSeq", required}, {"Expiration", optional}, {"BookNode", required}, {"OwnerNode", required}, {"PreviousTxnID", required}, {"LedgerIndex", optional}, {"BookDirectory", required}, {"TakerPays", required}, {"TakerGets", required}, {"Account", required}}
	ledgerEntryTypeSkyWellState    = [][]interface{}{{"LedgerIndex", optional}, {"LedgerEntryType", required}, {"Flags", required}, {"LedgerEntryType", required}, {"Flags", required}, {"PreviousTxnLgrSeq", required}, {"HighQualityIn", optional}, {"HighQualityOut", optional}, {"LowQualityIn", optional}, {"LowQualityOut", optional}, {"LowNode", optional}, {"HighNode", optional}, {"PreviousTxnID", required}, {"LedgerIndex", optional}, {"Balance", required}, {"LowLimit", required}, {"HighLimit", required}}
	ledgerEntryTypes               = map[uint8][][]interface{}{97: ledgerEntryTypeAccountRoot, 99: ledgerEntryTypeContract, 100: ledgerEntryTypeDirectoryNode, 102: ledgerEntryTypeEnabledFeatures, 115: ledgerEntryTypeFeeSettings, 103: ledgerEntryTypeGeneratorMap, 104: ledgerEntryTypeLedgerHashes, 110: ledgerEntryTypeNickName, 111: ledgerEntryTypeOffer, 114: ledgerEntryTypeSkyWellState}

	metaData = [][]interface{}{{"TransactionIndex", required}, {"TransactionResult", required}, {"AffectedNodes", required}}
)

//MemoInfo 备注
type MemoInfo struct {
	Memo *MemoDataInfo
}

//MemoDataInfo 备注
type MemoDataInfo struct {
	MemoData   string
	MemoFormat string
	MemoType   string
}

// type TxData struct {
// 	Flags           uint32
// 	Fee             interface{}
// 	Account         string
// 	TransactionType interface{}
// 	SendMax         interface{}
// 	Memos           []MemoInfo
// 	Paths           [][]PathComputed
// 	TransferRate    uint32
// 	MemoLen         interface{}
// 	Sequence        uint32
// 	Blob            string
// 	TxAmount        interface{}
// 	TakerPays       interface{}
// 	TakerGets       interface{}
// 	SigningPubKey   string
// }

//FromJSON 交易数据序列化。
func FromJSON(txData map[string]interface{}) (*Serializer, error) {
	var typedef [][]interface{}

	txType, ok := txData["TransactionType"]
	if ok {
		if _, ok := txType.(uint8); ok {

		} else if txTypeStr, ok := txType.(string); ok {
			typeInt, ok := txTypeStrMapNumber[txTypeStr]
			if !ok {
				return nil, fmt.Errorf("TransactionType (%s) invalid", txTypeStr)
			}

			typedef = transactionTypes[typeInt]

			txData["TransactionType"] = strconv.Itoa(int(typeInt))
		}
	}

	if len(typedef) == 0 {
		return nil, fmt.Errorf("Object to be serialized must contain either TransactionType, LedgerEntryType or AffectedNodes")
	}

	so := new(Serializer)
	so.Serialize(typedef, txData)

	if so.err != nil {
		return nil, so.err
	}

	return so, nil
}

//Serialize Object 序列化。
func (so *Serializer) Serialize(typedef [][]interface{}, txData map[string]interface{}) {
	STObject.Serialize(so, txData, true)
}

//Append Buffer append
func (so *Serializer) Append(v []byte) {
	so.Buffer = append(so.Buffer, v...)
	// fmt.Println(so.Buffer)
}

//Hash 序列化哈希。
func (so *Serializer) Hash(prefix uint32) []byte {
	sotemp := new(Serializer)
	STInt32.Serialize(sotemp, prefix, false)
	sotemp.Buffer = append(sotemp.Buffer, so.Buffer...)
	sh512 := jtUtils.NewSha512()
	sh512.Add(sotemp.Buffer)

	return sh512.Finish256() //jtUtils.ByteToHexString(sh512.Finish256())
}

//ToHex 序列化转 16 进制。
func (so *Serializer) ToHex() string {
	return hex.EncodeToString(so.Buffer)
}
