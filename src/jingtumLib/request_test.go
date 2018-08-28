/**
 * 请求测试类
 *
 * @FileName: request_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumlib

import (
	"encoding/json"
	"jingtumlib/constant"
	"sync"
	"testing"
)

//Test_ListenerEvent 监听账本消息
func Test_ListenerEvent(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Errorf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	options := make(map[string]interface{})
	gets := Amount{}
	gets.Currency = "SWT"
	pays := Amount{}
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
			t.Errorf("Fail request order book %s", err.Error())
			wg.Done()
			return
		}
		t.Log("Success request order book")
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
	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountTx(options)
	if err != nil {
		t.Fatalf("Fail request account tx %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request account tx %s", err.Error())
			wg.Done()
			return
		}
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountOffers(options)
	if err != nil {
		t.Fatalf("Fail request account offers %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request account offers %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust"}
	req, err := remote.RequestAccountRelations(options)
	if err != nil {
		t.Fatalf("Fail request account relations %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request account relations %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountTums(options)
	if err != nil {
		t.Fatalf("Fail request Account Tums %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request Account Tums %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	hash := "6537F72CE1DBD8043230C3FF64C6E5E95B11F6573D91EF6A13FEADE6940CB71A"
	req, err := remote.RequestTx(hash)
	if err != nil {
		t.Fatalf("Fail request tx %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request tx %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	options := map[string]interface{}{"transactions": true, "ledger_index": 969054, "ledger_hash": "AEE4B16B543D8C8924F09C1DB822C6419780B86019F5F5FF8DC2938E7E0E89D2"}

	req, err := remote.RequestLedger(options)
	if err != nil {
		t.Fatalf("Fail request ledger %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Fail request ledger %s", err.Error())
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
			t.Errorf("New remote fail : %s", err.Error())
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
			t.Errorf("Fail request ledger closed %s", err.Error())
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
			t.Errorf("New remote fail : %s", err.Error())
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
			t.Errorf("Fail request server info %s", err.Error())
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

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("New remote fail : %s", err.Error())
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

	//请求账号信息
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	req, err := remote.RequestAccountInfo(options)

	if err != nil {
		t.Fatalf("RequestAccountInfo fail : %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.SelectLedger(1065000)
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Errorf("Requst account info : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Requst account info result : %s", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}

/*
*以下为request性能测试用例
 */
/*
//BenchmarkListenerEvent 监听账本消息
func BenchmarkListenerEvent(b *testing.B) {
	for i := 0; i < b.N; i++ {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Errorf("New remote fail : %s", err.Error())
		continue
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
			continue
		}

		jsonBytes, _ := json.Marshal(result)

		b.Logf("Connect success : %s", jsonBytes)
	})

	if conErr != nil {
		b.Fatalf("Connect service fail : %s", conErr.Error())
		continue
	}

	defer remote.Disconnect()
	wg := sync.WaitGroup{}
	wg.Add(2)
	//监听所有账本消息
	remote.On(constanb.EventLedgerClosed, func(data interface{}) {
		jsonBytes, _ := json.Marshal(data)
		b.Logf("Success listener ledger closed : %s", string(jsonBytes))
		wg.Done()
	})
	wg.Wait()
	}
}*/

//BenchmarkRequestOrderBook 获得市场挂单列表
func BenchmarkRequestOrderBook(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := make(map[string]interface{})
		gets := Amount{}
		gets.Currency = "SWT"
		pays := Amount{}
		pays.Currency = "CNY"
		pays.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
		options["gets"] = gets
		options["pays"] = pays
		req, err := remote.RequestOrderBook(options)
		if err != nil {
			b.Fatalf("Fail request order book %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request order book %s", err.Error())
				wg.Done()
				return
			}
			b.Log("Success request order book")
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestAccountTx 获得账号交易列表
func BenchmarkRequestAccountTx(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}
	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
		req, err := remote.RequestAccountTx(options)
		if err != nil {
			b.Fatalf("Fail request account tx %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request account tx %s", err.Error())
				wg.Done()
				return
			}
			b.Log("Success request account tx")
			wg.Done()
		})

		wg.Wait()
	}
}

//RequestAccountOffers 获得账号挂单
func BenchmarkRequestAccountOffers(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
		req, err := remote.RequestAccountOffers(options)
		if err != nil {
			b.Fatalf("Fail request account offers %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request account offers %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request account offers %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestAccountRelations 获得账号关系
func BenchmarkRequestAccountRelations(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust"}
		req, err := remote.RequestAccountRelations(options)
		if err != nil {
			b.Fatalf("Fail request account relations %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request account relations %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request account relations %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestAccountTums 获得账号可接收和发送的货币
func BenchmarkRequestAccountTums(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
		req, err := remote.RequestAccountTums(options)
		if err != nil {
			b.Fatalf("Fail request Account Tums %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request Account Tums %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request Account Tums %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestTx 获得某一交易信息
func BenchmarkRequestTx(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		hash := "6537F72CE1DBD8043230C3FF64C6E5E95B11F6573D91EF6A13FEADE6940CB71A"
		req, err := remote.RequestTx(hash)
		if err != nil {
			b.Fatalf("Fail request tx %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request tx %s", err.Error())
				wg.Done()
				return
			}

			// jsonByte, _ := json.Marshal(result)
			b.Log("Success request tx")
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestLedger 获取某一账本
func BenchmarkRequestLedger(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		options := map[string]interface{}{"transactions": true, "ledger_index": 969054, "ledger_hash": "AEE4B16B543D8C8924F09C1DB822C6419780B86019F5F5FF8DC2938E7E0E89D2"}

		req, err := remote.RequestLedger(options)
		if err != nil {
			b.Fatalf("Fail request ledger %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request ledger %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request ledger %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

// BenchmarkRequestLedgerClosed 获取最新账本
func BenchmarkRequestLedgerClosed(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		req, err := remote.RequestLedgerClosed()
		if err != nil {
			b.Fatalf("Fail request ledger closed %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request ledger closed %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request ledger closed %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestServerInfo 获取服务器信息
func BenchmarkRequestServerInfo(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		req, err := remote.RequestServerInfo()
		if err != nil {
			b.Fatalf("Fail request server info %s", err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Fail request server info %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			b.Logf("Success request server info %s", jsonByte)
			wg.Done()
		})

		wg.Wait()
	}
}

//BenchmarkRequestAccountInfo 账号信息测试
func BenchmarkRequestAccountInfo(b *testing.B) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		b.Fatalf("New remote fail : %s", err.Error())
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			b.Errorf("New remote fail : %s", err.Error())
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
		//请求账号信息
		options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
		req, err := remote.RequestAccountInfo(options)

		if err != nil {
			b.Fatalf("RequestAccountInfo fail : %s", err.Error())
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		req.SelectLedger(1065000)
		req.Submit(func(err error, result interface{}) {
			if err != nil {
				b.Errorf("Requst account info : %s", err.Error())
				wg.Done()
				return
			}
			jsonBytes, _ := json.Marshal(result)
			b.Logf("Success Requst account info result : %s", jsonBytes)
			wg.Done()
		})

		wg.Wait()
	}
}
