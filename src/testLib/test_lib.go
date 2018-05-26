/***  初始化
 *** testLib.go
 *** 主要用于用于测试jingtumLib的各个实例
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package main

import (
	"fmt"
	jingtum "jingtumLib"
)

func main() {
	jingtum.Init()
	jingtum.Generate()
	fmt.Println("test_lib")
}
