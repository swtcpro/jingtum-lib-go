/***  初始化
 *** testLib.go
 *** 主要用于用于测试jingtumLib的各个实例
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package main

import (
	"fmt"
	"os"
	"sync"

	jingtum "jingtumLib"
)

func main() {
	err := jingtum.Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	remote, err := jingtum.NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		fmt.Printf("New remote fail : %s", err)
		return
	}

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}

		fmt.Println(result)
	})

	if cerr != nil {
		fmt.Printf("Connect service fail : %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	//请求账号信息
	// options := make(map[string]interface{})
	// options["account"] = "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	// req, err := remote.RequestAccountInfo(options)

	// if err != nil {
	// 	fmt.Printf("RequestAccountInfo fail : %v", err)
	// 	return
	// }

	// req.Submit(func(err error, result interface{}) {
	// 	if err != nil {
	// 		fmt.Printf("Requst account info : %v\n", err)
	// 		wg.Done()
	// 		return
	// 	}

	// 	fmt.Printf("Requst submit result : %v", result)
	// 	wg.Done()
	// })

	//支付请求
	// var v struct {
	// 	account string
	// 	secret  string
	// }
	// // v.account = "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"
	// // v.secret = "saNUs41BdTWSwBRqSTbkNdjnAVR8h"
	// // to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk" //"jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	// v.account = "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	// v.secret = "ssc5eiFivvU2otV6bSYmJeZrAsQK3"
	// to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	// amount := constant.Amount{}
	// amount.Currency = "SWT"
	// amount.Value = "0.0001"
	// tx, err := remote.BuildPaymentTx(v.account, to, amount)
	// if err != nil {
	// 	fmt.Printf("Build paymanet tx fail : %s\n", err)
	// 	return
	// }
	// tx.SetSecret(v.secret)
	// tx.AddMemo("支付0.0001SWT")
	// tx.Submit(func(err error, result interface{}) {
	// 	if err != nil {
	// 		fmt.Printf("Payment fail : %v\n", err)
	// 		wg.Done()
	// 		return
	// 	}

	// 	jsonByte, _ := json.Marshal(result)

	// 	fmt.Printf("Payment result : %s\n", jsonByte)
	// 	wg.Done()
	// })

	//请求服务信息
	// req, err := remote.RequestServerInfo()
	// if err != nil {
	// 	fmt.Printf("Fail request server info %s", err.Error())
	// }

	// req.Submit(func(err error, result interface{}) {
	// 	if err != nil {
	// 		fmt.Printf("Fail request server info %s", err.Error())
	// 		wg.Done()
	// 		return
	// 	}

	// 	jsonByte, _ := json.Marshal(result)
	// 	fmt.Printf("Success request server info %s", jsonByte)
	// 	wg.Done()
	// })

	//请求市场挂单
	// options := make(map[string]interface{}) //{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	// gets := constant.Amount{}
	// gets.Currency = "SWT"
	// pays := constant.Amount{}
	// pays.Currency = "CNY"
	// pays.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	// options["gets"] = gets
	// options["pays"] = pays
	// req, err := remote.RequestOrderBook(options)
	// if err != nil {
	// 	fmt.Printf("Fail request order book %s", err.Error())
	// 	return
	// }

	// if err != nil {
	// 	fmt.Printf("Fail request order book %s", err.Error())
	// 	return
	// }

	// req.Submit(func(err error, result interface{}) {
	// 	if err != nil {
	// 		fmt.Printf("Fail request order book %s", err.Error())
	// 		wg.Done()
	// 		return
	// 	}

	// 	jsonByte, _ := json.Marshal(result)
	// 	fmt.Printf("Success request order book %s", jsonByte)
	// 	wg.Done()
	// })
	wg.Add(3)
	// remote.On(constant.EventLedgerClosed, func(data interface{}) {
	// 	jsonBytes, _ := json.Marshal(data)
	// 	fmt.Printf("Success listener ledger closed : %s", string(jsonBytes))
	// 	wg.Done()
	// })
	wg.Wait()

	//	defer jingtum.Exits()
}
