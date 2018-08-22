# Design Spec

## Features
It has the same funtionalities as jingtum-lib-nodejs. 

https://github.com/swtcpro/jingtum-lib-nodejs

## References
* WebSocket (https://github.com/caivega/evtwebsocket) The jingtum-lib-go library is based on the ws protocol to connect with jingtum system. 
* Portable.BouncyCastle (http://github.com/btcsuite/btcd/btcec) The jingtum-lib-go library local sign depends on ECDSA signature.

## Models
* The inner Server class performs the websocket communication.
* The Remote class provides public APIs to create two kinds of objects: Request object, and Transaction object.
* The Request class is used to request info.
* The Transaction class is used to operate transactions. 
* Request class and Transacton class both use Submit(callback) method to submit data to server.
* The result can be handled by the callback.

```
|-----------|     |--------------|     |--------|     |--> [ Request Object ]
| WebSocket | --> | Server       | --> | Remote | --> |    
| Protocal  | <-- | Inner Class  | <-- | Class  |     |--> [ Transaction Object]
|-----------|     |--------------|     |--------|
```

## Stubs
* Account stub listen all the transactions in server, and then filter them for specfic account.
* OrderBook stub listen all the transactions in server, and then filter them for specfic gets/pays pair.

## Data
* The json string is sent to server for request operation.
* The json string is sent to server for transaction operation (server sign).
* The transaction data is serialized to blob string and then sent to server for transaction operation (local sign).
* The json string is reveived from server for reqeust and transaction operations.
* The callback result contains:
  * The raw message from server, in json format.
  * The exception message if the operation is refused by the server.
  * The result object if the operation is succeed. It is parsed from the json message.

## Local Sign
The local sign is implemented by serializing the json string into binary blob, and then send the blob string to server. 

The inner Serializer class performs the serialization. The data members are grouped into different categories, and then serialized as data type and data value pair.

The category contains:
* Int8
* Int16
* Int32
* Int64
* Hash128
* Hash160
* hash256
* Amount
* VL (string)
* Account (string)
* PathSet
* Object
* Array

# Dcuments
Usage for jingtum-lib-go. All classes are under the namespace JingTum.Lib. 

## Wallet struct
### Genreate()
Genereates a new wallet.

#### sample
```
newWallet, err := jingtumLib.Generate()
```

### FromSecret(secret)
Creates a wallet from existing secret. The secret is the private secret of jingtum wallet.

#### sample
```
wt, err := jingtumLib.FromSecret(secret)
```

## Remote class
Main function class in jingtum-lib-csharp. It creates a handle with jingtum, makes request to jingtum, subscribs event to jingtum, and gets info from jingtum.

* Connect(callback func(err error, result interface{})) error
* GetNowTime() string
* Disconnect()
* RequestServerInfo() (*Request, error)
* RequestLedgerClosed() (*Request, error)
* RequestLedger(options map[string]interface{}) (*Request, error)
* RequestTx(hash string) (*Request, error)
* RequestAccountInfo(options map[string]interface{}) (*Request, error)
* RequestAccountTums(options map[string]interface{}) (*Request, error)
* RequestAccountRelations(options map[string]interface{}) (*Request, error)
* RequestAccountOffers(options map[string]interface{}) (*Request, error)
* RequestAccountTx(options map[string]interface{}) (*Request, error)
* RequestOrderBook(options map[string]interface{}) (*Request, error)
* BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error)
* BuildRelationSet(options map[string]interface{}, tx *Transaction) error
* BuildTrustSet(options map[string]interface{}, tx *Transaction) error
* BuildRelationTx(options map[string]interface{}) (*Transaction, error)
* BuildAccountSet(options map[string]interface{}, tx *Transaction) error
* BuildDelegateKeySet(options map[string]interface{}, tx *Transaction) error
* BuildSignerSet(options map[string]interface{}, tx *Transaction) error
* BuildAccountSetTx(options map[string]interface{}) (*Transaction, error)
* BuildOfferCreateTx(options map[string]interface{}) (*Transaction, error)
* BuildOfferCancelTx(options map[string]interface{}) (*Transaction, error)
* DeployContractTx(options map[string]interface{}) (*Transaction, error)
* CallContractTx(options map[string]interface{}) (*Transaction, error)

### NewRemote(url, localSign)
#### options
* url: The jingtum websocket server url.
* localSign: Whether sign transaction in local.

#### sample
```
remote, err := NewRemote(wsurl, true)
```

### Connect(callback)
Each remote object should connect jingtum first. Now jingtum should connect manual, only then you can send request to backend.


#### sample
```
remote.Connect(func(err error, result interface{}) {
    	if err != nil {
			return
		}
		jsonByte, _ := json.Marshal(result)
		t.Logf("Connect to %s success. Result : %s.", wsurl, jsonByte)
})
```

### Disconnect()
Remote object can be disconnected manual, and no parameters are required.

#### sample
```
remote.Disconnect()
```

### RequestServerInfo()
Create request object and get server info from jingtum.

#### sample
```
req, _ := remote.RequestServerInfo()
req.Submit(func(err error, result interface{}) {
    t.Logf("Success request server info %s", result)
})
```

### RequestLedgerClosed()
Create request object and get last closed ledger in system.

#### sample
```
req, _ := remote.RequestLedgerClosed()
req.Submit(func(err error, result interface{}) {
    jsonByte, _ := json.Marshal(result)
	t.Logf("Success request ledger closed %s", jsonByte)
})
```

### RequestLedger(options)
Create request object and get ledger in system.

#### options
(If none is provided, then last closed ledger is returned.)
* ledger_index: The ledger index.
* ledger_hash: The ledger hash.
* transactions: Whether include the transactions list in ledger.

#### sample
```
options := map[string]interface{}{"transactions": true, "ledger_index": 969054, "ledger_hash": "AEE4B16B543D8C8924F09C1DB822C6419780B86019F5F5FF8DC2938E7E0E89D2"}
req, _ := remote.RequestLedger(options)
req.Submit(func(err error, result interface{}) {
	jsonByte, _ := json.Marshal(result)
	t.Logf("Success request ledger %s", jsonByte)

})
```

### RequestTx(options)
Query one transaction information.

#### options
* Hash: The transaction hash.

#### sample
```
req, _ := remote.RequestTx("6537F72CE1DBD8043230C3FF64C6E5E95B11F6573D91EF6A13FEADE6940CB71A")
req.Submit(func(err error, result interface{}) {
	jsonByte, _ := json.Marshal(result)
	t.Log("Success request tx")
})
```

### RequestAccountInfo(options)
Get account info.

#### options
* account: The wallet address.
* ledger: (optional) 

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
req, _ := remote.RequestAccountInfo(options)
req.Submit(func(err error, result interface{}) {
    jsonBytes, _ := json.Marshal(result)
	t.Logf("Success Requst account info result : %s", jsonBytes)
})
```

### RequestAccountTums(options)
Each account holds many jingtum tums, and the received and sent tums can be found by RequestAccountTums.

#### options
* account: The wallet address.
* ledger: (optional)

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
req, _ := remote.RequestAccountTums(options)
req.Submit(func(err error, result interface{}) {
	jsonByte, _ := json.Marshal(result)
	t.Logf("Success request Account Tums %s", jsonByte)
})
```

### RequestAccountRelations(options)
Jingtum wallet is connected by many relations. Now jingtum supports `trust`, `authorize` and `freeze` relation, all can be queried by requestAccountRelations.

#### options
* account: The wallet addres.
* type: Trust, Ahthorize, Freeze
* ledger: (optional)
* limit: (optional) Limit the return relations count.
* marker: (optional) Request from the marker position. It can be got from the response of previous request.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust"}
req, _ := remote.RequestAccountRelations(options)
req.Submit(func(err error, result interface{}) {
	jsonByte, _ := json.Marshal(result)
	t.Logf("Success request account relations %s", jsonByte)
	})
```

### RequestAccountOffers(options)
Query account's current offer that is suspended on jingtum system, and will be filled by other accounts.

#### options
* account: The wallet address.
* ledger: (optional)
* limit: (optional) Limit the return offers count.
* marker: (optional) Request from the marker position. It can be got from the response of previous request.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
req, _ := remote.RequestAccountOffers(options)
req.Submit(func(err error, result interface{}) {
	jsonByte, _ := json.Marshal(result)
	t.Logf("Success request account offers %s", jsonByte)
	})
```

### RequestAccountTx(options)
Query account transactions.

#### options
* account: The wallet address.
* ledger: (optional) 
* limit: (optional) Limit the return trancations count.
* marker: (optional) Request from the marker position. It can be got from the response of previous request.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
req, _ := remote.RequestAccountTx(options)
req.Submit(func(err error, result interface{}) {
    t.Logf("Success request account tx : %s",result)
})
```

### RequestOrderBook(options)
Query order book info.

Firstly, each order book has a currency pair, as AAA/BBB. When to query the bid orders, gets is AAA and pays is BBB. When to query the ask orders, gets is BBB and pays is AAA.
The result is array of orders.

#### options
* gets: Amount object. (ignore the Value)
* pays: Amount object. (ignore the Value)

#### sample
```
options := make(map[string]interface{})
gets := jingtumLib.Amount{}
gets.Currency = "SWT"
pays := jingtumLib.Amount{}
pays.Currency = "CNY"
pays.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
options["gets"] = gets
options["pays"] = pays
req, _ := remote.RequestOrderBook(options)
req.Submit(func(err error, result interface{}) {
	jsonBytes, _ := json.Marshal(result)
	t.Logf("Success request order book : %s",jsonBytes)
	})
```

### RequestPathFind(options)
Query path from one curreny to another.

#### options
* account: The payment source address.
* destination: The payment target address.
* amount: The payment amount.

#### sample
```
options := make(map[string]interface{})
amount := jingtumLib.Amount{}
amonnt.Currency = "CNT"
amount.Issuer="jGa9J9TkqtBcUoHe2zqhVFFbgUVED6o9or"
amount.Value = "0.5"
destination = "jB9eHCFeCaoxw6d9V9pBx5hiKUGW9K2fbs"
req, _ := remote.RequestOrderBook(options)
req.Submit(func(err error, result interface{}) {
    jsonBytes, _ := json.Marshal(result)
	t.Logf("Success request order book : %s",jsonBytes)
	})
```

In this path find, the user wants to send CNY to another account. The system provides one choice which is to use SWT.

In each choice, one `Key` is presented. Key is used to "SetPath" in transaction parameter setting.

### BuildPaymentTx(options)
Normal payment transaction. 

More parameters can be set by Transaction members. The secret is requried, and others are optional.

#### options
* account: The source address.
* to: The destination address.
* amount: The payment amount.

#### sample
```
var v struct {
    	account string
		secret  string
	}
v.account = "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
v.secret = "ssc5eiFivvU2otV6bSYmJeZrAsQK3"
to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
amount := constant.Amount{}
amount.Currency = "SWT"
amount.Value = "0.0001"
tx, _ := remote.BuildPaymentTx(v.account, to, amount)
tx.SetSecret(v.secret)
tx.AddMemo("支付0.0001SWT")
tx.Submit(func(err error, result interface{}) {
    jsonByte, _ := json.Marshal(result)
    t.Logf("Success Payment result : %s", jsonByte)
	})
```

### BuildRelationTx(options)
Build relation Transaction. Now Jingtum supports "trust", "authorize" and "freeze" relation setting.

Same as payment transaction parameter setting, secret is required and others are optional.

#### options
* account: The source address.
* target: The target address.
* type: The relation type. "Trust", "Authorize", "Freeze".
* limit: The limit amount.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust", "quality_out": 100, "quality_in": 10}
limit := jingtumLib.Amount{}
limit.Currency = "SWT"
limit.Value = "100.0001"
limit.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
options["limit"] = limit
req, _ := remote.BuildRelationTx(options)
req.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
req.Submit(func(err error, result interface{}) {
    jsonBytes, _ := json.Marshal(result)
	t.Logf("Success Build Relation Tx result : %s", jsonBytes)
})
```

### BuildAccountSetTx(options)
AccountSet Transaction is used to set account attribute. Now Jingtum supoorts three account attributes setting, as "property", "delegate" and "signer". "property" is used to set normal account info, "delegate" is used to set delegate account for this account, and "signer" is used to set signers for this acccount.

Same as payment transaction parameter setting, secret is required and others are optional.

#### options
* account: The source address.
* type: The property type. "Property", "Delegate", "Signer".
* set_flag: (optional) The attribute to set for property type.
* clear_flag: (optional) The attribute to remove for property type.
* delegate_key: (optional) The regualar address for delegate type.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "property", "set_flag": "asfRequireDest", "clear": "asfDisableMaster", "target": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"}
req, _ := remote.BuildAccountSetTx(options)
req.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
req.Submit(func(err error, result interface{}) {
	jsonBytes, _ := json.Marshal(result)
	t.Logf("Success Build AccountSet Tx result : %s", jsonBytes)
})
```

### BuildOfferCreateTx(options)
Create one offer and submit to system. 

#### options
* account: The source address.
* type: "Sell" or "Buy".
* gets: The amount to get by taker.
* pays: The amount to pay by taker.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "Sell",  "set_flag": "asfRequireDest", "clear": "asfDisableMaster"}
gets := jingtumLib.Amount{}
gets.Currency = "SWT"
pays := jingtumLib.Amount{}
pays.Currency = "CNY"
options["gets"] = gets
options["pays"] = pays
req, _ := remote.BuildOfferCreateTx(options)
req.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
req.Submit(func(err error, result interface{}) {
    jsonBytes, _ := json.Marshal(result)
	t.Logf("Success Build Offer Create Tx result : %s", jsonBytes)
})
```

### BuildOfferCancelTx(options)
Order can be canceled by order sequence. The sequence can be get when order is submitted or from offer query operation.

#### options
* account: The account address.
* sequence: The order sequence. It can be get from RequestAccountOffers operation.

#### sample
```
options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "sequence": uint32(26)}
tx, _ := remote.BuildOfferCancelTx(options)
tx.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
tx.Submit(func(err error, result interface{}) {
	jsonBytes, _ := json.Marshal(result)
	t.Logf("Success BuildOfferCancelTx : %s", jsonBytes)
})
```

### DeployContractTx(options)
Deploy contract to the system. The contract address is returned in the ContractState property.

#### options
* account: The source address.
* amount: The swt to active the contract address.
* paylaod: The lua scripts.
* params: (optional) The parameters.

#### sample
```
options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "amount": float64(100), "payload": fmt.Sprintf("%X", "result={}; function Init(t) result=scGetAccountBalance(t) return result end; function foo(t) result=scGetAccountBalance(t) return result end"), "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
tx, _ := remote.DeployContractTx(options)
tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
tx.Submit(func(err error, data interface{}) {
	jsonBytes, _ := json.Marshal(data)
	t.Logf("Success deploy contract : %s", string(jsonBytes))
})
```

### CallContractTx(options)
Call the contract. The call result is returned in the ContractState property.

#### options
* account: The source address.
* destination: The contract address.
* foo: The function name to call.
* params: (optional) The parameters.

#### sample
```
options := map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "destination": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72", "foo": "foo", "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
tx, _ := remote.CallContractTx(options)
tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
tx.Submit(func(err error, data interface{}) {
	jsonBytes, _ := json.Marshal(data)
	t.Logf("Success call contract Tx : %s", string(jsonBytes))
})
```

### Events

#### Transactions
* Listening all transactions occur in the system.

#### LedgerClosed
* Listening all last closed ledger event.

#### ServerStatusChanged
* Listening all server status change event.

## Request

Request is used to get server, account, orderbook and path info. Request is not secret required, and will be public to every one. All requests are asynchronized and should provide a callback. Each callback returns the raw json message, exception and parsed result.

* SelectLedger(ledger)
* Submit(callback)

### SelectLedger(ledger)

Select one ledger for current request, ledger can be follow options,

* ledger index: The ledger index.
* ledger hash: The ledger hash.
* ledger state: The ledger state. "Current", "Validated", "Closed". 

After ledger is selected, the result is for the specified ledger.

### Submit(callback)

Callback entry for request. Each callback returns error and parsed result.

* error: The exception for local argument validation or error message from the jingtum system.
* result: The parsed result object.


## Transaction

Transaction is used to make transaction and collect transaction parameter. Each transaction is secret required, and transaction can be signed local or remote. All transactions are should provide a callback. Each callback returns error and parsed result.

* Account (get)
* TransactionType (get)
* SetSecret(secret)
* AddMemo(memo)
* SetPath(key)
* SetSendMax(amount)
* SetTransferRate(rate)
* SetFlags(flags)
* Submit(callback)

### Account property
Each transaction has source address, and its secret should be set.

Account can be master account, delegate account or operation account.

### TransactionType property

Get transaction type. Now Jingtum supports `Payment`, `OfferCreate`, `OfferCancel`, `AccountSet` and so on. 

### SetSecret(secret)

Set Transaction secret, this method is required before transaction submit.

### AddMemo(memo)

Add one memo to transaction, memo is string and is limited to 2k.

### SetPath(key)

Set path for one transaction. The key parameter is request by RequestPathFind method. When the key is set, "SendMax" parameter is also set.

### SetSendMax(amount)

Set payment transaction max amount when needed. It is set by "SetPath" default.

### SetTransferRate(rate)

Set transaction transfer rate. It should be check with fee. 

### SetFlags(flags)

Set transaction flags. It is used to set Offer type mainly. As follows

```
SetFlags((UInt32)OfferCreateFlags.Sell)
```
    
### Submit(callback)

Submit entry for transaction. Each callback returns the error and parsed result.

* error: The exception for local argument validation or error message from the jingtum system.
* result: The parsed result object.