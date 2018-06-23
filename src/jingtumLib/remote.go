package jingtumLib

import (
	"common/golang.org/x/net/websocket"
	"errors"
	"time"
)

/*
type jtremote interface {
    //连接
    connect()
    //断开
    disconnect()
    //请求底层服务器信息
    requestServerInfo()
    //获取最新账本信息
    requestLedgerClosed()
    //获取某一账本具体信息
    requestLedger(options)
    //查询某一交易具体信息
    requestTx(options)
    //请求账号信息
    requestAccountInfo(options)
    //获得账号可接收和发送的货币
    requestAccountTums(options)
    //获得账号关系
    requestAccountRelations(options)
    //获得账号挂单
    requestAccountOffers(options)
    //获得账号交易列表
    requestAccountTx(options)
    //获得市场挂单列表
    requestOrderBook(options)
}
*/

type Conn struct {
	Ws   *websocket.Conn
	Host string
	Port string
}

/*
* 连接函数
 */
func Connect() (error, *Conn) {
	var host_port string
	var conn Conn
	conn.Host = ""
	conn.Port = ""
	host := JTConfig.Read("Service", "Host")
	if host == "" {
		Error("Config Service:Host is null.")
		return errors.New("Config|service:Host setting error"), &conn
	}
	port := JTConfig.Read("Service", "Port")
	if port == "" {
		Error("Config Service:Port is null.")
		return errors.New("Config|service:Port setting error"), &conn
	}
	host_port = host + ":" + port
	origin := "http://localhost/"
	ws, err := websocket.Dial(host_port, "", origin)
	if err != nil {
		Error("Connect ", host_port, "fail! errno = ", err)
	} else {
		Info("Connect ", host_port, " succ. ")
	}
	conn.Ws = ws
	conn.Host = host
	conn.Port = port
	return err, &conn
}

func Get_now_time() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

/*
* 发送
 */
func send(request string, conn *Conn) error {
	_, err := conn.Ws.Write([]byte(request))
	if err != nil {
		Error("Send ", request, "fail!", "errno:", err)
	} else {
		Info("Send: ", request, "succ.")
	}
	return err
}

/*
* 接收
 */

func read(conn *Conn) (error, string) {
	var msg = make([]byte, 409600)
	var n int
	n, err := conn.Ws.Read(msg)
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
func Disconnect(conn *Conn) {
	conn.Ws.Close()
}

/*
* 请求底层服务器信息
 */
func RequestServerInfo(conn *Conn) (error, string) {
	request := Pack_RequestServerInfo()
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Server Info succ: ", response)
	return nil, response
}

/*
* 获取最新账本信息
 */
func RequestLedgerClosed(conn *Conn) (error, string) {
	request := Pack_RequestLedgerClosed()
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
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
func RequestLedger(conn *Conn, ledger_index string, ledger_hash string, transactions bool) (error, string) {
	if ledger_index == "" && ledger_hash == "" {
		return errors.New("ledger_index|ledger_hash value error"), ""
	}
	request := Pack_RequestLedger(ledger_index, ledger_hash, transactions)
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
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
func RequestTx(conn *Conn, hash string) (error, string) {
	request := Pack_RequestTx(hash)
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
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
func RequestAccountInfo(conn *Conn, account string) (error, string) {
	request := Pack_RequestAccountInfo(account)
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
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
func RequestAccountTums(conn *Conn, account string) (error, string) {
	request := Pack_RequestAccountTums(account)
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
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
func RequestAccountRelations(conn *Conn, account string, atype string) (error, string) {
	request := Pack_RequestAccountRelations(account, atype)
	if request == "" {
		return errors.New("RequestAccountRelations type value error"), ""
	}
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Relations succ: ", response)
	return nil, response
}

/*
* 获得账号交易列表
 */
func RequestAccountTx(conn *Conn, account string, limit int) (error, string) {
	request := Pack_RequestAccountTx(account, limit)
	err := send(request, conn)
	if err != nil {
		Error("Send data fail!")
		return err, ""
	}
	err, response := read(conn)
	if err != nil {
		Error("Received data fail!")
		return err, ""
	}
	Info("Get Reqonse Account Tx succ: ", response)
	return nil, response
}
