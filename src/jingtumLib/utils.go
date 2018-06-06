/**
 *
 * 文件功能介绍
 *
 * @FileName: utils.go
 * @Auther : 13851485286
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumLib

import (
     "strconv"
     "regexp"
)

var (

	LEDGER_STATES = map[string]string{"current": "current", "closed": "closed", "validated": "validated"}//[3]string{"current", "closed", "validated"}
)

func Number(s string) bool {

    _, err := strconv.Atoi(s)

    return err != nil
}

func MatchString(patter string, str string) bool {
    match, _ := regexp.MatchString(patter, str)

    return match
}
