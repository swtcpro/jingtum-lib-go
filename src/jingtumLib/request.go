/**
 * Request类主管GET请求，包括获得服务器、账号、挂单、路径等信息。请求时不需要提供密
 * 钥，且对所有用户公开。所有的请求是异步的，会提供一个回调函数。每个回调函数有两个参
 * 数，一个是错误，另一个是结果。
 *
 * @FileName: request.go
 * @Auther : 13851485286
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-05-28 10:44:32
 * @UpdateTime: 2018-05-28 10:44:54
 */

package jingtumLib

import (
    _ "errors"

    log "common/github.com/blog4go"
)

type FilterFunc func(interface{}) interface{}

type Request struct {
	remote        Remote
	message       map[string]interface{}
    command       string
    filter        FilterFunc
}

func NewRequest(remote *Remote, command string,filter FilterFunc) (request *Request , err error) {

    request = new(Request)
    request.remote = remote
    request.command = command
    request.filter = filter
    request.message = mak(map[string]interface{})
	return &request, nil
}

//提交请求
func (request *Request) Submit(callback func(data interface{})) {
    
    for _, v := range request.message {
        //if v, ok := s.(string); ok {
		//fmt.Println(v)
        //if v, ok := interface{}(s).(string); ok {
		//fmt.Println(v)//}//}

        if _v, ok := v.(error); ok {
           callback(v)
           return
        }
    }

    request.remote.Submit(request.command, request.message, request.filter, callback)

}

func (request *Request) SelectLedger(ledger interface{}) {
    
    switch ledger.(type) {

    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
        request.message["ledger_index"] = ledger
    case string:  
        _, ok := LEDGER_STATES [ledger]

        if (ok) {
            request.message["ledger_index"] = ledger
        } else if MatchString("^[A-F0-9]+$",ledger) {
            request.message["ledger_hash"] = ledger;
        } else {
            request.message["ledger_index"] = "validated"
        }
    default:
        request.message["ledger_index"] = "validated"
    } 
}