/**
 * 全局变量定义类。
 *
 * @FileName: global.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package constant

import (
	"strconv"
)

//KeyValuePair 字典映射
type KeyValuePair struct {
	Key   int
	Value int
}

//Amount 支付金额
type Amount struct {
	//货币种类，三到六个字母或20字节的自定义货币
	Currency string `json:"currency"`
	//货币发行方
	Issuer string `json:"issuer"`
	//支付数量
	Value string `json:"value"`
}

var (
	//CFGCurrency 配置货币
	CFGCurrency string

	//LedgerStates 账本状态
	LedgerStates = map[string]string{"current": "current", "closed": "closed", "validated": "validated"}

	//InverseFieldsMap 反转字段映射
	InverseFieldsMap = map[string]*KeyValuePair{
		"LedgerEntryType":     &KeyValuePair{1, 1},
		"TransactionType":     &KeyValuePair{1, 2},
		"Flags":               &KeyValuePair{2, 2},
		"SourceTag":           &KeyValuePair{2, 3},
		"Sequence":            &KeyValuePair{2, 4},
		"PreviousTxnLgrSeq":   &KeyValuePair{2, 5},
		"LedgerSequence":      &KeyValuePair{2, 6},
		"CloseTime":           &KeyValuePair{2, 7},
		"ParentCloseTime":     &KeyValuePair{2, 8},
		"SigningTime":         &KeyValuePair{2, 9},
		"Expiration":          &KeyValuePair{2, 10},
		"TransferRate":        &KeyValuePair{2, 11},
		"WalletSize":          &KeyValuePair{2, 12},
		"OwnerCount":          &KeyValuePair{2, 13},
		"DestinationTag":      &KeyValuePair{2, 14},
		"Timestamp":           &KeyValuePair{2, 15},
		"HighQualityIn":       &KeyValuePair{2, 16},
		"HighQualityOut":      &KeyValuePair{2, 17},
		"LowQualityIn":        &KeyValuePair{2, 18},
		"LowQualityOut":       &KeyValuePair{2, 19},
		"QualityIn":           &KeyValuePair{2, 20},
		"QualityOut":          &KeyValuePair{2, 21},
		"StampEscrow":         &KeyValuePair{2, 22},
		"BondAmount":          &KeyValuePair{2, 23},
		"LoadFee":             &KeyValuePair{2, 24},
		"OfferSequence":       &KeyValuePair{2, 25},
		"FirstLedgerSequence": &KeyValuePair{2, 26},
		"LastLedgerSequence":  &KeyValuePair{2, 27},
		"TransactionIndex":    &KeyValuePair{2, 28},
		"OperationLimit":      &KeyValuePair{2, 29},
		"ReferenceFeeUnits":   &KeyValuePair{2, 30},
		"ReserveBase":         &KeyValuePair{2, 31},
		"ReserveIncrement":    &KeyValuePair{2, 32},
		"SetFlag":             &KeyValuePair{2, 33},
		"ClearFlag":           &KeyValuePair{2, 34},
		"RelationType":        &KeyValuePair{2, 35},
		"Method":              &KeyValuePair{2, 36},
		"Contracttype":        &KeyValuePair{2, 39},
		"IndexNext":           &KeyValuePair{3, 1},
		"IndexPrevious":       &KeyValuePair{3, 2},
		"BookNode":            &KeyValuePair{3, 3},
		"OwnerNode":           &KeyValuePair{3, 4},
		"BaseFee":             &KeyValuePair{3, 5},
		"ExchangeRate":        &KeyValuePair{3, 6},
		"LowNode":             &KeyValuePair{3, 7},
		"HighNode":            &KeyValuePair{3, 8},
		"EmailHash":           &KeyValuePair{4, 1},
		"LedgerHash":          &KeyValuePair{5, 1},
		"ParentHash":          &KeyValuePair{5, 2},
		"TransactionHash":     &KeyValuePair{5, 3},
		"AccountHash":         &KeyValuePair{5, 4},
		"PreviousTxnID":       &KeyValuePair{5, 5},
		"LedgerIndex":         &KeyValuePair{5, 6},
		"WalletLocator":       &KeyValuePair{5, 7},
		"RootIndex":           &KeyValuePair{5, 8},
		"AccountTxnID":        &KeyValuePair{5, 9},
		"BookDirectory":       &KeyValuePair{5, 16},
		"InvoiceID":           &KeyValuePair{5, 17},
		"Nickname":            &KeyValuePair{5, 18},
		"Amendment":           &KeyValuePair{5, 19},
		"TicketID":            &KeyValuePair{5, 20},
		"Amount":              &KeyValuePair{6, 1},
		"Balance":             &KeyValuePair{6, 2},
		"LimitAmount":         &KeyValuePair{6, 3},
		"TakerPays":           &KeyValuePair{6, 4},
		"TakerGets":           &KeyValuePair{6, 5},
		"LowLimit":            &KeyValuePair{6, 6},
		"HighLimit":           &KeyValuePair{6, 7},
		"Fee":                 &KeyValuePair{6, 8},
		"SendMax":             &KeyValuePair{6, 9},
		"MinimumOffer":        &KeyValuePair{6, 16},
		"JingtumEscrow":       &KeyValuePair{6, 17},
		"DeliveredAmount":     &KeyValuePair{6, 18},
		"PublicKey":           &KeyValuePair{7, 1},
		"MessageKey":          &KeyValuePair{7, 2},
		"SigningPubKey":       &KeyValuePair{7, 3},
		"TxnSignature":        &KeyValuePair{7, 4},
		"Generator":           &KeyValuePair{7, 5},
		"Signature":           &KeyValuePair{7, 6},
		"Domain":              &KeyValuePair{7, 7},
		"FundCode":            &KeyValuePair{7, 8},
		"RemoveCode":          &KeyValuePair{7, 9},
		"ExpireCode":          &KeyValuePair{7, 10},
		"CreateCode":          &KeyValuePair{7, 11},
		"MemoType":            &KeyValuePair{7, 12},
		"MemoData":            &KeyValuePair{7, 13},
		"MemoFormat":          &KeyValuePair{7, 14},
		"Payload":             &KeyValuePair{7, 15},
		"ContractMethod":      &KeyValuePair{7, 17},
		"Parameter":           &KeyValuePair{7, 18},
		"Account":             &KeyValuePair{8, 1},
		"Owner":               &KeyValuePair{8, 2},
		"Destination":         &KeyValuePair{8, 3},
		"Issuer":              &KeyValuePair{8, 4},
		"Target":              &KeyValuePair{8, 7},
		"RegularKey":          &KeyValuePair{8, 8},
		"undefined":           &KeyValuePair{15, 1},
		"TransactionMetaData": &KeyValuePair{14, 2},
		"CreatedNode":         &KeyValuePair{14, 3},
		"DeletedNode":         &KeyValuePair{14, 4},
		"ModifiedNode":        &KeyValuePair{14, 5},
		"PreviousFields":      &KeyValuePair{14, 6},
		"FinalFields":         &KeyValuePair{14, 7},
		"NewFields":           &KeyValuePair{14, 8},
		"TemplateEntry":       &KeyValuePair{14, 9},
		"Memo":                &KeyValuePair{14, 10},
		"Arg":                 &KeyValuePair{14, 11},
		"SigningAccounts":     &KeyValuePair{15, 2},
		"TxnSignatures":       &KeyValuePair{15, 3},
		"Signatures":          &KeyValuePair{15, 4},
		"Template":            &KeyValuePair{15, 5},
		"Necessary":           &KeyValuePair{15, 6},
		"Sufficient":          &KeyValuePair{15, 7},
		"AffectedNodes":       &KeyValuePair{15, 8},
		"Memos":               &KeyValuePair{15, 9},
		"Args":                &KeyValuePair{15, 10},
		"CloseResolution":     &KeyValuePair{16, 1},
		"TemplateEntryType":   &KeyValuePair{16, 2},
		"TransactionResult":   &KeyValuePair{16, 3},
		"TakerPaysCurrency":   &KeyValuePair{17, 1},
		"TakerPaysIssuer":     &KeyValuePair{17, 2},
		"TakerGetsCurrency":   &KeyValuePair{17, 3},
		"TakerGetsIssuer":     &KeyValuePair{17, 4},
		"Paths":               &KeyValuePair{18, 1},
		"Indexes":             &KeyValuePair{19, 1},
		"Hashes":              &KeyValuePair{19, 2},
		"Amendments":          &KeyValuePair{19, 3}}
)

//Integer int 包装结构。
type Integer struct {
	intv int
}

//ResponseData 区块连网络响应数据结构。
type ResponseData struct {
	ID           uint64                 `json:"id"`
	Status       string                 `json:"status"`
	Type         string                 `json:"type"`
	Result       map[string]interface{} `json:"result"`
	Request      map[string]interface{} `json:"request"`
	Validated    bool                   `json:"validated"`
	LedgerIndex  int                    `json:"ledger_index"`
	LedgerHash   string                 `json:"ledger_hash"`
	ErrorMessage string                 `json:"error_message"`
	ErrorCode    int                    `json:"error_code"`
	Error        string                 `json:"error"`
	Account      string                 `json:"account"`
}

//IntValue 基本类型int
func (integer *Integer) IntValue() int {
	return integer.intv
}

func (integer *Integer) String() string {
	return strconv.Itoa(integer.intv)
}

//NewInteger 创建int的包装类型
func NewInteger(intv int) *Integer {
	integer := new(Integer)
	integer.intv = intv
	return integer
}
