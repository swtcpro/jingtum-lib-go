//performance_test.go 性能测试
package jingtumlib

import (
	"encoding/json"
	//"fmt"
	//"os"
	"sync"
	"testing"
	"time"
)

var remote *Remote

/*
func TestMain(m *testing.M) {
	err := Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	remote, err = NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		return
	}

	conErr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}
	})

	if conErr != nil {
		return
	}
	ret := m.Run()
	os.Exit(ret)
}
*/
func BenchmarkBuildPaymentTx(B *testing.B) {
	B.Logf("Remote is null %t", remote == nil)
	wg := sync.WaitGroup{}
	wg.Add(B.N)
	for i := 0; i < B.N; i++ {
		//支付请求
		var v struct {
			account string
			secret  string
		}
		v.account = "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
		v.secret = "ssc5eiFivvU2otV6bSYmJeZrAsQK3"
		to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk" //"jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
		amount := Amount{}
		amount.Currency = "SWT"
		amount.Value = "0.0001"
		start := time.Now().Unix()
		B.Logf("current sec %d", start)
		tx, err := remote.BuildPaymentTx(v.account, to, amount)
		B.Logf("Builder payment cost %d", time.Now().Unix()-start)
		if err != nil {
			B.Fatalf("Build paymanet tx fail : %s", err.Error())
			continue
		}
		tx.SetSecret(v.secret)
		tx.AddMemo("支付0.0001SWT")
		start = time.Now().Unix()
		tx.Submit(func(err error, result interface{}) {
			if err != nil {
				B.Errorf("Fail Payment : %s", err.Error())
				wg.Done()
				return
			}

			jsonByte, _ := json.Marshal(result)
			B.Logf("Submit payment cost %d", time.Now().Unix()-start)
			B.Logf("Success Payment result : %s", jsonByte)
			wg.Done()
		})
	}
	wg.Wait()
}
