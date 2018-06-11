/***  初始化
*** init.go
*** 主要用于初始化运行前的准备工作，例如初始化日志，读取配置文件，初始化网络等
*** author:              1416205324@qq.com
*** last_modified_time:  2018-5-25 13:13:23
 */

package jingtumLib

import (
	log "common/github.com/blog4go"
	"fmt"
	"sync"
)

var (
	JTConfig = new(Config)
)

type MyHook struct {
	cnt     int
	level   log.LevelType
	message string

	l *sync.RWMutex
}

func NewMyHook() (hook *MyHook) {
	hook = new(MyHook)
	hook.cnt = 0
	hook.level = log.TRACE
	hook.message = ""
	hook.l = new(sync.RWMutex)

	return
}

func (hook *MyHook) Add() {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.cnt++
}

func (hook *MyHook) Cnt() int {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.cnt
}
func (hook *MyHook) Level() log.LevelType {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.level
}

func (hook *MyHook) SetLevel(level log.LevelType) {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.level = level
}

func (hook *MyHook) Message() string {
	hook.l.RLock()
	defer hook.l.RUnlock()
	return hook.message
}

func (hook *MyHook) SetMessage(message string) {
	hook.l.Lock()
	defer hook.l.Unlock()
	hook.message = message
}

func (hook *MyHook) Fire(level log.LevelType, tags map[string]string, args ...interface{}) {
	hook = NewMyHook()
	hook.Add()
	hook.SetLevel(level)
	hook.SetMessage(fmt.Sprint(args...))
}

func InitLog() (err error) {
	rerr := log.NewWriterFromConfigAsFile("../conf/jingtum-lib.xml")
	//rerr := log.NewWriterFromConfigAsFile("../conf/config.example.xml")
	if nil != rerr {
		fmt.Println(rerr.Error())
		return rerr
	}
	// initialize your hook instance
	hook := new(MyHook)
	log.SetHook(hook) // writersFromConfig can be replaced with writers
	log.SetHookLevel(log.INFO)
	log.SetHookAsync(true) // hook will be called in async mode
	return
}

// Debug static function for Debug
func Debug(args ...interface{}) {
	log.Debug(args...)
	log.Flush()
}

//Debugf static function for Debugf
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
	log.Flush()
}

// Info static function for Info
func Info(args ...interface{}) {
	log.Info(args...)
	log.Flush()
}

// Infof static function for Infof
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
	log.Flush()
}

// Warn static function for Warn
func Warn(args ...interface{}) {
	log.Warn(args...)
	log.Flush()
}

// Warnf static function for Warnf
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
	log.Flush()
}

// Error static function for Error
func Error(args ...interface{}) {
	log.Error(args...)
	log.Flush()
}

// Errorf static function for Errorf
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	log.Flush()
}

// Critical static function for Critical
func Critical(args ...interface{}) {
	log.Critical(args...)
	log.Flush()
}

// Criticalf static function for Criticalf
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
	log.Flush()
}

func Flush() {
	log.Flush()
}

func InitConfig() {
	JTConfig.InitConfig("conf/jing_tong_lib_config.txt")
}

func Init() (err error) {
	err := InitLog()
	if err != nil {
		return err
	}
	err = InitConfig()
	if err != nil {
		return err
	}
	return
}

//退出
func Exits() {
	defer log.Close()
}
