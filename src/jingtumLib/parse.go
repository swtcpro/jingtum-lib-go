package jingtumLib

import (
	"fmt"
)

const INT_MAX = int(^uint(0) >> 1)

func Get_Seq() {
	if Seq == INT_MAX {
		Seq = 1
	} else {
		Seq = Seq + 1
	}
}

func Pack_RequestServerInfo() string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"server_info\"}", Seq)
}

func Pack_RequestLedgerClosed() string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"ledger_closed\"}", Seq)
}

func Pack_RequestLedger(ledger_index string, ledger_hash string, transactions bool) string {
	Get_Seq()
	request := fmt.Sprintf("{\"id\":\"%d\",\"command\":\"ledger\",\"ledger_index\":\"%s\",\"transactions\":\"%t\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, ledger_index, transactions)
	if ledger_index == "" {
		request = fmt.Sprintf("{\"id\":\"%d\",\"command\":\"ledger\",\"ledger_hash\":\"%s\",\"transactions\":\"%t\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, ledger_hash, transactions)
	}
	return request
}

func Pack_RequestTx(hash string) string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"tx\",\"transaction\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, hash)
}

func Pack_RequestAccountInfo(account string) string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"account_tx\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, account)
}

func Pack_RequestAccountTums(account string) string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"account_currencies\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, account)
}

func Pack_RequestAccountTx(account string, limit int) string {
	Get_Seq()
	return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"account_tx\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1, \"limit\":\"%d\"}", Seq, account, limit)
}

func Pack_RequestAccountRelations(account string, atype string) string {
	Get_Seq()
	if atype == "trust" {
		return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"account_lines\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, account)
	} else if atype == "authorize" || atype == "freeze" {
		return fmt.Sprintf("{\"id\":\"%d\",\"command\":\"account_relation\",\"account\":\"%s\",\"ledger_index_min\":-1,\"ledger_index_max\":-1}", Seq, account)
	} else {
		return ""
	}
}
