/**
 * 底层区块链网络通信服务类，不对外部提供方法。
 *
 * @FileName: server.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-31 14:44:32
 * @UpdateTime: 2018-07-31 14:44:54
 */
package jingtumLib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"jingtumLib/constant"
	"jingtumLib/utils"

	"github.com/caivega/evtwebsocket"
)

//Server 区块链网络通信服务结构体。
type Server struct {
	id        uint64
	remote    *Remote
	connected bool
	opened    bool
	state     string
	conn      *evtwebsocket.Conn
	opts      map[string]interface{}
	url       string
	reqs      chan *ReqCtx
	l         *sync.RWMutex
}

type activeStates []string

var (
	onlineStates = activeStates{"syncing", "tracking", "proposing", "validating", "full", "connected"}
	domainRE     = "[A-Za-z0-9]+(\\.[A-Za-z0-9]){1,5}" //"^(?=.{1,255}$)[0-9A-Za-z](?:(?:[0-9A-Za-z]|[-_]){0,61}[0-9A-Za-z])?(?:\\.[0-9A-Za-z](?:(?:[0-9A-Za-z]|[-_]){0,61}[0-9A-Za-z])?)*\\.?$"
)

//NewServer 创建区块链网络通信服务。
func NewServer(remote *Remote, urlStr string) (*Server, error) {
	if urlStr == "" || remote == nil {
		return nil, constant.ERR_EMPTY_PARAM
	}

	server := new(Server)
	server.opts = make(map[string]interface{})

	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	h := strings.Split(urlParsed.Host, ":")
	if len(h) == 1 {
		server.opts["host"] = urlParsed.Host
	} else if len(h) == 2 {
		server.opts["host"] = h[0]
	}

	if !utils.MatchString(domainRE, server.opts["host"].(string)) {
		return nil, constant.ERR_SERVER_HOST_INCORRECT
	}

	if !utils.MatchString("^\\d+$", urlParsed.Port()) {
		return nil, constant.ERR_SERVER_PORT_ERROR
	}

	iport, err := strconv.Atoi(urlParsed.Port())
	if err != nil {
		return nil, err
	}

	if iport < 1 || iport > 65535 {
		return nil, constant.ERR_SERVER_PORT_OUT_OF_RANGE
	}

	server.opts["port"] = iport
	server.opts["protocol"] = urlParsed.Scheme

	if urlParsed.Scheme == "wss" {
		server.opts["secure"] = true
	} else {
		server.opts["secure"] = false
	}

	if server.opts["secure"].(bool) {
		server.url = "wss://" + server.opts["host"].(string) + ":" + urlParsed.Port()
	} else {
		server.url = "ws://" + server.opts["host"].(string) + ":" + urlParsed.Port()
	}

	server.id = 0
	server.remote = remote
	server.state = "offline"
	server.connected = false
	server.opened = false
	server.l = new(sync.RWMutex)
	server.reqs = make(chan *ReqCtx)
	return server, nil
}

//Disconnect 关闭连接
func (server *Server) Disconnect() bool {
	if server == nil || !server.connected {
		return true
	}

	req := server.remote.UnSubscribe([]string{"transactions", "ledger", "server"})
	req.Submit(func(err error, result interface{}) {})

	rc := new(ReqCtx)
	rc.command = constant.CommandDisconnect
	server.sendMessage(rc)
	server.remote.emit.Off("*")

	err := server.conn.Close()
	if err != nil {
		return false
	}

	server.state = "offline"
	server.connected = false
	server.opened = false
	return true
}

//IsConnected true已连接。
func (server *Server) IsConnected() bool {
	return server.connected
}

//GetCid 每次请求序列递增。
func (server *Server) GetCid() uint64 {
	server.l.Lock()
	defer server.l.Unlock()
	server.id++
	return server.id
}

func (server *Server) sendMessage(reqCtx *ReqCtx) {
	server.reqs <- reqCtx
}

func (server *Server) setState(state string) {
	if state == server.state {
		return
	}
	server.state = state
	server.connected = (state == "online")

	if !server.connected {
		server.opened = false
	}
}

func (server *Server) listeningSend() {
	for {
		req := <-server.reqs

		//终止消息监听线程
		if req.command == constant.CommandDisconnect {
			break
		}

		req.data["id"] = req.cid
		req.data["command"] = req.command
		jsonData, err := json.Marshal(req.data)
		if err != nil {
			req.callback(err, nil)
			continue
		}

		fmt.Printf("Request info %s\n", jsonData)
		bm := evtwebsocket.Msg{
			Body: jsonData,
			Callback: func(msg []byte, w *evtwebsocket.Conn) {
				// fmt.Printf("Response message : %s\n", msg)
			},
		}

		// 发送消息
		if err := server.conn.Send(bm); err != nil {
			server.Disconnect()
			req.callback(err, nil)
			break
		}
	}
}

func (server *Server) connect(callback func(err error, result interface{})) error {
	if server.connected {
		return nil
	}

	if server.conn != nil {
		server.Disconnect()
	}

	var once sync.Once
	wg := &sync.WaitGroup{}
	wg.Add(1)

	server.conn = &evtwebsocket.Conn{

		OnConnected: func(w *evtwebsocket.Conn) {
			server.connected = true
			server.opened = true
			server.state = "online"
			connectMsg := fmt.Sprintf("Connect to [%s] success.", server.url)
			once.Do(wg.Done)
			callback(nil, connectMsg)
			// go func() {
			// 	req := server.remote.Subscribe([]string{"transactions", "ledger", "server"})
			// 	req.Submit(func(err error, result interface{}) {
			// 	})
			// }()
		},

		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			fmt.Printf("On message %s\n", msg)
			server.remote.handleMessage(msg)
		},

		MatchMsg: func(req, resp []byte) bool {
			return true
		},

		OnError: func(err error) {
			Errorf("On error : %s", err.Error())
			//自动重连
			server.Disconnect()
			server.connect(func(err error, result interface{}) {
				if err != nil {
					Errorf("ReConnect fail. error : %v", err)
					server.Disconnect()
				}
			})
		},

		Reconnect: true,
	}
	err := server.conn.Dial(server.url, "")
	if err != nil {
		callback(err, nil)
		return err
	}

	go server.listeningSend()

	wg.Wait()

	return nil
}

func (status activeStates) contain(value string) bool {
	return status.indexOf(value) >= 0
}

func (status activeStates) indexOf(value string) int {
	for i, inv := range status {
		if inv == value {
			return i
		}
	}

	return -1
}
