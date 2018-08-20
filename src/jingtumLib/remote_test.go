/***  读取配置文件
*** remote_test.go
*** 主要用于测试remote
*** author: IPostMan
*** last_modified_time:  2018-08-15 23:13:23
 */

package jingtumLib

import (
	"encoding/json"
	"sync"
	"testing"

	"jingtumLib/constant"
)

//BuildRelationTx请求账号信息
func Test_BuildRelationTx(t *testing.T) {
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust", "quality_out": "100", "quality_in": "1"}
	limit := constant.Amount{}
	limit.Currency = "SWT"
	limit.Value = "100.0001"
	limit.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["limit"] = limit
	req, err := remote.BuildRelationTx(options)
	if err != nil {
		t.Fatalf("BuildRelationTx fail", err.Error())
		wg.Done()
		return
	}
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Build Relation Tx : ", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Build Relation Tx result : ", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}

//BuildAccountSetTx 设置账号属性
func Test_BuildAccountSetTx(t *testing.T) {
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "property", "set": "asfRequireDest", "clear": "asfDisableMaster", "target": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"}
	limit := constant.Amount{}
	limit.Currency = "SWT"
	limit.Value = "100.0001"
	limit.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["limit"] = limit
	req, err := remote.BuildAccountSetTx(options)
	if err != nil {
		t.Fatalf("Build AccountSet Tx fail", err.Error())
		wg.Done()
		return
	}
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Build AccountSet Tx : ", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Build AccountSet Tx result : ", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}

//BuildOfferCreateTx 挂单
func Test_BuildOfferCreateTx(t *testing.T) {
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "sell"}
	gets := constant.Amount{}
	gets.Currency = "SWT"
	pays := constant.Amount{}
	pays.Currency = "CNY"
	options["gets"] = gets
	options["pays"] = pays
	req, err := remote.BuildAccountSetTx(options)
	if err != nil {
		t.Fatalf("BuildOfferCreateTx fail", err.Error())
		wg.Done()
		return
	}
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Build Offer Create Tx : ", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Build Offer Create Tx result : ", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}

//BuildOfferCancelTx 取消挂单
func Test_BuildOfferCancelTx(t *testing.T) {
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
	options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	options["sequence"] = 1
	req, err := remote.BuildAccountSetTx(options)
	if err != nil {
		t.Fatalf("buildOfferCancelTx fail", err.Error())
		wg.Done()
		return
	}
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Build Offer Cancel Tx : ", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Success Build Offer Cancel Tx result : ", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}

