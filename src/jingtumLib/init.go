package jingtumLib

import (
	Log "common/github.com/blog4go"
	"fmt"
	"os"
)

func InitLog() {
	err := Log.NewFileWriter("../log/", true)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer Log.Close()
}

func Init() {
	InitLog()
}
