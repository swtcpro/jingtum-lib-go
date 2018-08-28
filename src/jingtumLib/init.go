/***  初始化
*** init.go
*** 主要用于初始化运行前的准备工作，例如初始化日志，读取配置文件，初始化网络等
*** author:              1416205324@qq.com
*** last_modified_time:  2018-5-25 13:13:23
 */

package jingtumlib

import (
	"jingtumlib/constant"
)

var (
	JTConfig = new(Config)
	Seq      = 1
)

func InitConfig() error {
	return JTConfig.InitConfig("../../conf/jingtum-lib-config.txt")
}

func Init() error {
	err := InitConfig()
	if err != nil {
		return err
	}

	constant.CFGCurrency = JTConfig.Read("Config", "currency")

	return nil
}

//退出
func Exits() {
}
