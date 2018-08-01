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

	//"golang.org/x/net/websocket"
	"github.com/caivega/evtwebsocket"
)

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

var (
	onlineStates = []string{"syncing", "tracking", "proposing", "validating", "full", "connected"}
	domainRE     = "^(?=.{1,255}$)[0-9A-Za-z](?:(?:[0-9A-Za-z]|[-_]){0,61}[0-9A-Za-z])?(?:\\.[0-9A-Za-z](?:(?:[0-9A-Za-z]|[-_]){0,61}[0-9A-Za-z])?)*\\.?$"
)

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

func (server *Server) Disconnect() bool {
	err := server.conn.Close()
	if err != nil {
		return false
	}

	server.state = "offline"
	server.connected = false
	server.opened = false
	return true
}

func (server *Server) IsConnected() bool {
	return server.connected
}

func (server *Server) GetCid() uint64 {
	server.l.Lock()
	defer server.l.Unlock()
	server.id++
	return server.id
}

func (server *Server) sendMessage(reqCtx *ReqCtx) {
	server.reqs <- reqCtx
}

func (server *Server) listeningSend() {
	for {
		req := <-server.reqs

		m, err := json.Marshal(req.data)
		if err != nil {
			req.callback(err, nil)
			continue
		}

		cmd := fmt.Sprintf("{\"id\":\"%d\",\"command\":\"%s\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1, \"limit\": %d, \"marker\":%s}", req.cid, req.command, req.data["account"].(string), req.data["limit"].(float64), m)

		bm := evtwebsocket.Msg{
			Body: []byte(cmd),
			Callback: func(msg []byte, w *evtwebsocket.Conn) {
				Infof("Response message : %s\n", msg)
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

func (server *Server) connect(callback func(err error, result interface{})) {
	if server.connected {
		return
	}

	if server.conn != nil {
		server.Disconnect()
	}

	server.conn = &evtwebsocket.Conn{

		OnConnected: func(w *evtwebsocket.Conn) {
			server.connected = true
			server.opened = true
		},

		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			server.remote.handleMessage(msg)
		},

		MatchMsg: func(req, resp []byte) bool {
			return true
		},

		OnError: func(err error) {
			callback(err, nil)
		},

		Reconnect: true,
	}
	err := server.conn.Dial(server.url, "")
	if err != nil {
		callback(err, nil)
		return
	}

	go server.listeningSend()
}
