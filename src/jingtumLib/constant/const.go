/**
 * 常量类。
 *
 * @FileName: errorCode.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package constant

//AccountPrefix AccountPrefix
const AccountPrefix uint8 = 0

//SeedPrefix SeedPrefix
const SeedPrefix uint8 = 33

//RegexCurrency RegexCurrency
const RegexCurrency = "^([a-zA-Z0-9]{3,6}|[A-F0-9]{40})$"

//TxJSONErrorKey TxJSONErrorKey
const TxJSONErrorKey = "Error"

//CommandDisconnect 服务关闭命令，用于底层链接断开后终止消息监听线程
const CommandDisconnect = "server_disconnect"

//CommandServerInfo CommandServerInfo
const CommandServerInfo = "server_info"

//CommandSubmit 提交命令
const CommandSubmit = "submit"

//CommandLedgerClosed 获取最新账本命令
const CommandLedgerClosed = "ledger_closed"

//CommandLedger 获取某一账本命令
const CommandLedger = "ledger"

//CommandTX 查询某一交易信息命令
const CommandTX = "tx"

//CommandAccountInfo 获取账号信息
const CommandAccountInfo = "account_info"

//CommandAccountCurrencies 获得账号可接收和发送的货币
const CommandAccountCurrencies = "account_currencies"

//CommandAccountLines CommandAccountLines
const CommandAccountLines = "account_lines"

//CommandAccountRelation CommandAccountRelation
const CommandAccountRelation = "account_relation"

//CommandAccountOffers 获得账号挂单
const CommandAccountOffers = "account_offers"

//CommandAccountTX 获得账号交易列表
const CommandAccountTX = "account_tx"

//CommandBookOffers 获得市场挂单列表
const CommandBookOffers = "book_offers"

//CommandSubscribe 订阅事件
const CommandSubscribe = "subscribe"

//CommandUnSubscribe 退订事件
const CommandUnSubscribe = "unsubscribe"

//AccountOne AccountOne
const AccountOne = "jjjjjjjjjjjjjjjjjjjjBZbvri"

//EventTX 交易事件
const EventTX = "transactions"

//EventLedgerClosed 账本事件
const EventLedgerClosed = "ledger_closed"
