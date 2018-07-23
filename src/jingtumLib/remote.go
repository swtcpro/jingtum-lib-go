package jingtumLib

import (
	"errors"
	"golang.org/x/net/websocket"
	"time"
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
	Wsconn    *WsConn
	Params    map[string]string
	Status    bool
	LocalSign bool
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
func NewRemote() (error, *Remote) {
	remote := new(Remote)
	remote.Wsconn = new(WsConn)
	remote.Wsconn.Host = JTConfig.Read("Service", "Host")
	if remote.Wsconn.Host == "" {
		Error("Config Service:Host is null.")
		return errors.New("Config|service:Host setting error"), remote
	}
	remote.Wsconn.Port = JTConfig.Read("Service", "Port")
	if remote.Wsconn.Port == "" {
		Error("Config Service:Port is null.")
		return errors.New("Config|service:Port setting error"), remote
	}
	remote.Params = make(map[string]string)
	remote.Status = false
	return nil, remote
}

/*
* 连接函数
 */
func (remote *Remote) Connect() error {
	if remote.Status {
		return nil
	}
	host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
	origin := "http://localhost/"
	ws, err := websocket.Dial(host_port, "", origin)
	if err != nil {
		Error("Connect ", host_port, "fail! errno = ", err)
	} else {
		Info("Connect ", host_port, " succ. ")
		remote.Status = true
	}
	remote.Wsconn.Ws = ws
	return err
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
func (remote *Remote) send(request string) error {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err
		}
	}
	if request == "" {
		return errors.New("Nothing to send")
	}
	_, err := remote.Wsconn.Ws.Write([]byte(request))
	if err != nil {
		Error("Send ", request, "fail!", "errno:", err)
	} else {
		Info("Send: ", request, "succ.")
	}
	return err
}

/*
* 接收
    params: 无
    return:
           接收的报文
*/

func (remote *Remote) read() (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	var msg = make([]byte, MAX_RECIVE_LEN)
	var n int
	n, err := remote.Wsconn.Ws.Read(msg)
	if err != nil {
		Error("Received data fail!", "errno:", err)
	} else {
		Info("Received: data succ. Len= ", n)
	}
	return err, string(msg[:n])
}

/*
*  断开连接
 */
func (remote *Remote) Disconnect() {
	remote.Wsconn.Ws.Close()
	remote.Status = false
}

/*
* 请求底层服务器信息
    return:
           response  返回结果
*/
func (remote *Remote) RequestServerInfo() (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}

	request := Pack_RequestServerInfo()
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Server Info succ: ", response)
	return nil, response
}

/*
* 获取最新账本信息
* return
 */
func (remote *Remote) RequestLedgerClosed() (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}

	request := Pack_RequestLedgerClosed()
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Ledger Closed succ: ", response)
	return nil, response
}

/*
* 获取某一账本具体信息
 */
func (remote *Remote) RequestLedger(ledger_index string, ledger_hash string, transactions bool) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}

	if ledger_index == "" && ledger_hash == "" {
		return errors.New("ledger_index|ledger_hash value error"), ""
	}
	request := Pack_RequestLedger(ledger_index, ledger_hash, transactions)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Ledger succ: ", response)
	return nil, response
}

/*
* 查询某一交易具体信息
 */
func (remote *Remote) RequestTx(hash string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestTx(hash)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Tx succ: ", response)
	return nil, response
}

/*
* 请求账号信息
 */
func (remote *Remote) RequestAccountInfo(account string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestAccountInfo(account)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Info succ: ", response)
	return nil, response
}

/*
* 获得账号可接收和发送的货币
 */
func (remote *Remote) RequestAccountTums(account string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestAccountTums(account)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Tums succ: ", response)
	return nil, response
}

/*
* 获得账号关系
 */
func (remote *Remote) RequestAccountRelations(account string, atype string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestAccountRelations(account, atype)
	if request == "" {
		return errors.New("RequestAccountRelations type value error"), ""
	}
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Relations succ: ", response)
	return nil, response
}

/*
* 获得账号挂单
 */
func (remote *Remote) RequestAccountOffers(account string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestAccountOffers(account)
	if request == "" {
		return errors.New("RequestAccountRelations type value error"), ""
	}
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Offers succ: ", response)
	return nil, response
}

/*
* 获得账号交易列表
 */
func (remote *Remote) RequestAccountTx(account string, limit int) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestAccountTx(account, limit)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Tx succ: ", response)
	return nil, response
}

/*
* 获得市场挂单列表
 */
func (remote *Remote) RequestOrderBook(account string, gets string, pays string) (error, string) {
	if !remote.Status {
		err := remote.Connect()
		if err != nil {
			host_port := remote.Wsconn.Host + ":" + remote.Wsconn.Port
			Error("Connect ", host_port, "fail! errno = ", err)
			return err, ""
		}
	}
	request := Pack_RequestOrderBook(account, gets, pays)
	err := remote.send(request)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := remote.read()
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Tx succ: ", response)
	return nil, response
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
