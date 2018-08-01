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

	ERR_INVALID_PARAM = errors.New("Invalid input.")

	//底层通信类相关错误码
	ERR_SERVER_HOST_INCORRECT = errors.New("server host incorrect.")

	ERR_SERVER_PORT_ERROR = errors.New("server port not a number.")

	ERR_SERVER_PORT_OUT_OF_RANGE = errors.New("server port out of range.")

	//支付相关错误码
	ERR_PAYMENT_INVALID_SRC_ADDR = errors.New("invalid source address.")

	ERR_PAYMENT_INVALID_DST_ADDR = errors.New("invalid destination address.")

	ERR_PAYMENT_INVALID_AMOUNT = errors.New("invalid amount.")

	ERR_PAYMENT_OUT_OF_AMOUNT = errors.New("invalid amount: amount's maximum value is 100000000000.")

	ERR_PAYMENT_MEMO_EMPTY = errors.New("Memo is empty.")

	ERR_PAYMENT_OUT_OF_MEMO_LEN = errors.New("The length of Memo shoule be less than or equal 2048.")

	ERR_PAYMENT_INVALID_SECRET = errors.New("invalid secret.")
)
