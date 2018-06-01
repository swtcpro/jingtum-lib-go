/***  初始化
 *** init.go
 *** 主要用于初始化运行前的准备工作，例如初始化日志，读取配置文件，初始化网络等
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package jingtumLib

import (
	Log "common/github.com/blog4go"
	"fmt"
	"os"
)

func InitLog() {
	err := Log.NewFileWriter("../../log/", true)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer Log.Close()
}

func Init() {
	InitLog()
    Log.Debugf("Good morning, %s", "eddie")
}
