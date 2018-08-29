/***  读取配置文件
*** config.go
*** 主要用于读取配置文件数据
*** author:              1416205324@qq.com
*** last_modified_time:  2018-6-6 23:13:23
 */

package jingtumlib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const middle = "========="

//Config 配置类
type Config struct {
	Mymap  map[string]string
	strcet string
}

//InitConfig 初始化配置
func (c *Config) InitConfig(path string) error {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err.Error())
			return err
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
	fmt.Printf("Init config succ.")
	return nil
}

func (c Config) Read(node, key string) string {
	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}

//ReadInt 读取int类型配置，可以使用默认值。
func (c Config) ReadInt(node, key string, defaultv int) int {
	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return defaultv
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic("Invalid string convert to int.")
	}
	return n
}
