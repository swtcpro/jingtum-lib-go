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
	"errors"
	"fmt"
	"strconv"

	jtUtils "jingtumLib/utils"
)

type Serializer struct {
	Buffer []byte
	err    error
}

type SerializedInt8 struct {
}

var (
	REQUIRED = 0
	OPTIONAL = 1
	DEFAULT  = 2

	//交易类型
	TRANSACTION_TYPE_ACCOUNT_SET     = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"EmailHash", OPTIONAL}, {"WalletLocator", OPTIONAL}, {"WalletSize", OPTIONAL}, {"MessageKey", OPTIONAL}, {"Domain", OPTIONAL}, {"TransferRate", OPTIONAL}}
	TRANSACTION_TYPE_TRUST_SET       = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"LimitAmount", OPTIONAL}, {"QualityIn", OPTIONAL}, {"QualityOut", OPTIONAL}}
	TRANSACTION_TYPE_OFFER_CREATE    = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"TakerPays", REQUIRED}, {"TakerGets", REQUIRED}, {"Expiration", OPTIONAL}}
	TRANSACTION_TYPE_OFFER_CANCEL    = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"OfferSequence", REQUIRED}}
	TRANSACTION_TYPE_SET_REGULARKEY  = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"RegularKey", REQUIRED}}
	TRANSACTION_TYPE_PAYMENT         = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Destination", REQUIRED}, {"Amount", REQUIRED}, {"SendMax", OPTIONAL}, {"Paths", DEFAULT}, {"InvoiceID", OPTIONAL}, {"DestinationTag", OPTIONAL}}
	TRANSACTION_TYPE_CONTRACT        = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Expiration", REQUIRED}, {"BondAmount", REQUIRED}, {"StampEscrow", REQUIRED}, {"JingtumEscrow", REQUIRED}, {"CreateCode", OPTIONAL}, {"FundCode", OPTIONAL}, {"RemoveCode", OPTIONAL}, {"ExpireCode", OPTIONAL}}
	TRANSACTION_TYPE_REMOVE_CONTRACT = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Target", REQUIRED}}
	TRANSACTION_TYPE_ENABLE_FEATURE  = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Feature", REQUIRED}}
	TRANSACTION_TYPE_SET_FEE         = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Features", REQUIRED}, {"BaseFee", REQUIRED}, {"ReferenceFeeUnits", REQUIRED}, {"ReserveBase", REQUIRED}, {"ReserveIncrement", REQUIRED}}
	TRANSACTION_TYPE_CONFIG_CONTRACT = [][]interface{}{{"TransactionType", REQUIRED}, {"Flags", OPTIONAL}, {"SourceTag", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"Account", REQUIRED}, {"Sequence", OPTIONAL}, {"Fee", REQUIRED}, {"OperationLimit", OPTIONAL}, {"SigningPubKey", OPTIONAL}, {"TxnSignature", OPTIONAL}, {"Method", REQUIRED}, {"Payload", OPTIONAL}, {"Destination", OPTIONAL}, {"Amount", OPTIONAL}, {"Contracttype", OPTIONAL}, {"ContractMethod", OPTIONAL}, {"Args", OPTIONAL}}
	TRANSACTION_TYPES                = map[uint8][][]interface{}{3: TRANSACTION_TYPE_ACCOUNT_SET, 20: TRANSACTION_TYPE_TRUST_SET, 7: TRANSACTION_TYPE_OFFER_CREATE, 8: TRANSACTION_TYPE_OFFER_CANCEL, 5: TRANSACTION_TYPE_SET_REGULARKEY, 0: TRANSACTION_TYPE_PAYMENT, 9: TRANSACTION_TYPE_CONTRACT, 10: TRANSACTION_TYPE_REMOVE_CONTRACT, 100: TRANSACTION_TYPE_ENABLE_FEATURE, 101: TRANSACTION_TYPE_SET_FEE, 30: TRANSACTION_TYPE_CONFIG_CONTRACT}
	//交易类型和数字映射
	TX_TYPE_STR_MAP_NUMBER = map[string]uint8{"AccountSet": 3, "TrustSet": 20, "OfferCreate": 7, "OfferCancel": 8, "SetRegularKey": 5, "Payment": 0, "Contract": 9, "RemoveContract": 10, "EnableFeature": 100, "SetFee": 101, "ConfigContract": 30}

	//账本类
	LEDGER_ENTRY_TYPE_ACCOUNT_ROOT     = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"Sequence", REQUIRED}, {"PreviousTxnLgrSeq", REQUIRED}, {"TransferRate", OPTIONAL}, {"WalletSize", OPTIONAL}, {"OwnerCount", REQUIRED}, {"EmailHash", OPTIONAL}, {"PreviousTxnID", REQUIRED}, {"AccountTxnID", OPTIONAL}, {"WalletLocator", OPTIONAL}, {"Balance", REQUIRED}, {"MessageKey", OPTIONAL}, {"Domain", OPTIONAL}, {"Account", REQUIRED}, {"RegularKey", OPTIONAL}}
	LEDGER_ENTRY_TYPE_CONTRACT         = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"PreviousTxnLgrSeq", REQUIRED}, {"Expiration", REQUIRED}, {"BondAmount", REQUIRED}, {"PreviousTxnID", REQUIRED}, {"Balance", REQUIRED}, {"FundCode", OPTIONAL}, {"RemoveCode", OPTIONAL}, {"ExpireCode", OPTIONAL}, {"CreateCode", OPTIONAL}, {"Account", REQUIRED}, {"Owner", REQUIRED}, {"Issuer", REQUIRED}}
	LEDGER_ENTRY_TYPE_DIRECTORY_NODE   = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"IndexNext", OPTIONAL}, {"IndexPrevious", OPTIONAL}, {"ExchangeRate", OPTIONAL}, {"RootIndex", REQUIRED}, {"Owner", OPTIONAL}, {"TakerPaysCurrency", OPTIONAL}, {"TakerPaysIssuer", OPTIONAL}, {"TakerGetsCurrency", OPTIONAL}, {"TakerGetsIssuer", OPTIONAL}, {"Indexes", REQUIRED}}
	LEDGER_ENTRY_TYPE_ENABLED_FEATURES = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"Features", REQUIRED}}
	LEDGER_ENTRY_TYPE_FEE_SETTINGS     = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"ReferenceFeeUnits", REQUIRED}, {"ReserveBase", REQUIRED}, {"ReserveIncrement", REQUIRED}, {"BaseFee", REQUIRED}, {"LedgerIndex", OPTIONAL}}
	LEDGER_ENTRY_TYPE_GENERATOR_MAP    = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"Generator", REQUIRED}}
	LEDGER_ENTRY_TYPE_LEDGER_HASHES    = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"FirstLedgerSequence", OPTIONAL}, {"LastLedgerSequence", OPTIONAL}, {"LedgerIndex", OPTIONAL}, {"Hashes", REQUIRED}}
	LEDGER_ENTRY_TYPE_NICKNAME         = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"LedgerIndex", OPTIONAL}, {"MinimumOffer", OPTIONAL}, {"Account", REQUIRED}}
	LEDGER_ENTRY_TYPE_OFFER            = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"Sequence", REQUIRED}, {"PreviousTxnLgrSeq", REQUIRED}, {"Expiration", OPTIONAL}, {"BookNode", REQUIRED}, {"OwnerNode", REQUIRED}, {"PreviousTxnID", REQUIRED}, {"LedgerIndex", OPTIONAL}, {"BookDirectory", REQUIRED}, {"TakerPays", REQUIRED}, {"TakerGets", REQUIRED}, {"Account", REQUIRED}}
	LEDGER_ENTRY_TYPE_SKYWELL_STATE    = [][]interface{}{{"LedgerIndex", OPTIONAL}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"LedgerEntryType", REQUIRED}, {"Flags", REQUIRED}, {"PreviousTxnLgrSeq", REQUIRED}, {"HighQualityIn", OPTIONAL}, {"HighQualityOut", OPTIONAL}, {"LowQualityIn", OPTIONAL}, {"LowQualityOut", OPTIONAL}, {"LowNode", OPTIONAL}, {"HighNode", OPTIONAL}, {"PreviousTxnID", REQUIRED}, {"LedgerIndex", OPTIONAL}, {"Balance", REQUIRED}, {"LowLimit", REQUIRED}, {"HighLimit", REQUIRED}}
	LEDGER_ENTRY_TYPES                 = map[uint8][][]interface{}{97: LEDGER_ENTRY_TYPE_ACCOUNT_ROOT, 99: LEDGER_ENTRY_TYPE_CONTRACT, 100: LEDGER_ENTRY_TYPE_DIRECTORY_NODE, 102: LEDGER_ENTRY_TYPE_ENABLED_FEATURES, 115: LEDGER_ENTRY_TYPE_FEE_SETTINGS, 103: LEDGER_ENTRY_TYPE_GENERATOR_MAP, 104: LEDGER_ENTRY_TYPE_LEDGER_HASHES, 110: LEDGER_ENTRY_TYPE_NICKNAME, 111: LEDGER_ENTRY_TYPE_OFFER, 114: LEDGER_ENTRY_TYPE_SKYWELL_STATE}

	//元数据
	METADATA = [][]interface{}{{"TransactionIndex", REQUIRED}, {"TransactionResult", REQUIRED}, {"AffectedNodes", REQUIRED}}
)

type MemoInfo struct {
	Memo *MemoDataInfo
}

type MemoDataInfo struct {
	MemoData   string
	MemoFormat string
	MemoType   string
}

type TxData struct {
	Flags           uint32
	Fee             interface{}
	Account         string
	TransactionType interface{}
	SendMax         interface{}
	Memos           []MemoInfo
	Paths           [][]PathComputed
	TransferRate    uint32
	MemoLen         interface{}
	Sequence        uint32
	Blob            string
	TxAmount        interface{}
	TakerPays       interface{}
	TakerGets       interface{}
	SigningPubKey   string
}

/**
 *  对象转字节序列化
 */
func FromJson(txData map[string]interface{}) (*Serializer, error) {
	var typedef [][]interface{}

	txType, ok := txData["TransactionType"]
	if ok {
		if _, ok := txType.(uint8); ok {

		} else if txTypeStr, ok := txType.(string); ok {
			typeInt, ok := TX_TYPE_STR_MAP_NUMBER[txTypeStr]
			if !ok {
				return nil, errors.New(fmt.Sprintf("TransactionType (%s) invalid.", txTypeStr))
			}

			typedef = TRANSACTION_TYPES[typeInt]

			txData["TransactionType"] = strconv.Itoa(int(typeInt))
		}
	}

	if len(typedef) == 0 {
		return nil, errors.New("Object to be serialized must contain either TransactionType, LedgerEntryType or AffectedNodes.")
	}

	so := new(Serializer)
	so.Serialize(typedef, txData)

	if so.err != nil {
		return nil, so.err
	}

	return so, nil
}

func (so *Serializer) Serialize(typedef [][]interface{}, txData map[string]interface{}) {
	STObject.Serialize(so, txData, true)
}

func (so *Serializer) Append(v []byte) {
	so.Buffer = append(so.Buffer, v...)
}

func (so *Serializer) Hash(prefix uint32) []byte {
	sotemp := new(Serializer)
	STInt32.Serialize(sotemp, prefix, false)
	sotemp.Buffer = append(sotemp.Buffer, so.Buffer...)
	sh512 := jtUtils.NewSha512()
	sh512.Add(sotemp.Buffer)
	return sh512.Finish256()
}

func (so *Serializer) ToHex() string {
	return hex.EncodeToString(so.Buffer)
}
