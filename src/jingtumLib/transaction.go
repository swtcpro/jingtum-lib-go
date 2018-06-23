/**
 * Transaction类主管POST请求，包括组装交易和交易参数。请求时需要提供密钥，且交易可以
 * 进行本地签名和服务器签名。目前支持服务器签名，本地签名支持主要的交易，还有部分参数
 * 不支持。所有的请求是异步的，会提供一个回调函数。每个回调函数有两个参数，一个是错误，
 * 另一个是结果。
 *
 * @FileName: transaction.go
 * @Auther : 13851485286
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
    remote        Remote
    filter        Filter
    secret        string
    tx_json       *TxData
}

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