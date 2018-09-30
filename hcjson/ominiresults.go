package hcjson

type OmniSendResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSenddexsellResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSenddexacceptResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendissuancecrowdsaleResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendissuancefixedResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendissuancemanagedResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendstoResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendgrantResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendrevokeResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendclosecrowdsaleResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendtradeResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendcanceltradesbypriceResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendcanceltradesbypairResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendcancelalltradesResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendchangeissuerResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendallResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendenablefreezingResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSenddisablefreezingResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendfreezeResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendunfreezeResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniSendrawtxResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniFundedSendResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniFundedSendallResult struct {
	/*
		"hash"  // (string) the hex-encoded transaction hash
	*/
}

type OmniGetinfoResult struct {
	/*
		{
		  "omnicoreversion_int" : xxxxxxx,      // (number) client version as integer
		  "omnicoreversion" : "x.x.x.x-xxx",    // (string) client version
		  "mastercoreversion" : "x.x.x.x-xxx",  // (string) client version (DEPRECIATED)
		  "bitcoincoreversion" : "x.x.x",       // (string) Bitcoin Core version
		  "commitinfo" : "xxxxxxx",             // (string) build commit identifier
		  "block" : nnnnnn,                     // (number) index of the last processed block
		  "blocktime" : nnnnnnnnnn,             // (number) timestamp of the last processed block
		  "blocktransactions" : nnnn,           // (number) Omni transactions found in the last processed block
		  "totaltransactions" : nnnnnnnn,       // (number) Omni transactions processed in total
		  "alerts" : [                          // (array of JSON objects) active protocol alert (if any)
			{
			  "alerttype" : n                       // (number) alert type as integer
			  "alerttype" : "xxx"                   // (string) alert type (can be "alertexpiringbyblock", "alertexpiringbyblocktime", "alertexpiringbyclientversion" or "error")
			  "alertexpiry" : "nnnnnnnnnn"          // (string) expiration criteria (can refer to block height, timestamp or client verion)
			  "alertmessage" : "xxx"                // (string) information about the alert
			},
			...
		  ]
		}
	*/
}

type OmniGetbalanceResult struct {
	/*
		{
		  "balance" : "n.nnnnnnnn",  // (string) the available balance of the address
		  "reserved" : "n.nnnnnnnn", // (string) the amount reserved by sell offers and accepts
		  "frozen" : "n.nnnnnnnn"    // (string) the amount frozen by the issuer (applies to managed properties only)
		}
	*/
}

type OmniGetallbalancesforidResult struct {
	/*
		[                          // (array of JSON objects)
		  {
			"address" : "address",     // (string) the address
			"balance" : "n.nnnnnnnn",  // (string) the available balance of the address
			"reserved" : "n.nnnnnnnn", // (string) the amount reserved by sell offers and accepts
			"frozen" : "n.nnnnnnnn"    // (string) the amount frozen by the issuer (applies to managed properties only)
		  },
		  ...
		]
	*/
}

type OmniGetallbalancesforaddressResult struct {
	/*
		[                          // (array of JSON objects)
		  {
			"propertyid" : n,          // (number) the property identifier
			"name" : "name",           // (string) the name of the property
			"balance" : "n.nnnnnnnn",  // (string) the available balance of the address
			"reserved" : "n.nnnnnnnn", // (string) the amount reserved by sell offers and accepts
			"frozen" : "n.nnnnnnnn"    // (string) the amount frozen by the issuer (applies to managed properties only)
		  },
		  ...
		]
	*/
}

type OmniGetwalletbalancesResult struct {
	/*
		[                           // (array of JSON objects)
		  {
			"propertyid" : n,         // (number) the property identifier
			"name" : "name",            // (string) the name of the token
			"balance" : "n.nnnnnnnn",   // (string) the total available balance for the token
			"reserved" : "n.nnnnnnnn"   // (string) the total amount reserved by sell offers and accepts
			"frozen" : "n.nnnnnnnn"     // (string) the total amount frozen by the issuer (applies to managed properties only)
		  },
		  ...
		]
	*/
}

type OmniGetwalletaddressbalancesResult struct {
	/*
		[                           // (array of JSON objects)
		  {
			"address" : "address",      // (string) the address linked to the following balances
			"balances" :
			[
			  {
				"propertyid" : n,         // (number) the property identifier
				"name" : "name",            // (string) the name of the token
				"balance" : "n.nnnnnnnn",   // (string) the available balance for the token
				"reserved" : "n.nnnnnnnn"   // (string) the amount reserved by sell offers and accepts
				"frozen" : "n.nnnnnnnn"     // (string) the amount frozen by the issuer (applies to managed properties only)
			  },
			  ...
			]
		  },
		  ...
		]
	*/
}

type OmniGettransactionResult struct {
	/*
		{
		  "txid" : "hash",                 // (string) the hex-encoded hash of the transaction
		  "sendingaddress" : "address",    // (string) the Bitcoin address of the sender
		  "referenceaddress" : "address",  // (string) a Bitcoin address used as reference (if any)
		  "ismine" : true|false,           // (boolean) whether the transaction involes an address in the wallet
		  "confirmations" : nnnnnnnnnn,    // (number) the number of transaction confirmations
		  "fee" : "n.nnnnnnnn",            // (string) the transaction fee in bitcoins
		  "blocktime" : nnnnnnnnnn,        // (number) the timestamp of the block that contains the transaction
		  "valid" : true|false,            // (boolean) whether the transaction is valid
		  "positioninblock" : n,           // (number) the position (index) of the transaction within the block
		  "version" : n,                   // (number) the transaction version
		  "type_int" : n,                  // (number) the transaction type as number
		  "type" : "type",                 // (string) the transaction type as string
		  [...]                            // (mixed) other transaction type specific properties
		}
	*/
}

type OmniListtransactionsResult struct {
	/*
		[                                // (array of JSON objects)
		  {
			"txid" : "hash",                 // (string) the hex-encoded hash of the transaction
			"sendingaddress" : "address",    // (string) the Bitcoin address of the sender
			"referenceaddress" : "address",  // (string) a Bitcoin address used as reference (if any)
			"ismine" : true|false,           // (boolean) whether the transaction involves an address in the wallet
			"confirmations" : nnnnnnnnnn,    // (number) the number of transaction confirmations
			"fee" : "n.nnnnnnnn",            // (string) the transaction fee in bitcoins
			"blocktime" : nnnnnnnnnn,        // (number) the timestamp of the block that contains the transaction
			"valid" : true|false,            // (boolean) whether the transaction is valid
			"positioninblock" : n,           // (number) the position (index) of the transaction within the block
			"version" : n,                   // (number) the transaction version
			"type_int" : n,                  // (number) the transaction type as number
			"type" : "type",                 // (string) the transaction type as string
			[...]                            // (mixed) other transaction type specific properties
		  },
		  ...
		]
	*/
}

type OmniListblocktransactionsResult struct {
	/*
		[      // (array of string)
		  "hash",  // (string) the hash of the transaction
		  ...
		]
	*/
}

type OmniListpendingtransactionsResult struct {
	/*
		[                                // (array of JSON objects)
		  {
			"txid" : "hash",                 // (string) the hex-encoded hash of the transaction
			"sendingaddress" : "address",    // (string) the Bitcoin address of the sender
			"referenceaddress" : "address",  // (string) a Bitcoin address used as reference (if any)
			"ismine" : true|false,           // (boolean) whether the transaction involes an address in the wallet
			"fee" : "n.nnnnnnnn",            // (string) the transaction fee in bitcoins
			"version" : n,                   // (number) the transaction version
			"type_int" : n,                  // (number) the transaction type as number
			"type" : "type",                 // (string) the transaction type as string
			[...]                            // (mixed) other transaction type specific properties
		  },
		  ...
		]
	*/
}

type OmniGetactivedexsellsResult struct {
	/*
		[                                  // (array of JSON objects)
		  {
			"txid" : "hash",                   // (string) the hash of the transaction of this offer
			"propertyid" : n,                  // (number) the identifier of the tokens for sale
			"seller" : "address",              // (string) the Bitcoin address of the seller
			"amountavailable" : "n.nnnnnnnn",  // (string) the number of tokens still listed for sale and currently available
			"bitcoindesired" : "n.nnnnnnnn",   // (string) the number of bitcoins desired in exchange
			"unitprice" : "n.nnnnnnnn" ,       // (string) the unit price (BTC/token)
			"timelimit" : nn,                  // (number) the time limit in blocks a buyer has to pay following a successful accept
			"minimumfee" : "n.nnnnnnnn",       // (string) the minimum mining fee a buyer has to pay to accept this offer
			"amountaccepted" : "n.nnnnnnnn",   // (string) the number of tokens currently reserved for pending "accept" orders
			"accepts": [                       // (array of JSON objects) a list of pending "accept" orders
			  {
				"buyer" : "address",               // (string) the Bitcoin address of the buyer
				"block" : nnnnnn,                  // (number) the index of the block that contains the "accept" order
				"blocksleft" : nn,                 // (number) the number of blocks left to pay
				"amount" : "n.nnnnnnnn"            // (string) the amount of tokens accepted and reserved
				"amounttopay" : "n.nnnnnnnn"       // (string) the amount in bitcoins needed finalize the trade
			  },
			  ...
			]
		  },
		  ...
		]
	*/
}

type OmniListpropertiesResult struct {
	/*
		[                               // (array of JSON objects)
		  {
			"propertyid" : n,               // (number) the identifier of the tokens
			"name" : "name",                // (string) the name of the tokens
			"category" : "category",        // (string) the category used for the tokens
			"subcategory" : "subcategory",  // (string) the subcategory used for the tokens
			"data" : "information",         // (string) additional information or a description
			"url" : "uri",                  // (string) an URI, for example pointing to a website
			"divisible" : true|false        // (boolean) whether the tokens are divisible
		  },
		  ...
		]
	*/
}

type OmniGetpropertyResult struct {
	/*
		{
		  "propertyid" : n,               // (number) the identifier
		  "name" : "name",                // (string) the name of the tokens
		  "category" : "category",        // (string) the category used for the tokens
		  "subcategory" : "subcategory",  // (string) the subcategory used for the tokens
		  "data" : "information",         // (string) additional information or a description
		  "url" : "uri",                  // (string) an URI, for example pointing to a website
		  "divisible" : true|false,       // (boolean) whether the tokens are divisible
		  "issuer" : "address",           // (string) the Bitcoin address of the issuer on record
		  "creationtxid" : "hash",        // (string) the hex-encoded creation transaction hash
		  "fixedissuance" : true|false,   // (boolean) whether the token supply is fixed
		  "managedissuance" : true|false, // (boolean) whether the token supply is managed by the issuer
		  "freezingenabled" : true|false, // (boolean) whether freezing is enabled for the property (managed properties only)
		  "totaltokens" : "n.nnnnnnnn"    // (string) the total number of tokens in existence
		}
	*/
}

type OmniGetactivecrowdsalesResult struct {
	/*
		[                                // (array of JSON objects)
		  {
			"propertyid" : n,                // (number) the identifier of the crowdsale
			"name" : "name",                 // (string) the name of the tokens issued via the crowdsale
			"issuer" : "address",            // (string) the Bitcoin address of the issuer on record
			"propertyiddesired" : n,         // (number) the identifier of the tokens eligible to participate in the crowdsale
			"tokensperunit" : "n.nnnnnnnn",  // (string) the amount of tokens granted per unit invested in the crowdsale
			"earlybonus" : n,                // (number) an early bird bonus for participants in percent per week
			"percenttoissuer" : n,           // (number) a percentage of tokens that will be granted to the issuer
			"starttime" : nnnnnnnnnn,        // (number) the start time of the of the crowdsale as Unix timestamp
			"deadline" : nnnnnnnnnn          // (number) the deadline of the crowdsale as Unix timestamp
		  },
		  ...
		]
	*/
}

type OmniGetcrowdsaleResult struct {
	/*
		{
		  "propertyid" : n,                    // (number) the identifier of the crowdsale
		  "name" : "name",                     // (string) the name of the tokens issued via the crowdsale
		  "active" : true|false,               // (boolean) whether the crowdsale is still active
		  "issuer" : "address",                // (string) the Bitcoin address of the issuer on record
		  "propertyiddesired" : n,             // (number) the identifier of the tokens eligible to participate in the crowdsale
		  "tokensperunit" : "n.nnnnnnnn",      // (string) the amount of tokens granted per unit invested in the crowdsale
		  "earlybonus" : n,                    // (number) an early bird bonus for participants in percent per week
		  "percenttoissuer" : n,               // (number) a percentage of tokens that will be granted to the issuer
		  "starttime" : nnnnnnnnnn,            // (number) the start time of the of the crowdsale as Unix timestamp
		  "deadline" : nnnnnnnnnn,             // (number) the deadline of the crowdsale as Unix timestamp
		  "amountraised" : "n.nnnnnnnn",       // (string) the amount of tokens invested by participants
		  "tokensissued" : "n.nnnnnnnn",       // (string) the total number of tokens issued via the crowdsale
		  "issuerbonustokens" : "n.nnnnnnnn",  // (string) the amount of tokens granted to the issuer as bonus
		  "addedissuertokens" : "n.nnnnnnnn",  // (string) the amount of issuer bonus tokens not yet emitted
		  "closedearly" : true|false,          // (boolean) whether the crowdsale ended early (if not active)
		  "maxtokens" : true|false,            // (boolean) whether the crowdsale ended early due to reaching the limit of max. issuable tokens (if not active)
		  "endedtime" : nnnnnnnnnn,            // (number) the time when the crowdsale ended (if closed early)
		  "closetx" : "hash",                  // (string) the hex-encoded hash of the transaction that closed the crowdsale (if closed manually)
		  "participanttransactions": [         // (array of JSON objects) a list of crowdsale participations (if verbose=true)
			{
			  "txid" : "hash",                     // (string) the hex-encoded hash of participation transaction
			  "amountsent" : "n.nnnnnnnn",         // (string) the amount of tokens invested by the participant
			  "participanttokens" : "n.nnnnnnnn",  // (string) the tokens granted to the participant
			  "issuertokens" : "n.nnnnnnnn"        // (string) the tokens granted to the issuer as bonus
			},
			...
		  ]
		}
	*/
}

type OmniGetgrantsResult struct {
	/*
		{
		  "propertyid" : n,              // (number) the identifier of the managed tokens
		  "name" : "name",               // (string) the name of the tokens
		  "issuer" : "address",          // (string) the Bitcoin address of the issuer on record
		  "creationtxid" : "hash",       // (string) the hex-encoded creation transaction hash
		  "totaltokens" : "n.nnnnnnnn",  // (string) the total number of tokens in existence
		  "issuances": [                 // (array of JSON objects) a list of the granted and revoked tokens
			{
			  "txid" : "hash",               // (string) the hash of the transaction that granted tokens
			  "grant" : "n.nnnnnnnn"         // (string) the number of tokens granted by this transaction
			},
			{
			  "txid" : "hash",               // (string) the hash of the transaction that revoked tokens
			  "grant" : "n.nnnnnnnn"         // (string) the number of tokens revoked by this transaction
			},
			...
		  ]
		}
	*/
}

type OmniGetstoResult struct {
	/*
		{
		  "txid" : "hash",               // (string) the hex-encoded hash of the transaction
		  "sendingaddress" : "address",  // (string) the Bitcoin address of the sender
		  "ismine" : true|false,         // (boolean) whether the transaction involes an address in the wallet
		  "confirmations" : nnnnnnnnnn,  // (number) the number of transaction confirmations
		  "fee" : "n.nnnnnnnn",          // (string) the transaction fee in bitcoins
		  "blocktime" : nnnnnnnnnn,      // (number) the timestamp of the block that contains the transaction
		  "valid" : true|false,          // (boolean) whether the transaction is valid
		  "positioninblock" : n,         // (number) the position (index) of the transaction within the block
		  "version" : n,                 // (number) the transaction version
		  "type_int" : n,                // (number) the transaction type as number
		  "type" : "type",               // (string) the transaction type as string
		  "propertyid" : n,              // (number) the identifier of sent tokens
		  "divisible" : true|false,      // (boolean) whether the sent tokens are divisible
		  "amount" : "n.nnnnnnnn",       // (string) the number of tokens sent to owners
		  "totalstofee" : "n.nnnnnnnn",  // (string) the fee paid by the sender, nominated in OMNI or TOMNI
		  "recipients": [                // (array of JSON objects) a list of recipients
			{
			  "address" : "address",         // (string) the Bitcoin address of the recipient
			  "amount" : "n.nnnnnnnn"        // (string) the number of tokens sent to this recipient
			},
			...
		  ]
		}
	*/
}

type OmniGettradeResult struct {
	/*
		{
		  "txid" : "hash",                              // (string) the hex-encoded hash of the transaction of the order
		  "sendingaddress" : "address",                 // (string) the Bitcoin address of the trader
		  "ismine" : true|false,                        // (boolean) whether the order involes an address in the wallet
		  "confirmations" : nnnnnnnnnn,                 // (number) the number of transaction confirmations
		  "fee" : "n.nnnnnnnn",                         // (string) the transaction fee in bitcoins
		  "blocktime" : nnnnnnnnnn,                     // (number) the timestamp of the block that contains the transaction
		  "valid" : true|false,                         // (boolean) whether the transaction is valid
		  "positioninblock" : n,                        // (number) the position (index) of the transaction within the block
		  "version" : n,                                // (number) the transaction version
		  "type_int" : n,                               // (number) the transaction type as number
		  "type" : "type",                              // (string) the transaction type as string
		  "propertyidforsale" : n,                      // (number) the identifier of the tokens put up for sale
		  "propertyidforsaleisdivisible" : true|false,  // (boolean) whether the tokens for sale are divisible
		  "amountforsale" : "n.nnnnnnnn",               // (string) the amount of tokens initially offered
		  "propertyiddesired" : n,                      // (number) the identifier of the tokens desired in exchange
		  "propertyiddesiredisdivisible" : true|false,  // (boolean) whether the desired tokens are divisible
		  "amountdesired" : "n.nnnnnnnn",               // (string) the amount of tokens initially desired
		  "unitprice" : "n.nnnnnnnnnnn..."              // (string) the unit price (shown in the property desired)
		  "status" : "status"                           // (string) the status of the order ("open", "cancelled", "filled", ...)
		  "canceltxid" : "hash",                        // (string) the hash of the transaction that cancelled the order (if cancelled)
		  "matches": [                                  // (array of JSON objects) a list of matched orders and executed trades
			{
			  "txid" : "hash",                              // (string) the hash of the transaction that was matched against
			  "block" : nnnnnn,                             // (number) the index of the block that contains this transaction
			  "address" : "address",                        // (string) the Bitcoin address of the other trader
			  "amountsold" : "n.nnnnnnnn",                  // (string) the number of tokens sold in this trade
			  "amountreceived" : "n.nnnnnnnn"               // (string) the number of tokens traded in exchange
			},
			...
		  ]
		}
	*/
}

type OmniGetorderbookResult struct {
	/*
		[                                             // (array of JSON objects)
		  {
			"address" : "address",                        // (string) the Bitcoin address of the trader
			"txid" : "hash",                              // (string) the hex-encoded hash of the transaction of the order
			"ecosystem" : "main"|"test",                  // (string) the ecosytem in which the order was made (if "cancel-ecosystem")
			"propertyidforsale" : n,                      // (number) the identifier of the tokens put up for sale
			"propertyidforsaleisdivisible" : true|false,  // (boolean) whether the tokens for sale are divisible
			"amountforsale" : "n.nnnnnnnn",               // (string) the amount of tokens initially offered
			"amountremaining" : "n.nnnnnnnn",             // (string) the amount of tokens still up for sale
			"propertyiddesired" : n,                      // (number) the identifier of the tokens desired in exchange
			"propertyiddesiredisdivisible" : true|false,  // (boolean) whether the desired tokens are divisible
			"amountdesired" : "n.nnnnnnnn",               // (string) the amount of tokens initially desired
			"amounttofill" : "n.nnnnnnnn",                // (string) the amount of tokens still needed to fill the offer completely
			"action" : n,                                 // (number) the action of the transaction: (1) "trade", (2) "cancel-price", (3) "cancel-pair", (4) "cancel-ecosystem"
			"block" : nnnnnn,                             // (number) the index of the block that contains the transaction
			"blocktime" : nnnnnnnnnn                      // (number) the timestamp of the block that contains the transaction
		  },
		  ...
		]
	*/
}

type OmniGettradehistoryforpairResult struct {
	/*
		[                                     // (array of JSON objects)
		  {
			"block" : nnnnnn,                     // (number) the index of the block that contains the trade match
			"unitprice" : "n.nnnnnnnnnnn..." ,    // (string) the unit price used to execute this trade (received/sold)
			"inverseprice" : "n.nnnnnnnnnnn...",  // (string) the inverse unit price (sold/received)
			"sellertxid" : "hash",                // (string) the hash of the transaction of the seller
			"address" : "address",                // (string) the Bitcoin address of the seller
			"amountsold" : "n.nnnnnnnn",          // (string) the number of tokens sold in this trade
			"amountreceived" : "n.nnnnnnnn",      // (string) the number of tokens traded in exchange
			"matchingtxid" : "hash",              // (string) the hash of the transaction that was matched against
			"matchingaddress" : "address"         // (string) the Bitcoin address of the other party of this trade
		  },
		  ...
		]
	*/
}

type OmniGettradehistoryforaddressResult struct {
	/*
		[                                             // (array of JSON objects)
		  {
			"txid" : "hash",                              // (string) the hex-encoded hash of the transaction of the order
			"sendingaddress" : "address",                 // (string) the Bitcoin address of the trader
			"ismine" : true|false,                        // (boolean) whether the order involes an address in the wallet
			"confirmations" : nnnnnnnnnn,                 // (number) the number of transaction confirmations
			"fee" : "n.nnnnnnnn",                         // (string) the transaction fee in bitcoins
			"blocktime" : nnnnnnnnnn,                     // (number) the timestamp of the block that contains the transaction
			"valid" : true|false,                         // (boolean) whether the transaction is valid
			"positioninblock" : n,                        // (number) the position (index) of the transaction within the block
			"version" : n,                                // (number) the transaction version
			"type_int" : n,                               // (number) the transaction type as number
			"type" : "type",                              // (string) the transaction type as string
			"propertyidforsale" : n,                      // (number) the identifier of the tokens put up for sale
			"propertyidforsaleisdivisible" : true|false,  // (boolean) whether the tokens for sale are divisible
			"amountforsale" : "n.nnnnnnnn",               // (string) the amount of tokens initially offered
			"propertyiddesired" : n,                      // (number) the identifier of the tokens desired in exchange
			"propertyiddesiredisdivisible" : true|false,  // (boolean) whether the desired tokens are divisible
			"amountdesired" : "n.nnnnnnnn",               // (string) the amount of tokens initially desired
			"unitprice" : "n.nnnnnnnnnnn..."              // (string) the unit price (shown in the property desired)
			"status" : "status"                           // (string) the status of the order ("open", "cancelled", "filled", ...)
			"canceltxid" : "hash",                        // (string) the hash of the transaction that cancelled the order (if cancelled)
			"matches": [                                  // (array of JSON objects) a list of matched orders and executed trades
			  {
				"txid" : "hash",                              // (string) the hash of the transaction that was matched against
				"block" : nnnnnn,                             // (number) the index of the block that contains this transaction
				"address" : "address",                        // (string) the Bitcoin address of the other trader
				"amountsold" : "n.nnnnnnnn",                  // (string) the number of tokens sold in this trade
				"amountreceived" : "n.nnnnnnnn"               // (string) the number of tokens traded in exchange
			  },
			  ...
			]
		  },
		  ...
		]
	*/
}

type OmniGetactivationsResult struct {
	/*
		{
		  "pendingactivations": [      // (array of JSON objects) a list of pending feature activations
			{
			  "featureid" : n,             // (number) the id of the feature
			  "featurename" : "xxxxxxxx",  // (string) the name of the feature
			  "activationblock" : n,       // (number) the block the feature will be activated
			  "minimumversion" : n         // (number) the minimum client version needed to support this feature
			},
			...
		  ]
		  "completedactivations": [    // (array of JSON objects) a list of completed feature activations
			{
			  "featureid" : n,             // (number) the id of the feature
			  "featurename" : "xxxxxxxx",  // (string) the name of the feature
			  "activationblock" : n,       // (number) the block the feature will be activated
			  "minimumversion" : n         // (number) the minimum client version needed to support this feature
			},
			...
		  ]
		}
	*/
}

type OmniGetpayloadResult struct {
	/*
		{
		  "payload" : "payloadmessage",  // (string) the decoded Omni payload message
		  "payloadsize" : n              // (number) the size of the payload
		}
	*/
}

type OmniGetseedblocksResult struct {
	/*
		[         // (array of numbers) a list of seed blocks
		  nnnnnnn,  // the block height of the seed block
		  ...
		]
	*/
}

type OmniGetcurrentconsensushashResult struct {
	/*
		{
		  "block" : nnnnnn,         // (number) the index of the block this consensus hash applies to
		  "blockhash" : "hash",     // (string) the hash of the corresponding block
		  "consensushash" : "hash"  // (string) the consensus hash for the block
		}
	*/
}

type OmniDecodetransactionResult struct {
	/*
		{
		  "txid" : "hash",                 // (string) the hex-encoded hash of the transaction
		  "fee" : "n.nnnnnnnn",            // (string) the transaction fee in bitcoins
		  "sendingaddress" : "address",    // (string) the Bitcoin address of the sender
		  "referenceaddress" : "address",  // (string) a Bitcoin address used as reference (if any)
		  "ismine" : true|false,           // (boolean) whether the transaction involes an address in the wallet
		  "version" : n,                   // (number) the transaction version
		  "type_int" : n,                  // (number) the transaction type as number
		  "type" : "type",                 // (string) the transaction type as string
		  [...]                            // (mixed) other transaction type specific properties
		}
	*/
}

type OmniCreaterawtxOpreturnResult struct {
	/*
		"rawtx"  // (string) the hex-encoded modified raw transaction
	*/
}

type OmniCreaterawtxMultisigResult struct {
	/*
		"rawtx"  // (string) the hex-encoded modified raw transaction
	*/
}

type OmniCreaterawtxInputResult struct {
	/*
		"rawtx"  // (string) the hex-encoded modified raw transaction
	*/
}

type OmniCreaterawtxReferenceResult struct {
	/*
		"rawtx"  // (string) the hex-encoded modified raw transaction
	*/
}

type OmniCreaterawtxChangeResult struct {
	/*
		"rawtx"  // (string) the hex-encoded modified raw transaction
	*/
}

type OmniCreatepayloadSimplesendResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadSendallResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadDexsellResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadDexacceptResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadStoResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadIssuancefixedResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadIssuancecrowdsaleResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadIssuancemanagedResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadClosecrowdsaleResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadGrantResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadRevokeResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadChangeissuerResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadTradeResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadCanceltradesbypriceResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadCanceltradesbypairResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadCancelalltradesResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadEnablefreezingResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadDisablefreezingResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadFreezeResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniCreatepayloadUnfreezeResult struct {
	/*
		"payload"  // (string) the hex-encoded payload
	*/
}

type OmniGetfeecacheResult struct {
	/*
		[                                  // (array of JSON objects)
		  {
			"propertyid" : nnnnnnn,        // (number) the property id
			"cachedfees" : "n.nnnnnnnn",   // (string) the amount of fees cached for this property
		  },
		...
		]
	*/
}

type OmniGetfeetriggerResult struct {
	/*
		[                                  // (array of JSON objects)
		  {
			"propertyid" : nnnnnnn,        // (number) the property id
			"feetrigger" : "n.nnnnnnnn",   // (string) the amount of fees required to trigger distribution
		  },
		...
		]
	*/
}

type OmniGetfeeshareResult struct {
	/*
		[                                  // (array of JSON objects)
		  {
			"address" : "address"          // (string) the adress that would receive a share of fees
			"feeshare" : "n.nnnn%",        // (string) the percentage of fees this address will receive based on the current state
		  },
		...
		]
	*/
}

type OmniGetfeedistributionResult struct {
	/*
		{
		  "distributionid" : n,            // (number) the distribution id
		  "propertyid" : n,                // (number) the property id of the distributed tokens
		  "block" : n,                     // (number) the block the distribution occurred
		  "amount" : "n.nnnnnnnn",         // (string) the amount that was distributed
		  "recipients": [                  // (array of JSON objects) a list of recipients
			{
			  "address" : "address",       // (string) the address of the recipient
			  "amount" : "n.nnnnnnnn"      // (string) the amount of fees received by the recipient
			},
			...
		  ]
		}
	*/
}

type OmniGetfeedistributionsResult struct {
	/*
		[                                  // (array of JSON objects)
		  {
			"distributionid" : n,          // (number) the distribution id
			"propertyid" : n,              // (number) the property id of the distributed tokens
			"block" : n,                   // (number) the block the distribution occurred
			"amount" : "n.nnnnnnnn",       // (string) the amount that was distributed
			"recipients": [                // (array of JSON objects) a list of recipients
			  {
				"address" : "address",       // (string) the address of the recipient
				"amount" : "n.nnnnnnnn"      // (string) the amount of fees received by the recipient
			  },
			  ...
			]
		  },
		  ...
		]
	*/
}

type OmniSetautocommitResult struct {
	/*
		true|false  // (boolean) the updated flag status
	*/
}
