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
	//	"encoding/json"
	//	"errors"
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
	tx.AddTxJson("Flags", 0)
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
	return tx.tx_json[key]
}

/**
 * 本地签名时需要设置私钥
 */
func (tx *Transaction) SetSecret(secret string) {
	if !IsValidSecret(secret) {
		tx.AddTxJson(constant.TXJSON_ERROR_KEY, constant.ERR_PAYMENT_INVALID_SECRET)
		return
	}

	tx.AddTxJson("secret", secret)
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

	//	memos := tx.tx_json["Memos"].([]serializer.MemoInfo)
	//
	//	if memos == nil {
	//		memos = make([]serializer.MemoInfo, 0) //[]MemoInfo
	//		tx.AddTxJson("Memos", memos)
	//		memos = append(memos, *mi)
	//	}
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

func (transaction *Transaction) Sign(callback func(param ...interface{})) {
	/*err, remote := NewRemote()
	if err != nil {
		panic(err.Error())
	}

	defer remote.Disconnect()

	remote.Connect(func(err error, result interface{}) {
		if err != nil {
			callback(err)
			return
		}

		errReq, req := transaction.remote.RequestAccountInfo(map[string]string{"account": transaction.tx_json.Account, "type": "trust"})

		if errReq != nil {
			panic(errReq.Error())
		}

		req.Submit(func(err error, data interface{}) {
			if nil != err {
				callback(err)
				return
			}
			//transaction.tx_json.Sequence = data.account_data.Sequence
			transaction.tx_json.Fee = transaction.tx_json.Fee.(float64) / 1000000

			//payment
			if nil != transaction.tx_json.TxAmount {
				if amount, ok := transaction.tx_json.TxAmount.(string); ok {
					//基础货币
					if number(amount) {
						f, err := strconv.ParseFloat(amount, 32)
						if err == nil {
							transaction.tx_json.TxAmount = f / 1000000
						}
					}
				}
			}

			if len(transaction.tx_json.Memos) > 0 {
				mdStr, errMd := hexToString(transaction.tx_json.Memos[0].Memo.MemoData)
				transaction.tx_json.Memos[0].Memo.MemoData = mdStr
			}
			if nil != transaction.tx_json.SendMax {
				if sendMax, ok := transaction.tx_json.SendMax.(string); ok {
					if number(sendMax) {
						f, err := strconv.ParseFloat(sendMax, 32)
						if err == nil {
							transaction.tx_json.SendMax = f / 1000000
						}
					}
				}
			}

			//order
			if nil != transaction.tx_json.TakerPays {
				//基础货币
				if takerPays, ok := transaction.tx_json.TakerPays.(string); ok {
					if number(takerPays) {
						f, err := strconv.ParseFloat(takerPays, 32)
						if err == nil {
							transaction.tx_json.TakerPays = f / 1000000
						}
					}
				}
			}

			if nil != transaction.tx_json.TakerGets {
				//基础货币
				if takerGets, ok := transaction.tx_json.TakerGets.(string); ok {
					if number(takerGets) {
						f, err := strconv.ParseFloat(takerGets, 32)
						if err == nil {
							transaction.tx_json.TakerGets = f / 1000000
						}
					}
				}
			}

			wt, err := jtbLib.FromSecret(transaction.secret)
			if err != nil {
				panic(err.Error())
			}

			transaction.tx_json.SigningPubKey = wt.GetPublicKey().BytesToHex()
			var prefix uint32 = 0x53545800
			hash, _ := jtSerz.(transaction.tx_json).Hash(prefix)
			transaction.tx_json.TxnSignature = wt.signTx(hash)
			transaction.tx_json.Blob = jtSerz.FromJson(transaction.tx_json).ToHex()
			transaction.local_sign = true
			callback(nil, transaction.tx_json.Blob)
		})
	})*/
}

func (transaction *Transaction) Submit(callback func(param ...interface{})) {
	/*err := checkTxError(transaction.tx_json)
	if err != nil {
		callback(err)
		return
	}

	if transaction.remote.LocalSign {
		transaction.Sign(func(err error, blob string) {
			if nil != err {
				callback(errors.New("sig error. " + err.Error()))
			} else {
				var data struct{ tx_blob string }
				data.tx_blob = transaction.tx_json.Blob
				transaction.remote.Submit("submit", data, transaction.filter, callback)
			}
		})
	} else if transaction.tx_json.TransactionType == "Signer" {
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

//func checkTxError(txJson *jtSerz.TxData) error {
//	fields := utils.GetFieldNames(txJson)
//	for _, fn := range fields {
//		v := utils.GetFieldValue(txJson, fn)
//		if err, ok := v.(error); ok {
//			return err
//		}
//	}
//	return nil
//}
