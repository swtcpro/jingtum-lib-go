/***  初始化
 *** testLib.go
 *** 主要用于用于测试jingtumLib的各个实例
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	jingtum "jingtumLib"
	"jingtumLib/constant"
)

func main() {
	err := jingtum.Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	// fields := []string{"Flags", "Fee", "TransactionType", "Account", "Amount", "Destination", "Memos", "Sequence", "SigningPubKey"}
	// fmt.Printf("%v", fields)
	// utils.SortByFieldName(fields)
	// fmt.Printf("%v", fields)
	// return

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
	//	options := make(map[string]interface{})
	//	options["account"] = "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	//	req, err := remote.RequestAccountInfo(options)
	//
	//	if err != nil {
	//		fmt.Printf("RequestAccountInfo fail : %v", err)
	//		return
	//	}
	//
	//	req.Submit(func(err error, result interface{}) {
	//		if err != nil {
	//			fmt.Printf("Requst account info : %v\n", err)
	//			wg.Done()
	//			return
	//		}
	//
	//		fmt.Printf("Requst submit result : %v", result)
	//		wg.Done()
	//	})

	//支付请求
	var v struct {
		account string
		secret  string
	}
	v.account = "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"
	v.secret = "saNUs41BdTWSwBRqSTbkNdjnAVR8h"
	to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk" //"jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	amount := constant.Amount{}
	amount.Currency = "SWT"
	amount.Value = "10"
	tx, err := remote.BuildPaymentTx(v.account, to, amount)
	if err != nil {
		fmt.Printf("Build paymanet tx fail : %s\n", err)
		return
	}
	tx.SetSecret(v.secret)
	tx.AddMemo("支付10SWT")
	tx.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Payment fail : %v\n", err)
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)

		fmt.Printf("Payment result : %s\n", jsonByte)
		wg.Done()
	})

	wg.Wait()

	/*
		isNumber := jingtum.Number("5445")
		fmt.Println(isNumber)

		jingtum.Generate()
	*/
	//	_, remote := jingtum.NewRemote()
	//	err = remote.Connect()
	//	if err != nil {
	//		fmt.Println("Connect service", remote.Wsconn.Host, remote.Wsconn.Port, "fail.", err)
	//		return
	//	}
	//	fmt.Println("Connect service", remote.Wsconn.Host, remote.Wsconn.Port, "succ.")
	//
	//	//请求底层服务器信息
	//	err, response := remote.RequestServerInfo()
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Server Info succ.len=", len(response), "data=", response)
	//
	//	//获取最新账本信息
	//	err, response = remote.RequestLedgerClosed()
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Ledger Closed succ.len=", len(response), "data=", response)
	//
	//	//获取某一账本具体信息
	//	var ledger_index string = "8488670"
	//	var ledger_hash string = ""
	//	var transactions bool = false
	//	err, response = remote.RequestLedger(ledger_index, ledger_hash, transactions)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Ledger succ.len=", len(response), "data=", response)
	//
	//	//获取某一账本具体信息
	//	var hash string = "084C7823C318B8921A362E39C67A6FB15ADA5BCCD0C7E9A3B13485B1EF2A4313"
	//	err, response = remote.RequestTx(hash)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Tx succ.len=", len(response), "data=", response)
	//
	//	//请求账号信息
	//	account := "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	err, response = remote.RequestAccountInfo(map[string]string{"account": account})
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Info succ.len=", len(response), "data=", response)
	//
	//	//获得账号可接收和发送的货币
	//	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	err, response = remote.RequestAccountTums(account)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Tums succ.len=", len(response), "data=", response)
	//
	//	//获得账号交易列表
	//	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	var limit int = 100
	//	err, response = remote.RequestAccountTx(account, limit)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Tx succ.len=", len(response), "data=", response)
	//
	//	//获得账号交易列表
	//	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	atype := "trust"
	//	err, response = remote.RequestAccountRelations(account, atype)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Relations succ.len=", len(response), "data=", response)
	//
	//	atype = "authorize"
	//	err, response = remote.RequestAccountRelations(account, atype)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Relations succ.len=", len(response), "data=", response)
	//
	//	//获得账号挂单
	//	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	err, response = remote.RequestAccountOffers(account)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Offers succ.len=", len(response), "data=", response)
	//
	//	//获得账号挂单
	//	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	//	gets := "SWT"
	//	pays := "CNY"
	//	err, response = remote.RequestOrderBook(account, gets, pays)
	//	if err != nil {
	//		fmt.Println("Get data:", response)
	//		return
	//	}
	//	fmt.Println("Get Response Account Order Book succ.len=", len(response), "data=", response)
	//
	//	/*
	//		//获得账号挂单
	//		err, response = remote.BuildPaymentTx()
	//		if err != nil {
	//			fmt.Println("Get data:", response)
	//			return
	//		}
	//		fmt.Println("Get Response Build Payment Tx succ.len=", len(response), "data=", response)
	//	*/
	//	defer jingtum.Exits()
}
