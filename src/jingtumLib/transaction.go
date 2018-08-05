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

package jingtumLib

import (
	"github.com/yangxuebo-138/decimal"
	//	"encoding/json"
	"errors"
	"fmt"

	//	"strconv"
	//	"strings"
	"container/list"

	"jingtumLib/constant"
	"jingtumLib/serializer"
	"jingtumLib/utils"
)

//type Transaction1 struct {
//	remote     *Remote
//	filter     Filter
//	secret     string
//	tx_json    *serializer.TxData
//	local_sign bool
//}

type Transaction struct {
	remote     *Remote
	tx_json    map[string]interface{}
	local_sign bool
	secret     string
}

type FlagClass map[string]uint32

var (
	TransactionFlags = map[string]FlagClass{"Universal": {"FullyCanonicalSig": 0x00010000}, "AccountSet": {"RequireDestTag": 0x00010000, "OptionalDestTag": 0x00020000, "RequireAuth": 0x00040000, "OptionalAuth": 0x00080000, "DisallowSWT": 0x00100000, "AllowSWT": 0x00200000}, "TrustSet": {"SetAuth": 0x00010000, "NoSkywell": 0x00020000, "SetNoSkywell": 0x00020000, "ClearNoSkywell": 0x00040000, "SetFreeze": 0x00100000, "ClearFreeze": 0x00200000}, "OfferCreate": {"Passive": 0x00010000, "ImmediateOrCancel": 0x00020000, "FillOrKill": 0x00040000, "Sell": 0x00080000}, "Payment": {"NoSkywellDirect": 0x00010000, "PartialPayment": 0x00020000, "LimitQuality": 0x00040000}, "RelationSet": {"Authorize": 0x00000001, "Freeze": 0x00000011}}
)

/**
 * 构造Transaction对象
 */
func NewTransaction(remote *Remote) (*Transaction, error) {
	if nil == remote {
		return nil, constant.ERR_EMPTY_PARAM
	}
	tx := new(Transaction)
	tx.remote = remote
	tx.tx_json = make(map[string]interface{})
	tx.AddTxJson("Flags", uint32(0))
	tx.AddTxJson("Fee", JTConfig.ReadInt("Config", "fee", 10000))
	return tx, nil
}

//func (tx *Transaction)

/**
 * 添加交易参数
 */
func (tx *Transaction) AddTxJson(key string, value interface{}) bool {
	if key == "" || tx.tx_json == nil {
		return false
	}

	tx.tx_json[key] = value
	return true
}

func (tx *Transaction) GetTxJson(key string) interface{} {
	v, ok := tx.tx_json[key]
	if ok {
		return v
	}

	return nil
}

/**
 * 本地签名时需要设置私钥
 */
func (tx *Transaction) SetSecret(secret string) {
	if !IsValidSecret(secret) {
		tx.AddTxJson(constant.TXJSON_ERROR_KEY, constant.ERR_PAYMENT_INVALID_SECRET)
		return
	}

	tx.secret = secret
}

/**
 * 设置备注
 */
func (tx *Transaction) AddMemo(memo string) {
	if memo == "" {
		tx.AddTxJson(constant.TXJSON_ERROR_KEY, constant.ERR_PAYMENT_MEMO_EMPTY)
		return
	}

	if len([]rune(memo)) > 2048 {
		tx.AddTxJson(constant.TXJSON_ERROR_KEY, constant.ERR_PAYMENT_OUT_OF_MEMO_LEN)
	}

	mi := new(serializer.MemoInfo)
	mdi := new(serializer.MemoDataInfo)
	mdi.MemoData = utils.StringToHex(memo)
	mi.Memo = mdi

	if nil == tx.tx_json["Memos"] {
		memos := list.New()
		memos.PushBack(mi)
		tx.AddTxJson("Memos", memos)
	} else {
		memos := tx.tx_json["Memos"].(*list.List)
		memos.PushBack(mi)
	}
}

//func (transaction *Transaction) ParseJson(jsonStr string) error {
//
//	err := json.Unmarshal([]byte(jsonStr), &transaction.tx_json)
//	return err
//}
//
//func (transaction *Transaction) GetAccount() string {
//	return transaction.tx_json.Account
//}
//
//func (transaction *Transaction) GetTransactionType() interface{} {
//	return transaction.tx_json.TransactionType
//}
//
//func (transaction *Transaction) SetSecret(secret string) {
//	transaction.secret = secret
//}
//
//
//func (transaction *Transaction) SetFee(fee uint32) {
//	if fee < 10 {
//		transaction.tx_json.Fee = errors.New("Fee should be great than or equal 10.")
//		return
//	}
//
//	transaction.tx_json.Fee = fee
//}
//
//func (transaction *Transaction) MaxAmount(amount interface{}) interface{} {
//	if mt, ok := amount.(string); ok {
//		if utils.IsNumberString(mt) {
//			f, err := strconv.ParseFloat(mt, 32)
//			if err != nil {
//				return errors.New("invalid amount to max")
//			}
//
//			return strconv.FormatFloat(f*1.0001, 'f', 10, 32)
//		}
//	}
//
//	if at, ok := amount.(constant.Amount); ok && utils.IsValidAmount(at) {
//		f, err := strconv.ParseFloat(at.Value, 32)
//		if err != nil {
//			return errors.New("invalid amount to max")
//		}
//		return strconv.FormatFloat(f*1.0001, 'f', 10, 32)
//	}
//
//	return errors.New("invalid amount to max")
//}
//
///**
// * set a path to payment
// * this path is repesented as a key, which is computed in path find
// * so if one path you computed self is not allowed
// * when path set, sendmax is also set.
// * @param path
// */
//func (transaction *Transaction) setPath(key string) (err error) {
//	if key == "" || (strings.Count(key, "")-1) != 40 {
//		err = errors.New("invalid path key")
//		return
//	}
//
//	item, ok := transaction.remote.Paths.Get(key)
//
//	if !ok {
//		err = errors.New("non exists path key")
//		return
//	}
//
//	//沒有支付路径，不需要传下面的参数
//	if item.(jtSerz.PathData).PathsComputed == nil || len(item.(jtSerz.PathData).PathsComputed) == 0 {
//		return
//	}
//	//var path [][]interface{}
//	//err = json.Unmarshal(item.(jtSerz.PathData).Pathcomputed, &path)
//	//if err != nil {
//	//err = errors.New("invalid path.")
//	//return
//	//}
//
//	transaction.tx_json.Paths = item.(jtSerz.PathData).PathsComputed
//	amount := transaction.MaxAmount(item.(jtSerz.PathData).Choice)
//	transaction.tx_json.SendMax = amount
//	return
//}
//
///**
// * limit send max amount
// */
//func (transaction *Transaction) setSendMax(amount constant.Amount) {
//	if !utils.IsValidAmount(amount) {
//		transaction.tx_json.SendMax = errors.New("invalid send max amount")
//		return
//	}
//
//	transaction.tx_json.SendMax = amount
//}
//
///**
// * transfer rate
// * between 0 and 1, type is number
// * @param rate
// */
//func (transaction *Transaction) setTransferRate(rate float32) (err error) {
//	if rate < 0 || rate > 1 {
//		err = errors.New("invalid transfer rate")
//		return
//	}
//
//	transaction.tx_json.TransferRate = uint32((rate + 1) * 1000000000)
//	return
//}
//
///**
// * set transaction flags
// *
// */
//func (transaction *Transaction) setFlags(flags interface{}) (err error) {
//	if fv, ok := flags.(uint32); ok {
//		transaction.tx_json.Flags = fv
//		return
//	}
//
//	if transType, isString := transaction.tx_json.TransactionType.(string); isString {
//		var transaction_flags = TransactionFlags[transType]
//		if transaction_flags != nil {
//			if flag_set, isArray := flags.([]string); isArray {
//				for i := 0; i < len(flag_set); i++ {
//					flag := flag_set[i]
//					transaction.tx_json.Flags += transaction_flags[flag]
//				}
//			}
//		}
//	}
//	err = errors.New("invalid flags")
//	return
//}
/**
 * 签名方法
 */
func signing(tx *Transaction) (string, error) {

	tx.AddTxJson("Fee", tx.GetTxJson("Fee").(int)/1000000)

	amount := tx.GetTxJson("Amount")
	if amount == nil {
		return "", errors.New("Amount not be empty.")
	}

	if amt64, ok := amount.(float64); ok {
		tx.AddTxJson("Amount", amt64/1000000)
	}

	if tx.GetTxJson("Memos") != nil {
		memos := tx.GetTxJson("Memos").(*list.List)
		for e := memos.Front(); e != nil; e = e.Next() {
			e.Value.(*serializer.MemoInfo).Memo.MemoData, _ = utils.HexToString(e.Value.(*serializer.MemoInfo).Memo.MemoData)
		}
	}

	if tx.GetTxJson("SendMax") != nil {
		if sendMax, ok := tx.GetTxJson("SendMax").(float64); ok {
			tx.AddTxJson("SendMax", sendMax/1000000)
		}
	}

	if tx.GetTxJson("TakerPays") != nil {
		if takerPays, ok := tx.GetTxJson("TakerPays").(float64); ok {
			tx.AddTxJson("TakerPays", takerPays/1000000)
		}
	}

	if tx.GetTxJson("TakerGets") != nil {
		if takerGets, ok := tx.GetTxJson("TakerGets").(float64); ok {
			tx.AddTxJson("TakerGets", takerGets/1000000)
		}
	}

	wt, err := FromSecret(tx.secret)
	if err != nil {
		return "", err
	}

	tx.AddTxJson("SigningPubKey", wt.GetPublicKey())
	var prefix uint32 = 0x53545800
	so, err := serializer.FromJSON(tx.tx_json)
	if err != nil {
		return "", err
	}
	hash := so.Hash(prefix)
	fmt.Println(hash)

	//transaction.tx_json.TxnSignature = wt.signTx(hash)
	//transaction.tx_json.Blob = jtSerz.FromJson(transaction.tx_json).ToHex()
	tx.local_sign = true
	//
	//    var wt = new base(self._secret);
	//    self.tx_json.SigningPubKey = wt.getPublicKey();
	//    var prefix = 0x53545800;
	//    var hash = jser.from_json(self.tx_json).hash(prefix);
	//    self.tx_json.TxnSignature = wt.signTx(hash);
	//    self.tx_json.blob =  jser.from_json(self.tx_json).to_hex();
	//    self._local_sign = true;
	//    return self.tx_json.blob;
	return "----", nil
}

func (tx *Transaction) sign(callback func(err error, blob string)) {

	if tx.GetTxJson("Sequence") == nil {
		//从服务端获取 Sequence 后再签名
		options := make(map[string]interface{})
		options["account"] = tx.GetTxJson("Account")
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

			fmt.Printf("Account info %v\n", result)

			if ret, ok := result.(map[string]interface{}); ok {
				fmt.Printf("Account info %v\n", ret)
				actData, ok := ret["account_data"]
				if !ok {
					callback(errors.New("account_data is null."), "")
					return
				}

				actDataMap, ok := actData.(map[string]interface{})
				if !ok {
					callback(errors.New(fmt.Sprintf("account_data type %T error.", actData)), "")
					return
				}
				seq, ok := actDataMap["Sequence"]
				if !ok {
					callback(errors.New("Get Sequence is null from server."), "")
					return
				}

				tx.AddTxJson("Sequence", uint32(decimal.NewFromFloat(seq.(float64)).IntPart()))
				blob, err := signing(tx)
				if err != nil {
					callback(err, "")
				} else {
					callback(nil, blob)
				}
			} else {
				callback(errors.New("Request account info fail."), "")
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

	//			transaction.tx_json.SigningPubKey = wt.GetPublicKey().BytesToHex()
	//			var prefix uint32 = 0x53545800
	//			hash, _ := jtSerz.(transaction.tx_json).Hash(prefix)
	//			transaction.tx_json.TxnSignature = wt.signTx(hash)
	//			transaction.tx_json.Blob = jtSerz.FromJson(transaction.tx_json).ToHex()
	//			transaction.local_sign = true
	//			callback(nil, transaction.tx_json.Blob)
	//		})
	//	})
}

/**
 * 提交交易数据
 */
func (tx *Transaction) Submit(callback func(err error, result interface{})) {
	if tx.checkTxError() {
		callback(tx.GetTxJson(constant.TXJSON_ERROR_KEY).(error), nil)
		return
	}

	if tx.remote.LocalSign {
		//本地签名
		tx.sign(func(err error, blob string) {
			if nil != err {
				callback(errors.New("sig error. "+err.Error()), nil)
			} else {
				//var data struct{ tx_blob string }
				//data.tx_blob = transaction.tx_json.Blob
				//transaction.remote.Submit("submit", data, transaction.filter, callback)
				fmt.Printf("blob %s\n", blob)
			}
		})
	} /*else if transaction.tx_json.TransactionType == "Signer" {
		//直接将blob传给底层
		var data struct{ tx_blob string }
		data.tx_blob = transaction.tx_json.Blob
		transaction.remote.Submit("submit", data, transaction.filter, callback)
	} else {
		//不签名交易传给底层
		var data struct {
			tx_json *jtSerz.TxData
			secret  string
		}
		data.tx_json = transaction.tx_json
		data.secret = transaction.secret
		transaction.remote.Submit("submit", data, transaction.filter, callback)
	}*/
}

func (tx *Transaction) checkTxError() bool {
	if tx.GetTxJson(constant.TXJSON_ERROR_KEY) != nil {
		return true
	}
	return false
}
