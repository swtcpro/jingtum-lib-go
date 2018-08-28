/**
 * Transaction类主管POST请求，包括组装交易和交易参数。请求时需要提供密钥，且交易可以
 * 进行本地签名和服务器签名。目前支持服务器签名，本地签名支持主要的交易，还有部分参数
 * 不支持。所有的请求是异步的，会提供一个回调函数。每个回调函数有两个参数，一个是错误，
 * 另一个是结果。
 *
 * @FileName: transaction.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-05-28 10:44:32
 * @UpdateTime: 2018-05-28 10:44:54
 */

package jingtumlib

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"jingtumlib/constant"
	"jingtumlib/serializer"
	"jingtumlib/utils"

	"github.com/yangxuebo-138/decimal"
)

//Transaction 交易结构体
type Transaction struct {
	remote    *Remote
	txJSON    map[string]interface{}
	localSign bool
	secret    string
	filter    Filter
}

//FlagClass FlagClass
type FlagClass map[string]uint32

//AccountSet AccountSet
type AccountSet map[string]uint32

var (
	//TransactionFlags 交易标识
	TransactionFlags = map[string]FlagClass{"Universal": {"FullyCanonicalSig": 0x00010000}, "AccountSet": {"RequireDestTag": 0x00010000, "OptionalDestTag": 0x00020000, "RequireAuth": 0x00040000, "OptionalAuth": 0x00080000, "DisallowSWT": 0x00100000, "AllowSWT": 0x00200000}, "TrustSet": {"SetAuth": 0x00010000, "NoSkywell": 0x00020000, "SetNoSkywell": 0x00020000, "ClearNoSkywell": 0x00040000, "SetFreeze": 0x00100000, "ClearFreeze": 0x00200000}, "OfferCreate": {"Passive": 0x00010000, "ImmediateOrCancel": 0x00020000, "FillOrKill": 0x00040000, "Sell": 0x00080000}, "Payment": {"NoSkywellDirect": 0x00010000, "PartialPayment": 0x00020000, "LimitQuality": 0x00040000}, "RelationSet": {"Authorize": 0x00000001, "Freeze": 0x00000011}, "RelationDel": {}}
	//SetClearFlags 清除标识
	SetClearFlags = map[uint32]AccountSet{uint32(1): {"asfRequireDest": uint32(1), "asfRequireAuth": uint32(2), "asfDisallowSWT": uint32(3), "asfDisableMaster": uint32(4), "asfNoFreeze": uint32(5), "asfGlobalFreeze": uint32(6)}}
)

//OfferTypes offer type 映射
var OfferTypes = map[string]int{"Sell": 1, "Buy": 2}

//RelationTypes relation type 映射
var RelationTypes = map[string]int{"trust": 1, "authorize": 2, "freeze": 3, "unfreeze": 4}

//AccountSetTypes account type 映射
var AccountSetTypes = map[string]int{"property": 1, "delegate": 2, "signer": 3}

//NewTransaction 构造Transaction对象
func NewTransaction(remote *Remote, filter Filter) (*Transaction, error) {
	if nil == remote {
		return nil, constant.ERR_EMPTY_PARAM
	}
	tx := new(Transaction)
	tx.remote = remote
	tx.txJSON = make(map[string]interface{})
	tx.AddTxJSON("Flags", uint32(0))
	tx.AddTxJSON("Fee", float32(JTConfig.ReadInt("Config", "fee", 10000)))
	if filter == nil {
		filter = func(data interface{}) interface{} {
			return data
		}
	}
	tx.filter = filter
	return tx, nil
}

//func (tx *Transaction)

//AddTxJSON 添加交易参数
func (tx *Transaction) AddTxJSON(key string, value interface{}) bool {
	if key == "" || tx.txJSON == nil {
		return false
	}

	tx.txJSON[key] = value
	return true
}

//GetTxJSON 获取交易参数
func (tx *Transaction) GetTxJSON(key string) interface{} {
	v, ok := tx.txJSON[key]
	if ok {
		return v
	}

	return nil
}

//SetSecret 本地签名时需要设置私钥
func (tx *Transaction) SetSecret(secret string) {
	if !IsValidSecret(secret) {
		tx.AddTxJSON(constant.TxJSONErrorKey, constant.ERR_PAYMENT_INVALID_SECRET)
		return
	}

	tx.secret = secret
}

//AddMemo 设置备注
func (tx *Transaction) AddMemo(memo string) {
	if memo == "" {
		tx.AddTxJSON(constant.TxJSONErrorKey, constant.ERR_PAYMENT_MEMO_EMPTY)
		return
	}

	if len([]rune(memo)) > 2048 {
		tx.AddTxJSON(constant.TxJSONErrorKey, constant.ERR_PAYMENT_OUT_OF_MEMO_LEN)
	}

	mi := new(serializer.MemoInfo)
	mdi := new(serializer.MemoDataInfo)
	mdi.MemoData = utils.StringToHex(memo)
	mi.Memo = mdi

	if nil == tx.txJSON["Memos"] {
		memos := list.New()
		memos.PushBack(mi)
		tx.AddTxJSON("Memos", memos)
	} else {
		memos := tx.txJSON["Memos"].(*list.List)
		memos.PushBack(mi)
	}
}

//GetAccount 获得交易账号
func (tx *Transaction) GetAccount() string {
	account, _ := tx.txJSON["Account"].(string)
	return account
}

//GetTransactionType 获得交易类型
func (tx *Transaction) GetTransactionType() string {
	txType, _ := tx.txJSON["TransactionType"].(string)

	return txType
}

func (tx *Transaction) setFee(fee float32) {
	if fee < 10 {
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("Fee should be great than or equal 10")
		return
	}

	tx.txJSON["Fee"] = fee
}

func maxAmount(amount interface{}) (interface{}, error) {
	if amt, ok := amount.(string); ok {
		if utils.IsNumberString(amt) {
			f, err := strconv.ParseFloat(amt, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid amount to max %s", amt)
			}

			return strconv.FormatFloat(f*1.0001, 'f', 10, 32), nil
		}
	}

	if amt, ok := amount.(constant.Amount); ok && utils.IsValidAmount(&amt) {
		f, err := strconv.ParseFloat(amt.Value, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid amount to max %s", amt.Value)
		}
		amt.Value = strconv.FormatFloat(f*1.0001, 'f', 10, 32)
		return amt, nil
	}

	return nil, fmt.Errorf("invalid amount to max")
}

//SetPath 设置支付路径
func (tx *Transaction) SetPath(key string) {
	if key == "" || len(key) != 40 {
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("invalid path key")
		return
	}

	item, ok := tx.remote.Paths.Get(key)

	if !ok {
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("non exists path key %s", key)
		return
	}

	//沒有支付路径，不需要传下面的参数
	if pd, ok := item.(serializer.PathData); ok {
		if len(pd.PathsComputed) == 0 {
			return
		}

		tx.txJSON["Paths"] = pd.PathsComputed
		amount, err := maxAmount(pd.Choice)
		if err != nil {
			tx.txJSON[constant.TxJSONErrorKey] = err
		}

		tx.txJSON["SendMax"] = amount
	}
}

//SetSendMax limit send max amount
func (tx *Transaction) SetSendMax(amount constant.Amount) {
	if !utils.IsValidAmount(&amount) {
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("invalid send max amount")
		return
	}

	tx.txJSON["SendMax"] = amount
}

//SetTransferRate 设置手续费汇率
func (tx *Transaction) SetTransferRate(rate float32) {
	if rate < 0 || rate > 1 {
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("invalid transfer rate %.f", rate)
		return
	}
	tx.txJSON["TransferRate"] = (rate + 1) * 1e9
}

//SetFlags flags 数据类型：uint32或[]string
func (tx *Transaction) SetFlags(flags interface{}) {
	switch v := flags.(type) {
	case uint32:
		tx.txJSON["Flags"] = v
	case []string:
		if txFlags, ok := TransactionFlags[tx.GetTransactionType()]; ok {
			flgTemp := uint32(0)
			for _, flgStr := range v {
				if flgv, ok := txFlags[flgStr]; ok {
					flgTemp += uint32(flgv)
				}
			}
			tx.txJSON["Flags"] = flgTemp
		}
	default:
		tx.txJSON[constant.TxJSONErrorKey] = fmt.Errorf("invalid flags")
	}
}

//signing 签名
func signing(tx *Transaction) (string, error) {
	fee, ok := decimal.NewFromFloat32(tx.GetTxJSON("Fee").(float32)).Div(decimal.NewFromFloat32(1000000)).Float64()
	if !ok {
		return "", fmt.Errorf("Fee / 1000000 float error")
	}
	tx.AddTxJSON("Fee", float32(fee))

	amount := tx.GetTxJSON("Amount")
	if amount != nil {
		if amt64, ok := amount.(float64); ok {
			amt, ok := decimal.NewFromFloat(amt64).Div(decimal.NewFromFloat(1000000)).Float64() //	NewFromFloat32(tx.GetTxJson("Fee").(float32)).Div(decimal.NewFromFloat32(1000000)).Float64()
			if !ok {
				return "", fmt.Errorf("Amount / 1000000 float error")
			}
			tx.AddTxJSON("Amount", amt)
		}
	}

	if tx.GetTxJSON("Memos") != nil {
		memos := tx.GetTxJSON("Memos").(*list.List)
		for e := memos.Front(); e != nil; e = e.Next() {
			e.Value.(*serializer.MemoInfo).Memo.MemoData, _ = utils.HexToString(e.Value.(*serializer.MemoInfo).Memo.MemoData)
		}
	}

	if tx.GetTxJSON("SendMax") != nil {
		if sendMax, ok := tx.GetTxJSON("SendMax").(float64); ok {
			sm, ok := decimal.NewFromFloat(sendMax).Div(decimal.NewFromFloat(1000000)).Float64()
			if !ok {
				return "", fmt.Errorf("SendMax / 1000000 float error")
			}
			tx.AddTxJSON("SendMax", sm)
		}
	}

	if tx.GetTxJSON("TakerPays") != nil {
		if takerPays, ok := tx.GetTxJSON("TakerPays").(float64); ok {
			tp, ok := decimal.NewFromFloat(takerPays).Div(decimal.NewFromFloat(1000000)).Float64()
			if !ok {
				return "", fmt.Errorf("TakerPays / 1000000 float error")
			}
			tx.AddTxJSON("TakerPays", tp)
		}
	}

	if tx.GetTxJSON("TakerGets") != nil {
		if takerGets, ok := tx.GetTxJSON("TakerGets").(float64); ok {
			tg, ok := decimal.NewFromFloat(takerGets).Div(decimal.NewFromFloat(1000000)).Float64()
			if !ok {
				return "", fmt.Errorf("TakerGets / 1000000 float error")
			}
			tx.AddTxJSON("TakerGets", tg)
		}
	}

	wt, err := FromSecret(tx.secret)
	if err != nil {
		return "", err
	}

	tx.AddTxJSON("SigningPubKey", wt.GetPublicKey())
	var prefix uint32 = 0x53545800
	so, err := serializer.FromJSON(utils.DeepCopy(tx.txJSON).(map[string]interface{}))
	if err != nil {
		return "", err
	}
	hash := so.Hash(prefix)
	// fmt.Println(utils.ByteToHexString(hash))
	signTx, err := wt.signTx(hash)
	if err != nil {
		return "", err
	}

	tx.AddTxJSON("TxnSignature", signTx)
	// fmt.Println(signTx)
	soBlog, err := serializer.FromJSON(utils.DeepCopy(tx.txJSON).(map[string]interface{}))

	if err != nil {
		return "", err
	}
	// fmt.Println(strings.ToUpper(soBlog.ToHex()))
	tx.AddTxJSON("blob", strings.ToUpper(soBlog.ToHex()))
	tx.localSign = true
	return tx.GetTxJSON("blob").(string), nil
}

//sign 签名方法
func (tx *Transaction) sign(callback func(err error, blob string)) {

	if tx.GetTxJSON("Sequence") == nil {
		//从服务端获取 Sequence 后再签名
		options := make(map[string]interface{})
		options["account"] = tx.GetTxJSON("Account")
		options["type"] = "trust"
		req, err := tx.remote.RequestAccountInfo(options)
		if err != nil {
			callback(err, "")
			return
		}
		req.Submit(func(err error, result interface{}) {
			if err != nil {
				callback(err, "")
				return
			}

			if ret, ok := result.(map[string]interface{}); ok {

				actData, ok := ret["account_data"]
				if !ok {
					callback(fmt.Errorf("account_data is null"), "")
					return
				}

				actDataMap, ok := actData.(map[string]interface{})
				if !ok {
					callback(fmt.Errorf("account_data type %T error", actData), "")
					return
				}
				seq, ok := actDataMap["Sequence"]
				if !ok {
					callback(fmt.Errorf("Get Sequence is null from server"), "")
					return
				}

				tx.AddTxJSON("Sequence", uint32(decimal.NewFromFloat(seq.(float64)).IntPart()))
				blob, err := signing(tx)
				if err != nil {
					callback(err, "")
				} else {
					callback(nil, blob)
				}
			} else {
				callback(fmt.Errorf("Request account info fail"), "")
			}

		})
	} else {
		blob, err := signing(tx)
		if err != nil {
			callback(err, "")
		} else {
			callback(nil, blob)
		}
	}
}

//Submit 提交交易数据
func (tx *Transaction) Submit(callback func(err error, result interface{})) {
	if !tx.remote.server.IsConnected() {
		callback(fmt.Errorf("Server not connected"), nil)
		return
	}
	if tx.checkTxError() {
		callback(tx.GetTxJSON(constant.TxJSONErrorKey).(error), nil)
		return
	}

	if tx.remote.LocalSign {
		//本地签名
		tx.sign(func(err error, blob string) {
			if nil != err {
				callback(errors.New("sig error. "+err.Error()), nil)
			} else {
				data := map[string]interface{}{"tx_blob": blob}
				tx.remote.Submit(constant.CommandSubmit, data, tx.filter, callback)
			}
		})
	} else if tx.GetTxJSON("TransactionType") == "Signer" {
		//直接将blob传给底层
		data := map[string]interface{}{"tx_blob": tx.GetTxJSON("blob")}
		tx.remote.Submit(constant.CommandSubmit, data, tx.filter, callback)
	} else {
		//不签名交易传给底层
		data := map[string]interface{}{"secret": tx.secret, "tx_json": tx.txJSON}
		tx.remote.Submit(constant.CommandSubmit, data, tx.filter, callback)
	}
}

func (tx *Transaction) checkTxError() bool {
	if tx.GetTxJSON(constant.TxJSONErrorKey) != nil {
		return true
	}
	return false
}
