package jingtumLib

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"strconv"
	"time"

	"jingtumLib/constant"
	jtLRU "jingtumLib/lruCache"
	"jingtumLib/utils"
)

//接收最长报文
var (
	MAX_RECIVE_LEN = 4096000
)

//websocket 连接类
type WsConn struct {
	Ws   *websocket.Conn
	Host string
	Port string
}

/*
* remote是跟井通底层交互最主要的类，它可以组装交易发送到底层、订阅事件及从底层拉取数据。
 */

type Remote struct {
	//	Wsconn   *WsConn
	requests map[uint64]*ReqCtx
	//	Status    bool
	LocalSign bool
	Paths     *jtLRU.LRU
	server    *Server
}

type ReqCtx struct {
	command  string
	data     map[string]interface{}
	callback func(err error, data interface{})
	cid      uint64
}

/*
* remoter 提供以下方法：
 */

type Remoter interface {
	//连接
	Connect() error
	//断开
	//获取当前时间
	Get_now_time() string
	// 发送
	send(request string) error
	// 接收
	read() (error, string)
	//断开连接
	Disconnect()
	//请求底层服务器信息
	RequestServerInfo() (error, string)
	//获取最新账本信息
	RequestLedgerClosed() (error, string)
	//获取某一账本具体信息
	RequestLedger(ledger_index string, ledger_hash string, transactions bool) (error, string)
	//询某一交易具体信息
	RequestTx(hash string) (error, string)
	//请求账号信息
	RequestAccountInfo(account string) (error, string)
	//得账号可接收和发送的货币
	RequestAccountTums(account string) (error, string)
	//得账号关系
	RequestAccountRelations(account string, atype string) (error, string)
	//获得账号挂单
	RequestAccountOffers(account string) (error, string)
	//获得账号交易列表
	RequestAccountTx(account string, limit int) (error, string)
	//获得市场挂单列表
	RequestOrderBook(account string, gets string, pays string) (error, string)

	//创建支付对象
	BuildPaymentTx(account string, to string, amount constant.Amount) (Transaction, error)
}

/*
*   构造函数 带形参
*   params:
*           host 主机名
*           hort  端口号
 */
/*
func NewRemote(host string, port string) (*Remote) {
	remote := new(Remote)
	remote.Wsconn = new(WsConn)
	//remote.Wsconn.Ws = websocket.Conn
	remote.Wsconn.Host = host
	remote.Wsconn.Port = port
	remote.Params = make(map[string]string)
	remote.Status = false
	return remote
}
*/
/*
*   构造函数 不带形参
*   params:
*           从配置文件读取
*           Service|Host
*           Service|Port
 */
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
	lru, err := jtLRU.NewLRU(100, time.Duration(5)*time.Minute, nil)
	if err != nil {
		return remote, err
	}
	remote.Paths = lru
	remote.LocalSign = localSign
	server, err := NewServer(remote, url)
	if err != nil {
		return remote, err
	}

	remote.server = server

	return remote, nil
}

//func NewRemoteByURL(url string, localSign bool) (*Remote, error) {
//	err, remote := NewRemote()
//	if err != nil {
//		return nil, err
//	}
//
//	remote.Wsconn.Host = url
//	remote.Wsconn.Port = port
//	remote.LocalSign = localSign
//	return remote, nil
//}

/*
* 连接函数
 */
func (remote *Remote) Connect(callback func(err error, result interface{})) error {
	if remote.server == nil {
		callback(constant.ERR_SERVER_NOT_READY, nil)
		return constant.ERR_SERVER_NOT_READY
	}

	return remote.server.connect(callback)
}

/*
*   获取当前时间
    params: 无
    return: string
            格式(2006-01-02 15:04:05)
*/
func (remote *Remote) Get_now_time() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

/*
* 发送
*     params:
             request 待发送的报文
*/
//func (remote *Remote) send(request string) error {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err
//		}
//	}
//	if request == "" {
//		return errors.New("Nothing to send")
//	}
//	_, err := remote.Wsconn.Ws.Write([]byte(request))
//	if err != nil {
//		Error("Send ", request, "fail!", "errno:", err)
//	} else {
//		Info("Send: ", request, "succ.")
//	}
//	return err
//}

/*
* 接收
    params: 无
    return:
           接收的报文
*/

//func (remote *Remote) read() (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	var msg = make([]byte, MAX_RECIVE_LEN)
//	var n int
//	n, err := remote.Wsconn.Ws.Read(msg)
//	if err != nil {
//		Error("Received data fail!", "errno:", err)
//	} else {
//		Info("Received: data succ. Len= ", n)
//	}
//	return err, string(msg[:n])
//}

/*
*  断开连接
 */
//func (remote *Remote) Disconnect() {
//	remote.Wsconn.Ws.Close()
//	remote.Status = false
//}

/*
* 请求底层服务器信息
    return:
           response  返回结果
*/
//func (remote *Remote) RequestServerInfo() (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//
//	request := Pack_RequestServerInfo()
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Server Info succ: ", response)
//	return nil, response
//}

/*
* 获取最新账本信息
* return
 */
//func (remote *Remote) RequestLedgerClosed() (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//
//	request := Pack_RequestLedgerClosed()
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Ledger Closed succ: ", response)
//	return nil, response
//}

/*
* 获取某一账本具体信息
 */
//func (remote *Remote) RequestLedger(ledger_index string, ledger_hash string, transactions bool) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//
//	if ledger_index == "" && ledger_hash == "" {
//		return errors.New("ledger_index|ledger_hash value error"), ""
//	}
//	request := Pack_RequestLedger(ledger_index, ledger_hash, transactions)
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Ledger succ: ", response)
//	return nil, response
//}

/*
* 查询某一交易具体信息
 */
//func (remote *Remote) RequestTx(hash string) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestTx(hash)
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Tx succ: ", response)
//	return nil, response
//}

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

/*
* 请求账号信息
 */
func (remote *Remote) RequestAccountInfo(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote)
	req.command = "account_info"
	account, ok := options["account"]
	peer, ok := options["peer"]
	retype, ok := options["type"]
	if ok {
		req.message["relation_type"] = getRelationType(retype.(string))
	}

	limit, ok := options["limit"]
	marker, ok := options["marker"]

	if account != nil {
		req.message["account"] = account.(string)
	}
	ledger, ok := options["ledger"]
	if ok {
		req.SelectLedger(ledger)
	}

	if peer != nil {
		if utils.IsValidAddress(peer.(string)) {
			req.message["peer"] = peer
		}
	}

	var checkedLimit interface{}

	if utils.IsNumberType(limit) {
		if limit.(float64) < 0 {
			checkedLimit = 0
		}

		if limit.(float64) > 1000000000 {
			checkedLimit = 1000000000
		}

	} else {
		if limit != nil {
			if v, ok := limit.(string); ok {
				if utils.IsNumberString(v) {
					lv, err := strconv.ParseFloat(v, 64)
					if err == nil {
						if lv < 0 {
							checkedLimit = 0
						}

						if lv > 1000000000 {
							checkedLimit = 1000000000
						}
					}
				}
			}
		}
	}

	if checkedLimit != nil {
		req.message["limit"] = checkedLimit
	}

	if marker != nil {
		req.message["marker"] = marker
	}

	return req, nil
}

/*
* 获得账号可接收和发送的货币
 */
//func (remote *Remote) RequestAccountTums(account string) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestAccountTums(account)
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Account Tums succ: ", response)
//	return nil, response
//}

/*
* 获得账号关系
 */
//func (remote *Remote) RequestAccountRelations(account string, atype string) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestAccountRelations(account, atype)
//	if request == "" {
//		return errors.New("RequestAccountRelations type value error"), ""
//	}
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Account Relations succ: ", response)
//	return nil, response
//}

/*
* 获得账号挂单
 */
//func (remote *Remote) RequestAccountOffers(account string) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestAccountOffers(account)
//	if request == "" {
//		return errors.New("RequestAccountRelations type value error"), ""
//	}
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Account Offers succ: ", response)
//	return nil, response
//}

/*
* 获得账号交易列表
 */
//func (remote *Remote) RequestAccountTx(account string, limit int) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestAccountTx(account, limit)
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Account Tx succ: ", response)
//	return nil, response
//}

/*
* 获得市场挂单列表
 */
//func (remote *Remote) RequestOrderBook(account string, gets string, pays string) (error, string) {
//	if !remote.Status {
//		err := remote.Connect()
//		if err != nil {
//			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
//			Error("Connect ", host_port, "fail! errno = ", err)
//			return err, ""
//		}
//	}
//	request := Pack_RequestOrderBook(account, gets, pays)
//	err := remote.send(request)
//	if err != nil {
//		Error("Send data fail!")
//		return err, ""
//	}
//	err, response := remote.read()
//	if err != nil {
//		Error("Received data fail!")
//		return err, ""
//	}
//	Info("Get Reqonse Account Tx succ: ", response)
//	return nil, response
//}

/**
 * 提交请求
 */
func (remote *Remote) Submit(command string, data map[string]interface{}, filter Filter, callback func(err error, data interface{})) {
	rc := new(ReqCtx)
	rc.command = command
	rc.data = data
	rc.callback = callback
	rc.cid = remote.server.GetCid()
	remote.requests[rc.cid] = rc

	remote.server.sendMessage(rc)
}

func (remote *Remote) handleResponse(data map[string]interface{}) {
	fmt.Printf("Handle response --> %v \n", data)
	cid, err := strconv.ParseUint(data["id"].(string), 10, 64)
	if err != nil {
		Errorf("Received msg parse id error : %v", err)
		return
	}

	request, ok := remote.requests[cid]

	if !ok {
		Errorf("Request id error %v", cid)
		return
	}

	delete(remote.requests, cid)

	//	 if (data.result && data.status === 'success'
	//            && data.result.server_status) {
	//        this._updateServerStatus(data.result);
	//    }

	status, ok := data["status"]
	if ok {
		if status == "success" {
			//				 var result = request.filter(data.result)
			result, ok := data["result"]
			if ok {
				request.callback(nil, result)
			} else {
				request.callback(errors.New("Response result is null."), nil)
			}

		} else if status == "error" {
			errMsg, ok := data["error_message"]
			errException, oke := data["error_exception"]
			if ok {
				request.callback(errors.New(errMsg.(string)), nil)
			} else if oke {
				request.callback(errors.New(errException.(string)), nil)
			}
		}
	}
}

func (remote *Remote) handlePathFind(data map[string]interface{}) {
	//    this.emit('path_find', data);
}

func (remote *Remote) handleTransaction(data map[string]interface{}) {
	//    var self = this;
	//    var tx = data.transaction.hash;
	//    if (self._cache.get(tx)) return;
	//    self._cache.set(tx, 1);
	//    this.emit('transactions', data);
}

func (remote *Remote) handleServerStatus(data map[string]interface{}) {
	// TODO check data format
	//    this._updateServerStatus(data);
	//    this.emit('server_status', data);
}

func (remote *Remote) handleLedgerClosed(data map[string]interface{}) {
	//    var self = this;
	//    if (data.ledger_index > self._status.ledger_index) {
	//        self._status.ledger_index = data.ledger_index;
	//        self._status.ledger_time = data.ledger_time;
	//        self._status.reserve_base = data.reserve_base;
	//        self._status.reserve_inc = data.reserve_inc;
	//        self._status.fee_base = data.fee_base;
	//        self._status.fee_ref = data.fee_ref;
	//        self.emit('ledger_closed', data);
	//    }
}

//消息处理方法
func (remote *Remote) handleMessage(msg []byte) {
	var data map[string]interface{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		Errorf("Received msg json Unmarshal error : %v", err)
		return
	}

	_, ok := data["error"]
	if ok {
		cid, err := strconv.ParseUint(data["id"].(string), 10, 64)
		if err != nil {
			Errorf("Received msg parse id error : %v", err)
			return
		}
		remote.requests[cid].callback(errors.New(data["error_message"].(string)), nil)
	} else {
		resType := data["type"].(string)
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
}

/**
 * 创建支付对象
 */
func (remote *Remote) BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error) {
	tx, err := NewTransaction(remote)
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

	tx.AddTxJson("TransactionType", "Payment")
	tx.AddTxJson("Account", account)

	toamount, err := utils.ToAmount(amount)

	if err != nil {
		return nil, err
	}

	tx.AddTxJson("Amount", toamount)
	tx.AddTxJson("Destination", to)

	return tx, nil
}

/*
* 支付
 */
/*
func (remote *Remote) BuildPaymentTx(account string) (error, string) {
    if !remote.Status {
        err := remote.Connect()
        if if err != nil {
            host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
            Error("Connect ", host_port, "fail! errno = ", err)
            return err, ""
        }
    }
	request := fmt.Sprintf("{\"id\":%d,\"command\":\"submit\",\"secret\":\"shXKYEcFwmxPSCeu4rRUz1ECEYtZP\",\"tx_json\":{\"Flags\":0,\"TransactionType\":\"Payment\",\"Account\":\"jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq\",\"Amount\":\"1000000\",\"Destination\":\"jQpP2UpTxECw74kuzrMXB6YEU3jsrnDDsc\"}}", seq)
	err := send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Build Payment Tx succ: ", response)
	return nil, response
}*/
