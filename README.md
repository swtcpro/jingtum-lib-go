# jingtum-lib-go
用于和井通区块链网络进行交互的jingtum-lib-go库。提供简单易用的go语言版的本地调用库。

## Source code
* /src/jingtumLib - 源码文件
* /src/testLib - 提供所有接口的一个集成测试包
* docs - jingtum-lib-go 使用文档

## 开发环境
* Windows 10
* go version go1.10.1 windows/amd64

## 概述
基于ws协议的jingtuml-lib-go库来连接井通区块链网络系统。提供了公共的api来创建两种对象：GET请求的Request对象和POST请求的Transaction对象，然后通过Submi方法提交。

## 使用示例
1) 创建 Remote实例。
```
    remote, err := NewRemote("wss://s.jingtum.com:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}
```
2) 链接服务
```
cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})
```
3) 通过Remote创建Request或Transaction 实例，调用Submit方法提交请求。下面是请求服务器信息示例： 
```
   req, err := remote.RequestServerInfo()
	if err != nil {
		t.Fatalf("Fail request server info %s", err.Error())
	}
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request server info %s", err.Error())
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request server info %s", jsonByte)
	})
```

4) 关闭链接。
```
    remote.Disconnect()
```

## Remote
Remote是跟井通底层交互最主要的类，它可以组装交易发送到底层、订阅事件及从底层拉取数据。

## Request
Request类主管GET请求，包括获得服务器、账号、挂单、路径等信息。请求时不需要提供密钥，且对所有用户公开。所有的请求是异步的，会提供一个回调函数。可以从回调函数获取错误或者结果。

## Transaction
Transaction主管POST请求，包括组装交易和交易参数。请求时需要提供密钥，且交易可以进行本地签名和服务器签名。目前支持服务器签名，本地签名支持主要的交易，还有部分参数不支持。所有的请求是异步的，会提供一个回调函数。每个回调函数包含错误信息和成功后的结果。

## Events
您可以监听服务器的事件。
* 监听系统发生的所有交易事件。
* 监听最新的账号事件。
* 监听所有服务器状态更改事件。