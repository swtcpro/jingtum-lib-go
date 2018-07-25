package jingtumBaseLib

import (
	"testing"
)

func Test_walletKeyPair(t *testing.T) {
	wt, _ := FromSecret("snsYqv2FsYLuibE9TGHdG5x5V5Qcn")
	t.Log("Wt public key : ", wt.GetPublicKey().ToAddress())
}
