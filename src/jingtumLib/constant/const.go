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

const CommandLedgerClosed = "ledger_closed"
