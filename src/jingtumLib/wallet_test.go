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
}
