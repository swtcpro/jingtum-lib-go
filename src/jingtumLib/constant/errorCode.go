/**
 * 错误定义类。
 *
 * @FileName: errorCode.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */

package constant

import (
	"errors"
)

var (

	//通用错误码
	ERR_EMPTY_PARAM = errors.New("Parameters cannot be empty.")

	ERR_INVALID_PARAM = errors.New("Invalid input")
)
