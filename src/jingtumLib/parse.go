package jingtumlib

import (
	"encoding/json"
	"strconv"
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
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "server_info"
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestLedgerClosed() string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "ledger_closed"
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestLedger(ledger_index string, ledger_hash string, transactions bool) string {
	Get_Seq()
	if ledger_index == "" && ledger_hash == "" {
		return ""
	}
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "ledger"
	mrequest["transactions"] = strconv.FormatBool(transactions)
	if ledger_index == "" {
		mrequest["ledger_hash"] = ledger_hash
	} else {
		mrequest["ledger_index"] = ledger_index
	}
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestTx(hash string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "tx"
	mrequest["transaction"] = hash
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestAccountInfo(account string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "account_tx"
	mrequest["account"] = account
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestAccountTums(account string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "account_currencies"
	mrequest["account"] = account
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestAccountRelations(account string, atype string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	switch atype {
	case "trust":
		mrequest["command"] = "account_lines"
	case "authorize":
		mrequest["command"] = "account_relation"
	case "freeze":
		mrequest["command"] = "account_relation"
	default:
		return ""
	}
	mrequest["account"] = account
	mrequest["relation_type"] = atype
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestAccountOffers(account string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "account_offers"
	mrequest["account"] = account
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestAccountTx(account string, limit int) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "account_tx"
	mrequest["account"] = account
	mrequest["limit"] = strconv.Itoa(limit)
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func Pack_RequestOrderBook(account string, gets string, pays string) string {
	Get_Seq()
	mrequest := make(map[string]string)
	mrequest["id"] = strconv.Itoa(Seq)
	mrequest["command"] = "book_offers"
	mrequest["account"] = account
	mrequest["taker_gets"] = gets
	mrequest["taker_pays"] = pays
	mjson, err := json.Marshal(mrequest)
	if err != nil {
		return ""
	}
	return string(mjson)
}
