/**
 * 钱包测试类
 *
 * @FileName: wallet_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumLib

import (
	"testing"
)

/**
 * 钱包创建测试用例
 */
func Test_Wallet(t *testing.T) {
	secret := "snsYqv2FsYLuibE9TGHdG5x5V5Qcn"

	//私钥合法性测试
	isOk := IsValidSecret(secret)

	if !isOk {
		t.Fatalf("Failure IsValidSecret(%s) is false", secret)
	}

	t.Logf("Success IsValidSecret(%s) is true", secret)

	//根据私钥创建测试
	wt, err := FromSecret(secret)

	if err != nil {
		t.Fatalf("Failure FromSecret : %s, err %v", secret, err)
	}

	t.Logf("Success FromSecret(%s). PublicKey : %s. Wallet address : %s", wt.GetSecret(), wt.GetPublicKey(), wt.GetAddress())

	//钱包地址合法性验证

	isOk = IsValidAddress(wt.GetAddress())

	if !isOk {
		t.Fatalf("Failure IsValidAddress(%s) is false", wt.GetAddress())
	}

	t.Logf("Success IsValidAddress(%s) is true", wt.GetAddress())

	//生成新钱包
	newWallet, err := Generate()
	isOk = IsValidSecret(newWallet.GetSecret())
	if !isOk {
		t.Fatalf("New secret IsValidSecret(%s) is false", newWallet.GetSecret())
	}

	isOk = IsValidAddress(newWallet.GetAddress())
	if !isOk {
		t.Fatalf("New address IsValidAddress(%s) is false", newWallet.GetAddress())
	}

	t.Logf("Success new secret (%s). address (%s)", newWallet.GetSecret(), newWallet.GetAddress())
}
