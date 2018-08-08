/**
 * 交易测试类
 *
 * @FileName: transaction_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumLib

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"jingtumLib/constant"
	"jingtumLib/serializer"
)

//Test_AddMemo 备注测试
func Test_AddMemo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", false)
	if err != nil {
		t.Fatalf("New remote error : %v.", err)
	}
	defer remote.Disconnect()
	remote.LocalSign = true

	t.Logf("Success remote : %v.", remote)
	tx, err := NewTransaction(remote)

	if err != nil {
		t.Fatalf("New transactino error : %s.", err.Error())
	}

	tx.AddTxJson("TransactionType", "Payment")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("我的测试")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("支付0.000001SWT")
	tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	t.Logf("Get tx TransactionType : %s. Flags : %d. Fee : %d. Secret : %s", tx.GetTxJson("TransactionType"), tx.GetTxJson("Flags"), tx.GetTxJson("Fee"), tx.GetTxJson("secret"))

	memos := tx.GetTxJson("Memos").(*list.List)

	for e := memos.Front(); e != nil; e = e.Next() {
		t.Logf("Get tx memos info %s.", e.Value.(*serializer.MemoInfo).Memo.MemoData)
	}
}

func Test_LocalSignPayment(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wsurl := "ws://123.57.219.57:5020"
	remote, err := NewRemote(wsurl, true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err)
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		t.Logf("Connect to %s success. Result : %s.", wsurl, jsonByte)
	})

	if cerr != nil {
		t.Fatalf("Connect to %s fail : %s", wsurl, err.Error())
		return
	}

	//支付请求
	var v struct {
		account string
		secret  string
	}
	v.account = "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	v.secret = "ssc5eiFivvU2otV6bSYmJeZrAsQK3"
	to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk" //"jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	amount := constant.Amount{}
	amount.Currency = "SWT"
	amount.Value = "0.0001"
	tx, err := remote.BuildPaymentTx(v.account, to, amount)
	if err != nil {
		t.Fatalf("Build paymanet tx fail : %s", err.Error())
		return
	}
	tx.SetSecret(v.secret)
	tx.AddMemo("支付0.0001SWT")
	tx.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail Payment : %s", err.Error())
			return
		}

		jsonByte, _ := json.Marshal(result)

		t.Logf("Success Payment result : %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

func TestMain(m *testing.M) {
	err := Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	ret := m.Run()
	os.Exit(ret)
}
