package jingtumLib

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"jingtumLib/constant"
	jtLRU "jingtumLib/lruCache"
	"jingtumLib/utils"

	"github.com/olebedev/emitter"
)

var (
	//MaxReciveLen 接收最长报文
	MaxReciveLen = 4096000
)

//Remote 是跟井通底层交互最主要的类，它可以组装交易发送到底层、订阅事件及从底层拉取数据。
type Remote struct {
	requests  map[uint64]*ReqCtx
	status    map[string]interface{}
	LocalSign bool
	Paths     *jtLRU.LRU
	cache     *jtLRU.LRU
	server    *Server
	emit      *emitter.Emitter
}

type ResData map[string]interface{}

type ParameterInfo struct {
	Parameter string
}

type ArgInfo struct {
	Arg *ParameterInfo
}

//ReqCtx 请求包装类
type ReqCtx struct {
	command  string
	data     map[string]interface{}
	callback func(err error, data interface{})
	cid      uint64
	filter   Filter
}

//Remoter 提供以下方法：
type Remoter interface {
	//Connect 连接
	Connect(callback func(err error, result interface{})) error

	//GetNowTime 获取当前时间
	GetNowTime() string

	//断开连接
	Disconnect()

	//RequestServerInfo 请求底层服务器信息
	RequestServerInfo() (*Request, error)

	//RequestLedgerClosed 获取最新账本信息
	RequestLedgerClosed() (*Request, error)

	//获取某一账本具体信息
	RequestLedger(options map[string]interface{}) (*Request, error)

	//RequestTx 询某一交易具体信息
	RequestTx(hash string) (*Request, error)

	//请求账号信息
	RequestAccountInfo(options map[string]interface{}) (*Request, error)

	//RequestAccountTums 得账号可接收和发送的货币
	RequestAccountTums(options map[string]interface{}) (*Request, error)

	//RequestAccountRelations 得账号关系
	RequestAccountRelations(options map[string]interface{}) (*Request, error)

	//RequestAccountOffers 获得账号挂单
	RequestAccountOffers(options map[string]interface{}) (*Request, error)

	//RequestAccountTx 获得账号交易列表
	RequestAccountTx(options map[string]interface{}) (*Request, error)

	//RequestOrderBook 获得市场挂单列表
	RequestOrderBook(options map[string]interface{}) (*Request, error)

	//BuildPaymentTx 创建支付对象
	BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error)
}

//NewRemote 创建Remote，url 为空是从配置文件获取server 地址
func NewRemote(url string, localSign bool) (*Remote, error) {
	remote := new(Remote)

	if url == "" {
		url = JTConfig.Read("Service", "Host")

		if url == "" {
			Error("Config Service:Host is null.")
			return remote, errors.New("Config|service:Host setting error")
		}

		port := JTConfig.Read("Service", "Port")

		if port == "" {
			Error("Config Service:Port is null.")
			return remote, errors.New("Config|service:Port setting error")
		}

		url += ":" + port
	}

	remote.requests = make(map[uint64]*ReqCtx)
	remote.status = make(map[string]interface{})
	lru, err := jtLRU.NewLRU(100, time.Duration(5)*time.Minute, nil)
	if err != nil {
		return remote, err
	}
	remote.Paths = lru

	remote.cache, err = jtLRU.NewLRU(100, time.Duration(5)*time.Minute, nil)
	if err != nil {
		return remote, err
	}
	remote.LocalSign = localSign
	server, err := NewServer(remote, url)
	if err != nil {
		return remote, err
	}

	remote.server = server
	remote.emit = &emitter.Emitter{}
	remote.emit.Use("*", emitter.Void)

	return remote, nil
}

//Connect 连接函数
func (remote *Remote) Connect(callback func(err error, result interface{})) error {
	if remote.server == nil {
		callback(constant.ERR_SERVER_NOT_READY, nil)
		return constant.ERR_SERVER_NOT_READY
	}

	return remote.server.connect(callback)
}

//GetNowTime 获取当前时间。格式(2006-01-02 15:04:05)
func (remote *Remote) GetNowTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

//Disconnect 关闭连接
func (remote *Remote) Disconnect() {
	if remote.server != nil && remote.server.Disconnect() {
		//清除请求缓存
		for id := range remote.requests {
			delete(remote.requests, id)
		}
	}
}

//RequestServerInfo 请求底层服务器信息
func (remote *Remote) RequestServerInfo() (*Request, error) {
	req := NewRequest(remote, constant.CommandServerInfo, func(data interface{}) interface{} {
		info := data.(map[string]interface{})["info"].(map[string]interface{})
		retData := map[string]interface{}{"version": "skywelld-" + info["build_version"].(string), "peers": info["peers"], "state": info["server_state"], "public_key": info["pubkey_node"], "complete_ledgers": info["complete_ledgers"], "ledger": info["validated_ledger"].(map[string]interface{})["hash"]}

		return retData
	})

	return req, nil
}

//RequestLedgerClosed 获取最新账本信息
func (remote *Remote) RequestLedgerClosed() (*Request, error) {
	req := NewRequest(remote, constant.CommandLedgerClosed, func(data interface{}) interface{} {
		retData := map[string]interface{}{"ledger_hash": data.(map[string]interface{})["ledger_hash"], "ledger_index": data.(map[string]interface{})["ledger_index"]}
		return retData
	})
	return req, nil
}

//RequestLedger 获取某一账本具体信息.
func (remote *Remote) RequestLedger(options map[string]interface{}) (*Request, error) {
	isFilter := true
	req := NewRequest(remote, constant.CommandLedger, func(data interface{}) interface{} {
		ledger, ok := data.(map[string]interface{})["ledger"]
		if !ok {
			if closed, ok := data.(map[string]interface{})["closed"]; ok {
				ledger, ok = closed.(map[string]interface{})["ledger"]
			}
		}
		if !isFilter {
			return ledger
		}

		if ledger == nil {
			return nil
		}

		retData := map[string]interface{}{"accepted": ledger.(map[string]interface{})["accepted"], "ledger_hash": ledger.(map[string]interface{})["hash"], "ledger_index": ledger.(map[string]interface{})["ledger_index"], "parent_hash": ledger.(map[string]interface{})["parent_hash"], "close_time": ledger.(map[string]interface{})["close_time_human"], "total_coins": ledger.(map[string]interface{})["total_coins"]}
		return retData
	})

	if ledgerIndex, ok := options["ledger_index"].(string); ok && utils.MatchString("^\\d+$", ledgerIndex) {
		ledgerIndexNum, err := strconv.Atoi(ledgerIndex)
		if err != nil {
			return nil, err
		}
		req.message["ledger_index"] = ledgerIndexNum
	}

	if ledgerHash, ok := options["ledger_hash"].(string); ok && utils.MatchString("^[A-F0-9]{64}$", ledgerHash) {
		req.message["ledger_hash"] = ledgerHash
	}

	if transactions, ok := options["transactions"].(bool); ok {
		req.message["transactions"] = transactions
		isFilter = false
	}

	return req, nil
}

//RequestTx 查询某一交易具体信息
func (remote *Remote) RequestTx(hash string) (*Request, error) {
	if hash == "" || !utils.MatchString("^[A-F0-9]{64}$", hash) {
		return nil, fmt.Errorf("Invalid tx hash")
	}

	req := NewRequest(remote, constant.CommandTX, nil)
	req.message["transaction"] = hash
	return req, nil
}

func getRelationType(relationType string) *constant.Integer {
	switch relationType {
	case "trustline":
		return constant.NewInteger(0)
	case "authorize":
		return constant.NewInteger(1)
	case "freeze":
		return constant.NewInteger(3)

	}
	return nil
}

func requestAccount(req *Request, options map[string]interface{}) {
	if retype, ok := options["type"].(string); ok {
		relationType := getRelationType(retype)
		if relationType != nil {
			req.message["relation_type"] = relationType.IntValue()
		}
	}

	if account, ok := options["account"].(string); ok {
		req.message["account"] = account
	}

	ledger, _ := options["ledger"]
	req.SelectLedger(ledger)

	if peer, ok := options["peer"].(string); ok {
		if utils.IsValidAddress(peer) {
			req.message["peer"] = peer
		}
	}

	if limit, ok := options["limit"].(int); ok {
		if limit < 0 {
			limit = 0
		}

		if limit > 1000000000 {
			limit = 1000000000
		}

		req.message["limit"] = limit

	}

	if marker, ok := options["marker"]; ok {
		req.message["marker"] = marker
	}
}

//RequestAccountInfo 请求账号信息
func (remote *Remote) RequestAccountInfo(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountInfo
	requestAccount(req, options)
	return req, nil
}

//RequestAccountTums 获得账号可接收和发送的货币
func (remote *Remote) RequestAccountTums(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountCurrencies
	requestAccount(req, options)
	return req, nil
}

//RequestAccountRelations 获得账号关系
func (remote *Remote) RequestAccountRelations(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	rtype, ok := options["type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid realtion type")
	}

	if _, okType := constant.RelationTypes[rtype]; !okType {
		return nil, fmt.Errorf("invalid realtion type %s", rtype)
	}

	switch rtype {
	case "trust":
		req.command = constant.CommandAccountLines
	case "authorize", "freeze":
		req.command = constant.CommandAccountRelation
	default:
		return nil, fmt.Errorf("relation should not go here %s", rtype)
	}

	requestAccount(req, options)
	return req, nil
}

//RequestAccountOffers 获得账号挂单
func (remote *Remote) RequestAccountOffers(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountOffers
	requestAccount(req, options)
	return req, nil
}

//RequestAccountTx 获得账号交易列表
func (remote *Remote) RequestAccountTx(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, constant.CommandAccountTX, func(data interface{}) interface{} {
		//過濾交易列表
		return data
	})

	if _, ok := options["limit"]; !ok {
		options["limit"] = 200
	}

	if account, ok := options["account"].(string); ok {
		if !utils.IsValidAddress(account) {
			return nil, fmt.Errorf("account parameter is invalid %s", account)
		}
		req.message["account"] = account
	}

	if ledgerMin, ok := options["ledger_min"].(int); ok {
		req.message["ledger_index_min"] = ledgerMin
	} else {
		req.message["ledger_index_min"] = 0
	}
	if ledgerMax, ok := options["ledger_max"].(int); ok {
		req.message["ledger_index_max"] = ledgerMax
	} else {
		req.message["ledger_index_max"] = -1
	}

	if limit, ok := options["limit"].(int); ok {
		req.message["limit"] = limit
	}

	if offset, ok := options["offset"].(int); ok {
		req.message["offset"] = offset
	}

	if marker, ok := options["offset"].(map[string]interface{}); ok {
		if _, ok = marker["ledger"].(int); ok {
			if _, ok = marker["seq"].(int); ok {
				req.message["marker"] = marker
			}
		}
	}
	if forward, ok := options["forward"].(bool); ok {
		//true 正向；false反向
		req.message["forward"] = forward
	}
	return req, nil
}

//RequestOrderBook 获得市场挂单列表
func (remote *Remote) RequestOrderBook(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, constant.CommandBookOffers, nil)

	if takerGets, ok := options["taker_gets"]; ok {
		getsAmount, ok := takerGets.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid taker_gets type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&getsAmount) {
			return nil, fmt.Errorf("invalid taker gets amount")
		}
		req.message["taker_gets"] = getsAmount
	} else if pays, ok := options["pays"]; ok {
		paysAmount, ok := pays.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid pays type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&paysAmount) {
			return nil, fmt.Errorf("invalid taker gets amount")
		}
		req.message["taker_gets"] = paysAmount
	}

	if takerPays, ok := options["taker_pays"]; ok {
		paysAmount, ok := takerPays.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid taker_pays type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&paysAmount) {
			return nil, fmt.Errorf("invalid taker pays amount")
		}
		req.message["taker_pays"] = paysAmount

	} else if gets, ok := options["gets"]; ok {
		getsAmount, ok := gets.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid gets type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&getsAmount) {
			return nil, fmt.Errorf("invalid gets amount")
		}
		req.message["taker_pays"] = getsAmount
	}

	if limit, ok := options["limit"].(int); ok {
		req.message["limit"] = limit
	}

	if taker, ok := options["taker"]; ok {
		req.message["taker"] = taker
	} else {
		req.message["taker"] = constant.AccountOne
	}
	return req, nil
}

//Subscribe 订阅服务
func (remote *Remote) Subscribe(streams []string) *Request {
	req := NewRequest(remote, constant.CommandSubscribe, nil)

	if len(streams) > 0 {
		req.message["streams"] = streams
	}
	return req
}

//UnSubscribe 退订服务
func (remote *Remote) UnSubscribe(streams []string) *Request {
	req := NewRequest(remote, constant.CommandUnSubscribe, nil)
	if len(streams) > 0 {
		req.message["streams"] = streams
	}
	return req
}

//Submit 提交请求
func (remote *Remote) Submit(command string, data map[string]interface{}, filter Filter, callback func(err error, data interface{})) {
	rc := new(ReqCtx)
	rc.command = command
	rc.data = data
	rc.callback = callback
	rc.filter = filter
	rc.cid = remote.server.GetCid()
	remote.requests[rc.cid] = rc

	remote.server.sendMessage(rc)
}

//On 监听特定的事件消息
func (remote *Remote) On(eventName string, callback func(data interface{})) {
	remote.emit.On(eventName, func(event *emitter.Event) {
		if len(event.Args) > 0 {
			callback(event.Args[0])
		}
	})
}

func (remote *Remote) handleResponse(data ResData) {
	request, ok := remote.requests[data.getUint64("id")]

	if !ok {
		Errorf("Request id error %d", data.getUint64("id"))
		return
	}

	delete(remote.requests, data.getUint64("id"))

	if data.getString("status") == "success" {
		result := request.filter(data.getMap("result"))
		request.callback(nil, result)
	} else if data.getString("status") == "error" {
		errMsg := data.getString("error_message")
		if errMsg != "" {
			request.callback(errors.New(errMsg), nil)
		}
	}
}

func (remote *Remote) handlePathFind(data ResData) {
	go remote.emit.Emit(constant.EventPathFind, data)
}

func (remote *Remote) handleTransaction(data ResData) {
	if txHash, ok := data.getMap("transaction")["hash"].(string); ok {
		remote.cache.Add(txHash, 1)
		go remote.emit.Emit(constant.EventTX, data)
	}
}

func (remote *Remote) updateServerStatus(data ResData) {
	remote.status["load_base"] = data.getObj("load_base")
	remote.status["load_factor"] = data.getObj("load_factor")
	if data.getObj("pubkey_node") != nil {
		remote.status["pubkey_node"] = data.getObj("pubkey_node")
	}
	remote.status["server_status"] = data.getObj("server_status")
	serverStatus := data.getString("server_status")
	online := "offline"
	if onlineStates.contain(serverStatus) {
		online = "online"
	}
	remote.server.setState(online)
}

func (remote *Remote) handleServerStatus(data ResData) {
	remote.updateServerStatus(data)
	go remote.emit.Emit(constant.EventServerStatus, data)
}

func (remote *Remote) handleLedgerClosed(data ResData) {
	stsIdx, ok := remote.status["ledger_index"]
	if !ok {
		remote.status["ledger_index"] = data.getFloat64("ledger_index")
		go remote.emit.Emit(constant.EventLedgerClosed, data)
	} else if data.getFloat64("ledger_index") > stsIdx.(float64) {
		remote.status["ledger_time"] = data.getObj("ledger_time")
		remote.status["reserve_base"] = data.getObj("reserve_base")
		remote.status["reserve_inc"] = data.getObj("reserve_inc")
		remote.status["fee_base"] = data.getObj("fee_base")
		remote.status["fee_ref"] = data.getObj("fee_ref")
		go remote.emit.Emit(constant.EventLedgerClosed, data)
	}
}

//消息处理方法
func (remote *Remote) handleMessage(msg []byte) {
	var data ResData
	err := json.Unmarshal(msg, &data)
	if err != nil {
		Errorf("Received msg json Unmarshal error : %v", err)
		return
	}

	if data.getString("error") != "" {
		remote.requests[data.getUint64("id")].callback(errors.New(data.getString("error_message")), nil)
	} else {
		resType := data.getString("type")
		switch resType {
		case "ledgerClosed":
			remote.handleLedgerClosed(data)
		case "serverStatus":
			remote.handleServerStatus(data)
		case "response":
			remote.handleResponse(data)
		case "transaction":
			remote.handleTransaction(data)
		case "path_find":
			remote.handlePathFind(data)
		}
	}
	// }
}

//BuildPaymentTx 创建支付对象
func (remote *Remote) BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error) {
	tx, err := NewTransaction(remote, nil)
	if err != nil {
		return nil, err
	}

	if !utils.IsValidAddress(account) {
		return nil, constant.ERR_PAYMENT_INVALID_SRC_ADDR
	}

	if !utils.IsValidAddress(to) {
		return nil, constant.ERR_PAYMENT_INVALID_DST_ADDR
	}

	if !utils.IsValidAmount(&amount) {
		return nil, constant.ERR_PAYMENT_INVALID_AMOUNT
	}

	tx.AddTxJSON("TransactionType", "Payment")
	tx.AddTxJSON("Account", account)

	toamount, err := utils.ToAmount(amount)

	if err != nil {
		return nil, err
	}

	tx.AddTxJSON("Amount", toamount)
	tx.AddTxJSON("Destination", to)

	return tx, nil
}

//BuildRelationSet
// func (remote *Remote) BuildRelationSetBuildRelationSet(options map[string]interface{}) (*Transaction, error) {
// {
// 	src :=  options["source"]
// 	if src == "" {
// 		src = options["from"]
// 	}
// 	if src == "" {
// 		src = options["account"]
// 	}

// 	des := options["target"]
// 	limit := options["limit"]
// 	if !utils.isValidAddress(src.(string)) {
// 		return tx, Error("invalid source address")
// 	}

// 	if !utils.isValidAddress(des.(string)) {
// 		return tx, Error("invalid target address")
// 	}

// 	if !utils.isValidAmount(limit) {
// 		return tx, Error("invalid amount")
// 	}
// 	transactionType :=  ""
// 	if options["type"] == "unfreeze" {
// 		transactionType = "RelationDel"
// 	} else {
// 		transactionType = "RelationSet"
// 	}
// 	tx.AddTxJSON("TransactionType", transactionType)
// 	tx.AddTxJSON("AccountAccount", src)
// 	tx.AddTxJSON("Target", des)
// 	relationType := ""
// 	if options["type"] == "authorize" {
// 		relationType = "1"
// 	} else {
// 		relationType = "3"
// 	}
// 	tx.AddTxJSON("RelationType", relationType)
// 	if limit != 0 {
// 		tx.AddTxJSON("LimitAmount", limit)
// 	}
// 	return tx, nil
// }

// func(remote *Remote) BuildTrustSet(options map[string]interface{}) (*Transaction, error) {
// 	tx, err := NewTransaction(remote, nil)
// 	if src == "" {
// 		src = options["from"]
// 	}
// 	if src == "" {
// 		src = options.account
// 	}
// 	limit := options.limit
// 	quality_out := options.quality_out
// 	quality_in := options.quality_in
// 	if !utils.isValidAddress(src) {
// 		return tx, Error("invalid source address")
// 	}
// 	if !utils.isValidAmount(limit) {
// 		return tx, Error("invalid amount")
// 	}
// 	tx.AddTxJSON("TrustSet", TransactionType)
// 	tx.AddTxJSON("Account", src)
// 	if limit != 0 {
// 		tx.AddTxJSON("LimitAmount", limit)
// 	}
// 	if quality_in {
// 		tx.AddTxJSON("QualityInQualityIn", quality_in)
// 	}
// 	if quality_out {
// 		tx.AddTxJSON("QualityOut", quality_out)
// 	}
// 	return tx,nil
// }

// //创建关系对象
// func (remote *Remote) BuildRelationTx(options map[string]interface{}) (*Transaction, error) {
// 	tx, err := NewTransaction(remote, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if _, ok := RelationTypes[options.type]; !ok {
// 		return tx, Error("invalid relation type")
// 	}

// 	switch (options.type) {
// 		case "trust":
// 			return BuildTrustSet(options, tx)
// 		case "authorize":
// 		case "freeze":
// 		case "unfreeze":
// 			return BuildRelationSet(options, tx)
// 	}
// 	Errorf("build relation set should not go here")
// 	return tx, Error("build relation set error")
// }

// //BuildAccountSet
// func (remote *Remote) BuildAccountSet(options map[string]interface{}) (*Transaction, error) {
// 	src := options.source
// 	if options.source == "" {
// 		src = options.from
// 	}
// 	if src == "" {
// 		src = options.account
// 	}
//     set_flag := options.set_flag
// 	if options.set_flag == "" {
// 		set_flag = options.set
// 	}
//     clear_flag := options.clear_flag
// 	if clear_flag == "" {
// 		clear_flag = options.clear
// 	}
//     if (!utils.isValidAddress(src)) {
//         return tx, Error("invalid source address")
//     }
//     tx.AddTxJSON("TransactionType", "AccountSet")
//     tx.AddTxJSON("Account", src)

//     SetClearFlags := Set_clear_flags[1]
// 	_set_flag := ""
// 	if IsNumberType(set_flag) {
// 		_set_flag = set_flag
// 	}
// 	else if SetClearFlags[set_flag] == "" {
// 		_set_flag = SetClearFlags["asf" + set_flag]
// 	}
// 	else {
// 		_set_flag = SetClearFlags[set_flag]
// 	}

// 	if set_flag == "" {
// 		set_flag = _set_flag
// 	}
// 	tx.AddTxJSON("SetFlag", set_flag)

// 	_clear_flag := ""
// 	if IsNumberType(clear_flag) {
// 		_clear_flag = clear_flag
// 	}
// 	else if SetClearFlags[clear_flag] == "" {
// 		_clear_flag = SetClearFlags["asf" + clear_flag]
// 	}
// 	else {
// 		_clear_flag =  SetClearFlags[clear_flag]
// 	}
//     if clear_flag == "" {
// 		clear_flag = _clear_flag
// 	}
// 	tx.addTxJSON("ClearFlag", clear_flag)
// 	return tx, nil
// }

// //BuildDelegateKeySet
// func (remote *Remote) BuildDelegateKeySet(options map[string]interface{}) (*Transaction, error) {
// 	src := options.source
// 	if options.source == "" {
// 		src = options.from
// 	}
// 	if src = "" {
// 		src = options.account
// 	}
// 	delegate_key := options.delegate_key
// 	if !utils.isValidAddress(src) {
// 		return tx, Error("invalid source address")
// 	}
// 	if !utils.isValidAddress(delegate_key) {
// 		return tx, Error("invalid regular key address")
// 	}
// 	tx.addTxJSON("TransactionType", "SetRegularKey")
// 	tx.addTxJSON("Account", src)
// 	tx.addTxJSON("RegularKey", delegate_key)
// 	return tx, nil
// }

// //BuildSignerSet
// func (remote *Remote) BuildSignerSet(options map[string]interface{}) (*Transaction, error) {
// 	// TODO
// 	tx, err := NewTransaction(remote, nil)
// 	return tx, nil
// }

// //创建属性对象
// func (remote *Remote) BuildAccountSetTx(options map[string]interface{}) (*Transaction, error) {
// 	tx, err := NewTransaction(remote, nil)
// 	if _, ok := AccountSetTypes[options.type]; !ok {
// 		return tx, Error("invalid account set type")
// 	}

//     switch(options.type) {
//         case "property":
//             return BuildAccountSet(options, tx)
//         case "delegate":
//             return BuildDelegateKeySet(options, tx)
//         case "signer":
//             return BuildSignerSet(options, tx)
//     }

// 	Errorf("build account set should not go here")
//     return tx, Error("build account set should not go here")
// }

// //挂单
// func (remote *Remote) BuildOfferCreateTx(options map[string]interface{}) (*Transaction, error) {
//  	tx := NewTransaction(remote, nil)
//     offer_type := options.type;
// 	src := options.source
//     if options.source == "" {
//         src = options.from
//     }
//     if src = "" {
//         src = options.account
//     }
//     taker_gets := options.taker_gets
// 	if taker_gets == "" {
// 		taker_gets = options.pays
// 	}
//     taker_pays := options.taker_pays
// 	if taker_pays == ""{
// 		taker_pays = options.gets;
// 	}
//     if (!utils.isValidAddress(src)) {
//         return tx, Error("invalid source address");
//     }

// 	if _, ok := OfferTypesOfferTypes[options.type]; !ok {
// 		return tx, Error("invalid offer type")
// 	}
// 	if !IsStringType(offer_typroffer_typr) {
// 		 return tx, Error("invalid offer type")
// 	}

// 	if IsStringType(taker_gets) && !IsNumberString(taker_gets) {
// 		return tx, Error("invalid to pays amount")
// 	}
// 	if  !utils.isValidAmount(taker_gets) {
// 		return tx, Error("invalid to pays amount object")
// 	}

// 	if IsStringType(taker_pays) && !IsNumberString(taker_pays){
// 		return tx, Error("invalid to gets amount")
// 	}
// 	if !utils.isValidAmount(taker_pays) {
// 		return tx, Error("invalid to gets amount object")
// 	}

//     tx.AddTxJSON("TransactionType", "OfferCreate")
//     if (offer_type == "Sell") {
// 		tx.SetFlags(offer_type)
// 	}
//     tx.AddTxJSON("Account", src)
//     tx.AddTxJSON("TakerPays", utils.ToAmount(taker_pays)
//     tx.AddTxJSON("TakerGets", utils.ToAmount(taker_gets)
//     return tx, nil;
// }

// //取消挂单
// func (remote *Remote) BuildOfferCancelTx(options map[string]interface{}) (*Transaction, error) {
// {
// 	tx := NewTransaction(remote, nil)
// 	src := options.source
// 	if options.source == "" {
// 		src = options.from
// 	}
// 	if src = "" {
// 		src = options.account
// 	}
// 	sequence := options.sequence;
// 	if !utils.isValidAddress(src) {
// 		return tx, Error("invalid source address");
// 	}

// 	if !utils.IsNumberString(sequence) {
// 		return tx, Error("invalid sequence param");
// 	}
// 	tx.AddTxJSON("TransactionType", "OfferCancel")
// 	tx.AddTxJSON("Account", src)
// 	tx.AddTxJSON("OfferSequence", trconv.Atoi(sequencesequence))
// 	return tx;
// }

//DeployContractTx 部署合约
func (remote *Remote) DeployContractTx(options map[string]interface{}) (*Transaction, error) {
	tx, err := NewTransaction(remote, nil)
	if err != nil {
		return nil, err
	}

	if account, ok := options["account"].(string); ok {
		if !utils.IsValidAddress(account) {
			return tx, fmt.Errorf("invalid address %s", account)
		}
		tx.AddTxJSON("Account", account)
	} else {
		return tx, fmt.Errorf("invalid address")
	}

	if amount, ok := options["amount"]; ok {
		if amtStr, ok := amount.(string); ok {
			amtFlat, err := strconv.ParseFloat(amtStr, 64)
			if err != nil {
				return tx, err
			}

			tx.AddTxJSON("Amount", (amtFlat * 1000000))
		} else if amtFlat64, ok := amount.(float64); ok {
			tx.AddTxJSON("Amount", (amtFlat64 * 1000000))
		} else {
			return tx, fmt.Errorf("amount type must be float64 or string")
		}
	}

	if payload, ok := options["payload"].(string); ok {
		tx.AddTxJSON("Payload", payload)
	} else {
		return tx, fmt.Errorf("invalid payload: type error")
	}

	if params, ok := options["params"]; ok {
		if paramArray, ok := params.([]string); ok {
			args := list.New() //[]map[string]string
			for _, v := range paramArray {
				argInfo := new(ArgInfo)
				obj := &ParameterInfo{Parameter: fmt.Sprintf("%X", v)}
				argInfo.Arg = obj
				args.PushBack(argInfo)
			}
			tx.AddTxJSON("Args", args)
		} else {
			return tx, fmt.Errorf("invalid options type")
		}
	}
	tx.AddTxJSON("TransactionType", "ConfigContract")
	tx.AddTxJSON("Method", 0)
	return tx, nil
}

//CallContractTx 执行合约
func (remote *Remote) CallContractTx(options map[string]interface{}) (*Transaction, error) {
	tx, err := NewTransaction(remote, nil)
	if err != nil {
		return nil, err
	}
	if account, ok := options["account"].(string); ok {
		if !utils.IsValidAddress(account) {
			return tx, fmt.Errorf("invalid address %s", account)
		}

		tx.AddTxJSON("Account", account)
	} else {
		return tx, fmt.Errorf("invalid address")
	}

	if des, ok := options["destination"].(string); ok {
		if !utils.IsValidAddress(des) {
			return tx, fmt.Errorf("invalid destination %s", des)
		}
		tx.AddTxJSON("Destination", des)
	}

	if params, ok := options["params"]; ok {
		if paramArray, ok := params.([]string); ok {
			var Args []map[string]string
			for _, v := range paramArray {
				obj := make(map[string]string)
				obj["Parameter"] = fmt.Sprintf("%X", v)
				Args = append(Args, obj)
			}
			tx.AddTxJSON("Args", Args)
		} else {
			return tx, fmt.Errorf("invalid options type")
		}
	}

	if foo, ok := options["foo"].(string); ok {
		tx.AddTxJSON("ContractMethod", fmt.Sprintf("%X", foo))
	} else {
		return tx, fmt.Errorf("foo must be string")
	}

	tx.AddTxJSON("TransactionType", "ConfigContract")
	tx.AddTxJSON("Method", 1)
	return tx, nil
}

func (resData ResData) getUint64(key string) uint64 {
	if ret, ok := (resData)[key]; ok {
		switch v := ret.(type) {
		case float64:
			return uint64(v)
		case float32:
			return uint64(v)
		case int:
			return uint64(v)
		case int8:
			return uint64(v)
		case int32:
			return uint64(v)
		case int64:
			return uint64(v)
		case uint:
			return uint64(v)
		case uint8:
			return uint64(v)
		case uint32:
			return uint64(v)
		}
	}
	return 0
}

func (resData ResData) getString(key string) string {
	if ret, ok := resData[key].(string); ok {
		return ret
	}
	return ""
}

func (resData ResData) getMap(key string) map[string]interface{} {
	if ret, ok := resData[key].(map[string]interface{}); ok {
		return ret
	}
	return nil
}

func (resData ResData) getFloat64(key string) float64 {
	if ret, ok := resData[key].(float64); ok {
		return ret
	}
	return 0
}

func (resData ResData) getObj(key string) interface{} {
	if ret, ok := resData[key]; ok {
		return ret
	}
	return nil
}
