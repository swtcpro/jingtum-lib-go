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

//Test_BuildOfferCancelTx 取消挂单
func Test_BuildOfferCancelTx(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err)
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		t.Logf("Connect to %s success. Result : %s.", "ws://123.57.219.57:5020", jsonByte)
	})

	if conErr != nil {
		t.Fatalf("Connect to %s fail : %s", "ws://123.57.219.57:5020", conErr.Error())
		return
	}

	defer remote.Disconnect()

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "sequence": uint32(26)}
	tx, err := remote.BuildOfferCancelTx(options)
	if err != nil {
		t.Fatalf("Fail BuildOfferCancelTx : %s", err.Error())
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	tx.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail BuildOfferCancelTx : %s", err.Error())
			wg.Done()
		} else {
			jsonBytes, _ := json.Marshal(result)
			t.Logf("Success BuildOfferCancelTx : %s", jsonBytes)
			wg.Done()
		}
	})

	wg.Wait()
}

//Test_DeployContractTx 部署合约测试
func Test_DeployContractTx(t *testing.T) {
	remote, err := NewRemote("ws://139.129.194.175:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if conErr != nil {
		t.Fatalf("Connect service fail : %s", conErr.Error())
		return
	}

	defer remote.Disconnect()

	wg := sync.WaitGroup{}
	wg.Add(1)
	//部署合约
	options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "amount": float64(100), "payload": fmt.Sprintf("%X", "result={}; function Init(t) result=scGetAccountBalance(t) return result end; function foo(t) result=scGetAccountBalance(t) return result end"), "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
	tx, err := remote.DeployContractTx(options)
	if err != nil {
		t.Errorf("Fail request deploy contract %s", err.Error())
		wg.Done()
	} else {
		tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
		tx.Submit(func(err error, data interface{}) {
			if err != nil {
				t.Errorf("Fail request deploy contract %s", err.Error())
			} else {
				jsonBytes, _ := json.Marshal(data)
				t.Logf("Success deploy contract : %s", string(jsonBytes))
			}
			wg.Done()
		})
	}
	wg.Wait()
}

//Test_CallContractTx 执行合约
func Test_CallContractTx(t *testing.T) {
	//执行合约
	remote, err := NewRemote("ws://139.129.194.175:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if conErr != nil {
		t.Fatalf("Connect service fail : %s", conErr.Error())
		return
	}

	defer remote.Disconnect()

	wg := sync.WaitGroup{}
	wg.Add(1)
	options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "destination": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72", "foo": "foo", "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
	tx, err := remote.CallContractTx(options)
	if err != nil {
		t.Errorf("Fail request call contract Tx %s", err.Error())
		wg.Done()
	}
	tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
	tx.Submit(func(err error, data interface{}) {
		if err != nil {
			t.Errorf("Fail request call contract Tx %s", err.Error())
		} else {
			jsonBytes, _ := json.Marshal(data)
			t.Logf("Success call contract Tx : %s", string(jsonBytes))
		}
		wg.Done()
	})
	wg.Wait()
}

//Test_AddMemo 备注测试
func Test_AddMemo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", false)
	if err != nil {
		t.Fatalf("New remote error : %v.", err)
	}
	defer remote.Disconnect()
	remote.LocalSign = true

	t.Logf("Success remote : %v.", remote)
	tx, err := NewTransaction(remote, nil)

	if err != nil {
		t.Fatalf("New transactino error : %s.", err.Error())
	}

	tx.AddTxJSON("TransactionType", "Payment")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("我的测试")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("支付0.000001SWT")
	tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	t.Logf("Get tx TransactionType : %s. Flags : %d. Fee : %d. Secret : %s", tx.GetTxJSON("TransactionType"), tx.GetTxJSON("Flags"), tx.GetTxJSON("Fee"), tx.GetTxJSON("secret"))

	memos := tx.GetTxJSON("Memos").(*list.List)

	for e := memos.Front(); e != nil; e = e.Next() {
		t.Logf("Get tx memos info %s.", e.Value.(*serializer.MemoInfo).Memo.MemoData)
	}
}

func Test_LocalSignPayment(t *testing.T) {
	wsurl := "ws://123.57.219.57:5020"
	remote, err := NewRemote(wsurl, true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err)
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		t.Logf("Connect to %s success. Result : %s.", wsurl, jsonByte)
	})

	if conErr != nil {
		t.Fatalf("Connect to %s fail : %s", wsurl, conErr.Error())
		return
	}

	defer remote.Disconnect()

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
	wg := sync.WaitGroup{}
	wg.Add(1)
	tx.SetSecret(v.secret)
	tx.AddMemo("支付0.0001SWT")
	tx.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail Payment : %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)

		t.Logf("Success Payment result : %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

/*
*以下为remote 性能测试用例
 */

//BenchmarkBuildOfferCancelTx 取消挂单
func BenchmarkBuildOfferCancelTx(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err)
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		b.Logf("Connect to %s success. Result : %s.", "ws://123.57.219.57:5020", jsonByte)
	})

	if conErr != nil {
		b.Fatalf("Connect to %s fail : %s", "ws://123.57.219.57:5020", conErr.Error())
		return
	}

	defer remote.Disconnect()

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "sequence": uint32(26)}
	for i := 0; i < b.N; i++ {
		tx, err := remote.BuildOfferCancelTx(options)
		if err != nil {
			b.Fatalf("Fail BuildOfferCancelTx : %s", err.Error())
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
		tx.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail BuildOfferCancelTx : %s", err.Error())
				wg.Done()
			} else {
				jsonBytes, _ := json.Marshal(result)
				b.Logf("Success BuildOfferCancelTx : %s", jsonBytes)
				wg.Done()
			}
		})

		wg.Wait()
	}
}

//BenchmarkDeployContractTx 部署合约测试
func BenchmarkDeployContractTx(b *testing.B) {
	remote, err := NewRemote("ws://139.129.194.175:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		b.Logf("Connect success : %s", jsonBytes)
	})

	if conErr != nil {
		b.Fatalf("Connect service fail : %s", conErr.Error())
		return
	}

	defer remote.Disconnect()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(1)
		//部署合约
		options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "amount": float64(100), "payload": fmt.Sprintf("%X", "result={}; function Init(t) result=scGetAccountBalance(t) return result end; function foo(t) result=scGetAccountBalance(t) return result end"), "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
		tx, err := remote.DeployContractTx(options)
		if err != nil {
			b.Errorf("Fail request deploy contract %s", err.Error())
			wg.Done()
		} else {
			tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
			tx.Submit(func(err error, data interface{}) {
				if err != nil {
					b.Errorf("Fail request deploy contract %s", err.Error())
				} else {
					jsonBytes, _ := json.Marshal(data)
					b.Logf("Success deploy contract : %s", string(jsonBytes))
				}
				wg.Done()
			})
		}
		wg.Wait()
	}
}

//BenchmarkCallContractTx 执行合约
func BenchmarkCallContractTx(b *testing.B) {
	//执行合约
	remote, err := NewRemote("ws://139.129.194.175:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		b.Logf("Connect success : %s", jsonBytes)
	})

	if conErr != nil {
		b.Fatalf("Connect service fail : %s", conErr.Error())
		return
	}

	defer remote.Disconnect()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(1)
		options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "destination": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72", "foo": "foo", "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
		tx, err := remote.CallContractTx(options)
		if err != nil {
			b.Errorf("Fail request call contract Tx %s", err.Error())
			wg.Done()
		}
		tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
		tx.Submit(func(err error, data interface{}) {
			if err != nil {
				b.Errorf("Fail request call contract Tx %s", err.Error())
			} else {
				jsonBytes, _ := json.Marshal(data)
				b.Logf("Success call contract Tx : %s", string(jsonBytes))
			}
			wg.Done()
		})
		wg.Wait()
	}
}

//BenchmarkAddMemo 备注测试
func BenchmarkAddMemo(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", false)
	if err != nil {
		b.Fatalf("New remote error : %v.", err)
	}
	defer remote.Disconnect()
	remote.LocalSign = true

	b.Logf("Success remote : %v.", remote)
	for i := 0; i < b.N; i++ {
		tx, err := NewTransaction(remote, nil)

		if err != nil {
			b.Fatalf("New transactino error : %s.", err.Error())
		}

		tx.AddTxJSON("TransactionType", "Payment")
		tx.AddMemo("支付0.000001SWT")
		tx.AddMemo("我的测试")
		tx.AddMemo("支付0.000001SWT")
		tx.AddMemo("支付0.000001SWT")
		tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
		b.Logf("Get tx TransactionType : %s. Flags : %d. Fee : %d. Secret : %s", tx.GetTxJSON("TransactionType"), tx.GetTxJSON("Flags"), tx.GetTxJSON("Fee"), tx.GetTxJSON("secret"))

		memos := tx.GetTxJSON("Memos").(*list.List)

		for e := memos.Front(); e != nil; e = e.Next() {
			b.Logf("Get tx memos info %s.", e.Value.(*serializer.MemoInfo).Memo.MemoData)
		}
	}
}

func BenchmarkLocalSignPayment(b *testing.B) {
	wsurl := "ws://123.57.219.57:5020"
	remote, err := NewRemote(wsurl, true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err)
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		b.Logf("Connect to %s success. Result : %s.", wsurl, jsonByte)
	})

	if conErr != nil {
		b.Fatalf("Connect to %s fail : %s", wsurl, conErr.Error())
		return
	}

	defer remote.Disconnect()

	for i := 0; i < b.N; i++ {
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
			b.Fatalf("Build paymanet tx fail : %s", err.Error())
			continue
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		tx.SetSecret(v.secret)
		tx.AddMemo("支付0.0001SWT")
		tx.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail Payment : %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)

			b.Logf("Success Payment result : %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
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
