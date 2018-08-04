package jingtumLib

import (
	"container/list"
	"fmt"
	"jingtumLib/serializer"
	"os"
	"testing"
)

func Test_AddMemo(t *testing.T) {
	remote, err := NewRemote("",false)
	if err != nil {
		t.Fatalf("New remote error : %v.", err)
	}

	remote.LocalSign = true

	t.Logf("Success remote : %v.", remote)
	tx, err := NewTransaction(remote)

	if err != nil {
		t.Fatalf("New transactino error : %v.", err)
	}

	tx.AddTxJson("TransactionType", "Payment")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("我的测试")
	tx.AddMemo("支付0.000001SWT")
	tx.AddMemo("支付0.000001SWT")
	tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	t.Logf("Get tx TransactionType : %s. Flags : %d. Fee : %d. Secret : %s", tx.GetTxJson("TransactionType"), tx.GetTxJson("Flags"), tx.GetTxJson("Fee"), tx.GetTxJson("secret"))

	memos := tx.GetTxJson("Memos").(*list.List)

	for e := memos.Front(); e != nil; e = e.Next() {
		t.Logf("Get tx memos info %s.", e.Value.(*serializer.MemoInfo).Memo.MemoData)
	}

	//	for e := memos.Front(); e != nil; e = memos.Next() {
	//		t.Logf("Get tx memos info %s.", e.Memo.MemoData)
	//	}

	//	for i, v := range memos {
	//		t.Logf("Get tx memos info %d = %s.", i, v.Memo.MemoData)
	//	}
}

func TestMain(m *testing.M) {
	err := Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	ret := m.Run()
	os.Exit(ret)
}
