/**
 * 请求测试类
 *
 * @FileName: request_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumLib

import (
	"encoding/json"
	"fmt"
	"jingtumLib/constant"
	"sync"
	"testing"
)

//Test_DeployContractTx 部署合约测试
func Test_DeployContractTx(t *testing.T) {
	remote, err := NewRemote("ws://139.129.194.175:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	//部署合约
	options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "amount": float64(100), "payload": fmt.Sprintf("%X", "result={}; function Init(t) result=scGetAccountBalance(t) return result end; function foo(t) result=scGetAccountBalance(t) return result end"), "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
	tx, err := remote.DeployContractTx(options)
	if err != nil {
		t.Fatalf("Fail request deploy contract %s", err.Error())
		wg.Done()
	} else {
		tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
		tx.Submit(func(err error, data interface{}) {
			if err != nil {
				t.Fatalf("Fail request deploy contract %s", err.Error())
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
		t.Fatalf("Fail request call contract Tx %s", err.Error())
		wg.Done()
	}
	tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
	tx.Submit(func(err error, data interface{}) {
		if err != nil {
			t.Fatalf("Fail request call contract Tx %s", err.Error())
		} else {
			jsonBytes, _ := json.Marshal(data)
			t.Logf("Success call contract Tx : %s", string(jsonBytes))
		}
		wg.Done()
	})
	wg.Wait()
}

//Test_ListenerEvent 监听账本消息
func Test_ListenerEvent(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	//监听所有账本消息
	remote.On(constant.EventLedgerClosed, func(data interface{}) {
		jsonBytes, _ := json.Marshal(data)
		t.Logf("Success listener ledger closed : %s", string(jsonBytes))
		wg.Done()
	})
	wg.Wait()
}

//Test_RequestOrderBook 获得市场挂单列表
func Test_RequestOrderBook(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := make(map[string]interface{})
	gets := constant.Amount{}
	gets.Currency = "SWT"
	pays := constant.Amount{}
	pays.Currency = "CNY"
	pays.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["gets"] = gets
	options["pays"] = pays
	req, err := remote.RequestOrderBook(options)
	if err != nil {
		t.Fatalf("Fail request order book %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request order book %s", err.Error())
			wg.Done()
			return
		}

		// jsonByte, _ := json.Marshal(result)
		t.Logf("Success request order book")
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestAccountTx 获得账号交易列表
func Test_RequestAccountTx(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountTx(options)
	if err != nil {
		t.Fatalf("Fail request account tx %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request account tx %s", err.Error())
			wg.Done()
			return
		}

		// jsonByte, _ := json.Marshal(result)
		t.Log("Success request account tx")
		wg.Done()
	})

	wg.Wait()
}

//RequestAccountOffers 获得账号挂单
func Test_RequestAccountOffers(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountOffers(options)
	if err != nil {
		t.Fatalf("Fail request account offers %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request account offers %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request account offers %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestAccountRelations 获得账号关系
func Test_RequestAccountRelations(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust"}
	req, err := remote.RequestAccountRelations(options)
	if err != nil {
		t.Fatalf("Fail request account relations %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request account relations %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request account relations %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestAccountTums 获得账号可接收和发送的货币
func Test_RequestAccountTums(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountTums(options)
	if err != nil {
		t.Fatalf("Fail request Account Tums %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request Account Tums %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request Account Tums %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestTx 获得某一交易信息
func Test_RequestTx(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	hash := "6537F72CE1DBD8043230C3FF64C6E5E95B11F6573D91EF6A13FEADE6940CB71A"
	req, err := remote.RequestTx(hash)
	if err != nil {
		t.Fatalf("Fail request tx %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request tx %s", err.Error())
			wg.Done()
			return
		}

		// jsonByte, _ := json.Marshal(result)
		t.Log("Success request tx")
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestLedger 获取某一账本
func Test_RequestLedger(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"transactions": true, "ledger_index": 969054, "ledger_hash": "AEE4B16B543D8C8924F09C1DB822C6419780B86019F5F5FF8DC2938E7E0E89D2"}

	req, err := remote.RequestLedger(options)
	if err != nil {
		t.Fatalf("Fail request ledger %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request ledger %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request ledger %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

// Test_RequestLedgerClosed 获取最新账本
func Test_RequestLedgerClosed(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
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

	req, err := remote.RequestLedgerClosed()
	if err != nil {
		t.Fatalf("Fail request ledger closed %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request ledger closed %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request ledger closed %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestServerInfo 获取服务器信息
func Test_RequestServerInfo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
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

	req, err := remote.RequestServerInfo()
	if err != nil {
		t.Fatalf("Fail request server info %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request server info %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request server info %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestAccountInfo 账号信息测试
func Test_RequestAccountInfo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	//请求账号信息
	options := make(map[string]interface{})
	options["account"] = "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	req, err := remote.RequestAccountInfo(options)

	if err != nil {
		t.Fatalf("RequestAccountInfo fail : %s", err.Error())
		wg.Done()
		return
	}

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Requst account info : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Requst account info result : %s", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}
