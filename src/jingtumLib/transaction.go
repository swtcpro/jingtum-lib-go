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
     "fmt"
     "encoding/json"
     "strings"
     "errors"
)

type FlagClass map[string] uint32

type Amount struct {
    Currency        string
    Issuer          string
    Value           string
}

type MemoInfo struct {
    MemoData        string
}

type TxData struct {
    Flags           uint32
    Fee             interface{}
    Account         string
    TransactionType interface{}
    SendMax         interface{}
    Memos           []MemoInfo
    Paths           interface{}
    SendMax         interface{}
    TransferRate    uint32
    MemoType        interface{}
    MemoLen         interface{}
}

type Transaction struct {
    remote        *Remote
    filter        Filter
    secret        string
    tx_json       *TxData
}

var (
     TransactionFlags = map[string]FlagClass {"Universal":{"FullyCanonicalSig":0x00010000},"AccountSet":{"RequireDestTag":0x00010000,"OptionalDestTag":0x00020000,"RequireAuth":0x00040000,"OptionalAuth":0x00080000,"DisallowSWT":0x00100000,"AllowSWT":0x00200000},"TrustSet":{"SetAuth":0x00010000,"NoSkywell":0x00020000,"SetNoSkywell":0x00020000,"ClearNoSkywell":0x00040000,"SetFreeze":0x00100000,"ClearFreeze":0x00200000},"OfferCreate":{"Passive":0x00010000,"ImmediateOrCancel":0x00020000,"FillOrKill":0x00040000,"Sell":0x00080000},"Payment":{"NoSkywellDirect":0x00010000,"PartialPayment":0x00020000,"LimitQuality":0x00040000},"RelationSet":{"Authorize":0x00000001,"Freeze":0x00000011}}
)

func NewTransaction(remote *Remote, filter Filter) (transaction *Transaction , err error) {
    transaction = new(Transaction)
    transaction.remote = remote
    transaction.tx_json = new(TxData)
    transaction.tx_json.Flags = 0
    transaction.tx_json.Fee = JTConfig.ReadInt("Config","fee", 10000);
    transaction.filter = filter
}

func (transaction *Transaction) ParseJson (jsonStr string) (err error) {

    err:=json.Unmarshal([]byte(jsonStr),&transaction.tx_json)
}

func (transaction *Transaction) GetAccount() (stirng) {
    return transaction.tx_json.Account
}

func (transaction *Transaction) GetTransactionType() (interface{}) {
    return transaction.tx_json.TransactionType
}

func (transaction *Transaction) SetSecret (secret string) {
    transaction.secret = secret
}

func (transaction *Transaction AddMemo(memo string) {
    if (len(memo) > 2048) {
       transaction.tx_json.MemoLen = errors.New("The length of Memo shoule be less than or equal 2048.")
       return
    }

    var _memo MemoInfo = new(MemoInfo)
    _memo.MemoData = stringToHex(memo)

    append(transaction.tx_json.Memos, _memo)
}

func (transaction *Transaction) SetFee(fee uint32) {
    if (fee < 10) {
        transaction.tx_json.Fee = errors.New("Fee should be great than or equal 10.")
        return
    }

    transaction.tx_json.Fee = fee
}

func (transaction *Transaction) MaxAmount(amount interface{}) (interface{}) {
    if _, ok := amount.(string); ok {
        if (number(amount)) {
            f, err := strconv.ParseFloat(amount, 32)
            if (err != nil) {
                return errors.New("invalid amount to max")
            }

            return strconv.FormatFloat(f * 1.0001, 'f', 10, 32)
        }
    }

    if _, ok := amount.(Amount); ok && isValidAmount(amount) {
        f, err := strconv.ParseFloat(amount.value, 32)
           if (err != nil) {
               return errors.New("invalid amount to max")
           }
        return strconv.FormatFloat(f * 1.0001, 'f', 10, 32)
    }

    return errors.New("invalid amount to max")
}
/**
 * set a path to payment
 * this path is repesented as a key, which is computed in path find
 * so if one path you computed self is not allowed
 * when path set, sendmax is also set.
 * @param path
 */
func (transaction * Transaction) setPath(key string) (err error) {
    if key == "" ||  (strings.Count(key,"")-1) != 40 {
        err = errors.New("invalid path key")
        return
    }

    item := transaction.remote.paths.get(key)

    if item == nil {
        err = errors.New("non exists path key")
        return
    }

    //沒有支付路径，不需要传下面的参数
    if(item.path == nil || len(item.path) == 0) {
        return;
    }
    var path []interface{}
    err := json.Unmarshal(path,&path) 
    if err != nil { 
        err = errors.New("invalid path.")
        return
    } 

    transaction.tx_json.Paths = path;
    var amount = MaxAmount(item.choice);
    transaction.tx_json.SendMax = amount;
}

/**
 * limit send max amount
*/
func (transaction *Transaction) setSendMax(amount Amount) {
    if !isValidAmount(amount) {
       transaction.tx_json.SendMax = errors.New("invalid send max amount")
       return
    }

    transaction.tx_json.SendMax = amount
}

/**
 * transfer rate
 * between 0 and 1, type is number
 * @param rate
 */
func (transaction *Transaction) setTransferRate(rate float32) (err error) {
    if (rate < 0 || rate > 1) {
        err = errors.New("invalid transfer rate")
        return
    }

    transaction.tx_json.TransferRate = (rate + 1) * 1000000000;
}

/**
 * set transaction flags
 *
 */
func (transaction *Transaction) setFlags(flags interface{}) (err error) {
    if _, ok := flags.(uint32); ok {
        transaction.tx_json.Flags = flags
        return
    }

    if _, arrayOk := flags.([]string); arrayOk {
       
    }

    err = errors.New("invalid flags")

    var transaction_flags = Transaction.flags[this.getTransactionType()] || {};
    var flag_set = Array.isArray(flags) ? flags : [].concat(flags);
    for (var i = 0; i < flag_set.length; ++i) {
        var flag = flag_set[i];
        if (transaction_flags.hasOwnProperty(flag)) {
            this.tx_json.Flags += transaction_flags[flag];
        }
    }
}