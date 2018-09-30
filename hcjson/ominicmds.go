package hcjson

// OmniSend // Create and broadcast a simple send transaction.
// example: $ omnicore-cli "omni_send" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" "37FaKponF7zqoMLUjEiko25pDiuVH5YLEa" 1 "100.0"
type OmniSendCmd struct {
	Fromaddress     string  `json:"fromaddress" desc:"the address to send from"`
	Toaddress       string  `json:"toaddress" desc:"the address of the receiver"`
	Propertyid      int64   `json:"propertyid" desc:"the identifier of the tokens to send"`
	Amount          string  `json:"amount" desc:"the amount to send"`
	Redeemaddress   *string `json:"redeemaddress" desc:"an address that can spend the transaction dust (sender by default)"`
	Referenceamount *string `json:"referenceamount" desc:"a bitcoin amount that is sent to the receiver (minimal by default)"`
}

func NewOmniSendCmd(fromaddress string, toaddress string, propertyid int64, amount string, redeemaddress *string, referenceamount *string) *OmniSendCmd {
	return &OmniSendCmd{
		Fromaddress:     fromaddress,
		Toaddress:       toaddress,
		Propertyid:      propertyid,
		Amount:          amount,
		Redeemaddress:   redeemaddress,
		Referenceamount: referenceamount,
	}
}

// OmniSenddexsell // Place, update or cancel a sell offer on the traditional distributed OMNI/BTC exchange.
// example: $ omnicore-cli "omni_senddexsell" "37FaKponF7zqoMLUjEiko25pDiuVH5YLEa" 1 "1.5" "0.75" 25 "0.0005" 1
type OmniSenddexsellCmd struct {
	Fromaddress       string `json:"fromaddress" desc:"the address to send from"`
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale (must be 1 for OMNI or 2for TOMNI)"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to list for sale"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of bitcoins desired"`
	Paymentwindow     int64  `json:"paymentwindow" desc:"a time limit in blocks a buyer has to pay following a successful accepting order"`
	Minacceptfee      string `json:"minacceptfee" desc:"a minimum mining fee a buyer has to pay to accept the offer"`
	Action            int64  `json:"action" desc:"the action to take (1 for new offers, 2 to update, 3 to cancel)"`
}

func NewOmniSenddexsellCmd(fromaddress string, propertyidforsale int64, amountforsale string, amountdesired string, paymentwindow int64, minacceptfee string, action int64) *OmniSenddexsellCmd {
	return &OmniSenddexsellCmd{
		Fromaddress:       fromaddress,
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Amountdesired:     amountdesired,
		Paymentwindow:     paymentwindow,
		Minacceptfee:      minacceptfee,
		Action:            action,
	}
}

// OmniSenddexaccept // Create and broadcast an accept offer for the specified token and amount.
// example: $ omnicore-cli "omni_senddexaccept" \     "35URq1NN3xL6GeRKUP6vzaQVcxoJiiJKd8" "37FaKponF7zqoMLUjEiko25pDiuVH5YLEa" 1 "15.0"
type OmniSenddexacceptCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from"`
	Toaddress   string `json:"toaddress" desc:"the address of the seller"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the token to purchase"`
	Amount      string `json:"amount" desc:"the amount to accept"`
	Override    *bool  `json:"override" desc:"override minimum accept fee and payment window checks (use with caution!)"`
}

func NewOmniSenddexacceptCmd(fromaddress string, toaddress string, propertyid int64, amount string, override bool) *OmniSenddexacceptCmd {
	return &OmniSenddexacceptCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
		Amount:      amount,
		Override:    &override,
	}
}

// OmniSendissuancecrowdsale // Create new tokens as crowdsale.
// example: $ omnicore-cli "omni_sendissuancecrowdsale" \     "3JYd75REX3HXn1vAU83YuGfmiPXW7BpYXo" 2 1 0 "Companies" "Bitcoin Mining" \     "Quantum Miner" "" "" 2 "100" 1483228800 30 2
type OmniSendissuancecrowdsaleCmd struct {
	Fromaddress       string `json:"fromaddress" desc:"the address to send from"`
	Ecosystem         int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo              int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid        int64  `json:"previousid" desc:"an identifier of a predecessor token (0 for new crowdsales)"`
	Category          string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory       string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name              string `json:"name" desc:"the name of the new tokens to create"`
	Url               string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data              string `json:"data" desc:"a description for the new tokens (can be "")"`
	Propertyiddesired int64  `json:"propertyiddesired" desc:"the identifier of a token eligible to participate in the crowdsale"`
	Tokensperunit     string `json:"tokensperunit" desc:"the amount of tokens granted per unit invested in the crowdsale"`
	Deadline          int64  `json:"deadline" desc:"the deadline of the crowdsale as Unix timestamp"`
	Earlybonus        int64  `json:"earlybonus" desc:"an early bird bonus for participants in percent per week"`
	Issuerpercentage  int64  `json:"issuerpercentage" desc:"a percentage of tokens that will be granted to the issuer"`
}

func NewOmniSendissuancecrowdsaleCmd(fromaddress string, ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string, propertyiddesired int64, tokensperunit string, deadline int64, earlybonus int64, issuerpercentage int64) *OmniSendissuancecrowdsaleCmd {
	return &OmniSendissuancecrowdsaleCmd{
		Fromaddress:       fromaddress,
		Ecosystem:         ecosystem,
		Typo:              typo,
		Previousid:        previousid,
		Category:          category,
		Subcategory:       subcategory,
		Name:              name,
		Url:               url,
		Data:              data,
		Propertyiddesired: propertyiddesired,
		Tokensperunit:     tokensperunit,
		Deadline:          deadline,
		Earlybonus:        earlybonus,
		Issuerpercentage:  issuerpercentage,
	}
}

// OmniSendissuancefixed // Create new tokens with fixed supply.
// example: $ omnicore-cli "omni_sendissuancefixed" \     "3Ck2kEGLJtZw9ENj2tameMCtS3HB7uRar3" 2 1 0 "Companies" "Bitcoin Mining" \     "Quantum Miner" "" "" "1000000"
type OmniSendissuancefixedCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from"`
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo        int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid  int64  `json:"previousid" desc:"an identifier of a predecessor token (0 for new tokens)"`
	Category    string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name        string `json:"name" desc:"the name of the new tokens to create"`
	Url         string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data        string `json:"data" desc:"a description for the new tokens (can be "")"`
	Amount      string `json:"amount" desc:"the number of tokens to create"`
}

func NewOmniSendissuancefixedCmd(fromaddress string, ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string, amount string) *OmniSendissuancefixedCmd {
	return &OmniSendissuancefixedCmd{
		Fromaddress: fromaddress,
		Ecosystem:   ecosystem,
		Typo:        typo,
		Previousid:  previousid,
		Category:    category,
		Subcategory: subcategory,
		Name:        name,
		Url:         url,
		Data:        data,
		Amount:      amount,
	}
}

// OmniSendissuancemanaged // Create new tokens with manageable supply.
// example: $ omnicore-cli "omni_sendissuancemanaged" \     "3HsJvhr9qzgRe3ss97b1QHs38rmaLExLcH" 2 1 0 "Companies" "Bitcoin Mining" "Quantum Miner" "" ""
type OmniSendissuancemanagedCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from"`
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo        int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid  int64  `json:"previousid" desc:"an identifier of a predecessor token (0 for new tokens)"`
	Category    string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name        string `json:"name" desc:"the name of the new tokens to create"`
	Url         string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data        string `json:"data" desc:"a description for the new tokens (can be "")"`
}

func NewOmniSendissuancemanagedCmd(fromaddress string, ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string) *OmniSendissuancemanagedCmd {
	return &OmniSendissuancemanagedCmd{
		Fromaddress: fromaddress,
		Ecosystem:   ecosystem,
		Typo:        typo,
		Previousid:  previousid,
		Category:    category,
		Subcategory: subcategory,
		Name:        name,
		Url:         url,
		Data:        data,
	}
}

// OmniSendsto // Create and broadcast a send-to-owners transaction.
// example: $ omnicore-cli "omni_sendsto" \     "32Z3tJccZuqQZ4PhJR2hxHC3tjgjA8cbqz" "37FaKponF7zqoMLUjEiko25pDiuVH5YLEa" 3 "5000"
type OmniSendstoCmd struct {
	Fromaddress          string  `json:"fromaddress" desc:"the address to send from"`
	Propertyid           int64   `json:"propertyid" desc:"the identifier of the tokens to distribute"`
	Amount               string  `json:"amount" desc:"the amount to distribute"`
	Redeemaddress        *string `json:"redeemaddress" desc:"an address that can spend the transaction dust (sender by default)"`
	Distributionproperty *int64  `json:"distributionproperty" desc:"the identifier of the property holders to distribute to"`
}

func NewOmniSendstoCmd(fromaddress string, propertyid int64, amount string, redeemaddress *string, distributionproperty *int64) *OmniSendstoCmd {
	return &OmniSendstoCmd{
		Fromaddress:          fromaddress,
		Propertyid:           propertyid,
		Amount:               amount,
		Redeemaddress:        redeemaddress,
		Distributionproperty: distributionproperty,
	}
}

// OmniSendgrant // Issue or grant new units of managed tokens.
// example: $ omnicore-cli "omni_sendgrant" "3HsJvhr9qzgRe3ss97b1QHs38rmaLExLcH" "" 51 "7000"
type OmniSendgrantCmd struct {
	Fromaddress string  `json:"fromaddress" desc:"the address to send from"`
	Toaddress   string  `json:"toaddress" desc:"the receiver of the tokens (sender by default, can be "")"`
	Propertyid  int64   `json:"propertyid" desc:"the identifier of the tokens to grant"`
	Amount      string  `json:"amount" desc:"the amount of tokens to create"`
	Memo        *string `json:"memo" desc:"a text note attached to this transaction (none by default)"`
}

func NewOmniSendgrantCmd(fromaddress string, toaddress string, propertyid int64, amount string, memo *string) *OmniSendgrantCmd {
	return &OmniSendgrantCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
		Amount:      amount,
		Memo:        memo,
	}
}

// OmniSendrevoke // Revoke units of managed tokens.
// example: $ omnicore-cli "omni_sendrevoke" "3HsJvhr9qzgRe3ss97b1QHs38rmaLExLcH" "" 51 "100"
type OmniSendrevokeCmd struct {
	Fromaddress string  `json:"fromaddress" desc:"the address to send from"`
	Propertyid  int64   `json:"propertyid" desc:"the identifier of the tokens to revoke"`
	Amount      string  `json:"amount" desc:"the amount of tokens to revoke"`
	Memo        *string `json:"memo" desc:"a text note attached to this transaction (none by default)"`
}

func NewOmniSendrevokeCmd(fromaddress string, propertyid int64, amount string, memo *string) *OmniSendrevokeCmd {
	return &OmniSendrevokeCmd{
		Fromaddress: fromaddress,
		Propertyid:  propertyid,
		Amount:      amount,
		Memo:        memo,
	}
}

// OmniSendclosecrowdsale // Manually close a crowdsale.
// example: $ omnicore-cli "omni_sendclosecrowdsale" "3JYd75REX3HXn1vAU83YuGfmiPXW7BpYXo" 70
type OmniSendclosecrowdsaleCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address associated with the crowdsale to close"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the crowdsale to close"`
}

func NewOmniSendclosecrowdsaleCmd(fromaddress string, propertyid int64) *OmniSendclosecrowdsaleCmd {
	return &OmniSendclosecrowdsaleCmd{
		Fromaddress: fromaddress,
		Propertyid:  propertyid,
	}
}

// OmniSendtrade // Place a trade offer on the distributed token exchange.
// example: $ omnicore-cli "omni_sendtrade" "3BydPiSLPP3DR5cf726hDQ89fpqWLxPKLR" 31 "250.0" 1 "10.0"
type OmniSendtradeCmd struct {
	Fromaddress       string `json:"fromaddress" desc:"the address to trade with"`
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to list for sale"`
	Propertiddesired  int64  `json:"propertiddesired" desc:"the identifier of the tokens desired in exchange"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of tokens desired in exchange"`
}

func NewOmniSendtradeCmd(fromaddress string, propertyidforsale int64, amountforsale string, propertiddesired int64, amountdesired string) *OmniSendtradeCmd {
	return &OmniSendtradeCmd{
		Fromaddress:       fromaddress,
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Propertiddesired:  propertiddesired,
		Amountdesired:     amountdesired,
	}
}

// OmniSendcanceltradesbyprice // Cancel offers on the distributed token exchange with the specified price.
// example: $ omnicore-cli "omni_sendcanceltradesbyprice" "3BydPiSLPP3DR5cf726hDQ89fpqWLxPKLR" 31 "100.0" 1 "5.0"
type OmniSendcanceltradesbypriceCmd struct {
	Fromaddress       string `json:"fromaddress" desc:"the address to trade with"`
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens listed for sale"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to listed for sale"`
	Propertiddesired  int64  `json:"propertiddesired" desc:"the identifier of the tokens desired in exchange"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of tokens desired in exchange"`
}

func NewOmniSendcanceltradesbypriceCmd(fromaddress string, propertyidforsale int64, amountforsale string, propertiddesired int64, amountdesired string) *OmniSendcanceltradesbypriceCmd {
	return &OmniSendcanceltradesbypriceCmd{
		Fromaddress:       fromaddress,
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Propertiddesired:  propertiddesired,
		Amountdesired:     amountdesired,
	}
}

// OmniSendcanceltradesbypair // Cancel all offers on the distributed token exchange with the given currency pair.
// example: $ omnicore-cli "omni_sendcanceltradesbypair" "3BydPiSLPP3DR5cf726hDQ89fpqWLxPKLR" 1 31
type OmniSendcanceltradesbypairCmd struct {
	Fromaddress       string `json:"fromaddress" desc:"the address to trade with"`
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens listed for sale"`
	Propertiddesired  int64  `json:"propertiddesired" desc:"the identifier of the tokens desired in exchange"`
}

func NewOmniSendcanceltradesbypairCmd(fromaddress string, propertyidforsale int64, propertiddesired int64) *OmniSendcanceltradesbypairCmd {
	return &OmniSendcanceltradesbypairCmd{
		Fromaddress:       fromaddress,
		Propertyidforsale: propertyidforsale,
		Propertiddesired:  propertiddesired,
	}
}

// OmniSendcancelalltrades // Cancel all offers on the distributed token exchange.
// example: $ omnicore-cli "omni_sendcancelalltrades" "3BydPiSLPP3DR5cf726hDQ89fpqWLxPKLR" 1
type OmniSendcancelalltradesCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to trade with"`
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem of the offers to cancel (1 for main ecosystem, 2 for test ecosystem)"`
}

func NewOmniSendcancelalltradesCmd(fromaddress string, ecosystem int64) *OmniSendcancelalltradesCmd {
	return &OmniSendcancelalltradesCmd{
		Fromaddress: fromaddress,
		Ecosystem:   ecosystem,
	}
}

// OmniSendchangeissuer // Change the issuer on record of the given tokens.
// example: $ omnicore-cli "omni_sendchangeissuer" \     "1ARjWDkZ7kT9fwjPrjcQyvbXDkEySzKHwu" "3HTHRxu3aSDV4deakjC7VmsiUp7c6dfbvs" 3
type OmniSendchangeissuerCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address associated with the tokens"`
	Toaddress   string `json:"toaddress  " desc:"the address to transfer administrative control to"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniSendchangeissuerCmd(fromaddress string, toaddress string, propertyid int64) *OmniSendchangeissuerCmd {
	return &OmniSendchangeissuerCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
	}
}

// OmniSendall // Transfers all available tokens in the given ecosystem to the recipient.
// example: $ omnicore-cli "omni_sendall" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" "37FaKponF7zqoMLUjEiko25pDiuVH5YLEa" 2
type OmniSendallCmd struct {
	Fromaddress     string  `json:"fromaddress" desc:"the address to send from"`
	Toaddress       string  `json:"toaddress  " desc:"the address of the receiver"`
	Ecosystem       int64   `json:"ecosystem" desc:"the ecosystem of the tokens to send (1 for main ecosystem, 2 for test ecosystem)"`
	Redeemaddress   *string `json:"redeemaddress" desc:"an address that can spend the transaction dust (sender by default)"`
	Referenceamount *string `json:"referenceamount" desc:"a bitcoin amount that is sent to the receiver (minimal by default)"`
}

func NewOmniSendallCmd(fromaddress string, toaddress string, ecosystem int64, redeemaddress *string, referenceamount *string) *OmniSendallCmd {
	return &OmniSendallCmd{
		Fromaddress:     fromaddress,
		Toaddress:       toaddress,
		Ecosystem:       ecosystem,
		Redeemaddress:   redeemaddress,
		Referenceamount: referenceamount,
	}
}

// OmniSendenablefreezing // Enables address freezing for a centrally managed property.
// example: $ omnicore-cli "omni_sendenablefreezing" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" 2
type OmniSendenablefreezingCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from (must be issuer of a managed property)"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniSendenablefreezingCmd(fromaddress string, propertyid int64) *OmniSendenablefreezingCmd {
	return &OmniSendenablefreezingCmd{
		Fromaddress: fromaddress,
		Propertyid:  propertyid,
	}
}

// OmniSenddisablefreezing // Disables address freezing for a centrally managed property.
// IMPORTANT NOTE:  Disabling freezing for a property will UNFREEZE all frozen addresses for that property!
// example: $ omnicore-cli "omni_senddisablefreezing" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" 2
type OmniSenddisablefreezingCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from (must be issuer of a managed property)"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniSenddisablefreezingCmd(fromaddress string, propertyid int64) *OmniSenddisablefreezingCmd {
	return &OmniSenddisablefreezingCmd{
		Fromaddress: fromaddress,
		Propertyid:  propertyid,
	}
}

// OmniSendfreeze // Freeze an address for a centrally managed token.
// Note: Only the issuer may freeze tokens, and only if the token is of the managed type with the freezing option enabled.
// example: $ omnicore-cli "omni_sendfreeze" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" "3HTHRxu3aSDV4deakjC7VmsiUp7c6dfbvs" 2 1000
type OmniSendfreezeCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from (must be issuer of a managed property with freezing enabled"`
	Toaddress   string `json:"toaddress" desc:"the address to freeze"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens to freeze"`
	Amount      string `json:"amount" desc:"the amount to freeze (note: currently unused, frozen addresses cannot transact the property)"`
}

func NewOmniSendfreezeCmd(fromaddress string, toaddress string, propertyid int64, amount string) *OmniSendfreezeCmd {
	return &OmniSendfreezeCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
		Amount:      amount,
	}
}

// OmniSendunfreeze // Unfreeze an address for a centrally managed token.
// Note: Only the issuer may unfreeze tokens
// example: $ omnicore-cli "omni_sendunfreeze" "3M9qvHKtgARhqcMtM5cRT9VaiDJ5PSfQGY" "3HTHRxu3aSDV4deakjC7VmsiUp7c6dfbvs" 2 1000
type OmniSendunfreezeCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from (must be issuer of a managed property with freezing enabled"`
	Toaddress   string `json:"toaddress" desc:"the address to unfreeze"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens to unfreeze"`
	Amount      string `json:"amount" desc:"the amount to unfreeze (note: currently unused"`
}

func NewOmniSendunfreezeCmd(fromaddress string, toaddress string, propertyid int64, amount string) *OmniSendunfreezeCmd {
	return &OmniSendunfreezeCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
		Amount:      amount,
	}
}

// OmniSendrawtx // Broadcasts a raw Omni Layer transaction.
// example: $ omnicore-cli "omni_sendrawtx" \     "1MCHESTptvd2LnNp7wmr2sGTpRomteAkq8" "000000000000000100000000017d7840" \     "1EqTta1Rt8ixAA32DuC29oukbsSWU62qAV"
type OmniSendrawtxCmd struct {
	Fromaddress      string  `json:"fromaddress" desc:"the address to send from"`
	Rawtransaction   string  `json:"rawtransaction" desc:"the hex-encoded raw transaction"`
	Referenceaddress *string `json:"referenceaddress" desc:"a reference address (none by default)"`
	Redeemaddress    *string `json:"redeemaddress" desc:"an address that can spend the transaction dust (sender by default)"`
	Referenceamount  *string `json:"referenceamount" desc:"a bitcoin amount that is sent to the receiver (minimal by default)"`
}

func NewOmniSendrawtxCmd(fromaddress string, rawtransaction string, referenceaddress *string, redeemaddress *string, referenceamount *string) *OmniSendrawtxCmd {
	return &OmniSendrawtxCmd{
		Fromaddress:      fromaddress,
		Rawtransaction:   rawtransaction,
		Referenceaddress: referenceaddress,
		Redeemaddress:    redeemaddress,
		Referenceamount:  referenceamount,
	}
}

// OmniFundedSend // Creates and sends a funded simple send transaction.
// All bitcoins from the sender are consumed and if there are bitcoins missing, they are taken from the specified fee source. Change is sent to the fee source!
// example: $ omnicore-cli "omni_funded_send" "1DFa5bT6KMEr6ta29QJouainsjaNBsJQhH" \     "15cWrfuvMxyxGst2FisrQcvcpF48x6sXoH" 1 "100.0" \     "15Jhzz4omEXEyFKbdcccJwuVPea5LqsKM1"
type OmniFundedSendCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from"`
	Toaddress   string `json:"toaddress" desc:"the address of the receiver"`
	Propertyid  int64  `json:"propertyid" desc:"the identifier of the tokens to send"`
	Amount      string `json:"amount" desc:"the amount to send"`
	Feeaddress  string `json:"feeaddress" desc:"the address that is used to pay for fees, if needed"`
}

func NewOmniFundedSendCmd(fromaddress string, toaddress string, propertyid int64, amount string, feeaddress string) *OmniFundedSendCmd {
	return &OmniFundedSendCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Propertyid:  propertyid,
		Amount:      amount,
		Feeaddress:  feeaddress,
	}
}

// OmniFundedSendall // Creates and sends a transaction that transfers all available tokens in the given ecosystem to the recipient.
// All bitcoins from the sender are consumed and if there are bitcoins missing, they are taken from the specified fee source. Change is sent to the fee source!
// example: $ omnicore-cli "omni_funded_sendall" "1DFa5bT6KMEr6ta29QJouainsjaNBsJQhH" \     "15cWrfuvMxyxGst2FisrQcvcpF48x6sXoH" 1 "15Jhzz4omEXEyFKbdcccJwuVPea5LqsKM1"
type OmniFundedSendallCmd struct {
	Fromaddress string `json:"fromaddress" desc:"the address to send from"`
	Toaddress   string `json:"toaddress" desc:"the address of the receiver"`
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem of the tokens to send (1 for main ecosystem, 2 for test ecosystem)"`
	Feeaddress  string `json:"feeaddress" desc:"the address that is used to pay for fees, if needed"`
}

func NewOmniFundedSendallCmd(fromaddress string, toaddress string, ecosystem int64, feeaddress string) *OmniFundedSendallCmd {
	return &OmniFundedSendallCmd{
		Fromaddress: fromaddress,
		Toaddress:   toaddress,
		Ecosystem:   ecosystem,
		Feeaddress:  feeaddress,
	}
}

// OmniGetinfo // Returns various state information of the client and protocol.
// example: $ omnicore-cli "omni_getinfo"
type OmniGetinfoCmd struct {
}

func NewOmniGetinfoCmd() *OmniGetinfoCmd {
	return &OmniGetinfoCmd{}
}

// OmniGetbalance // Returns the token balance for a given address and property.
// example: $ omnicore-cli "omni_getbalance", "1EXoDusjGwvnjZUyKkxZ4UHEf77z6A5S4P" 1
type OmniGetbalanceCmd struct {
	Address    string `json:"address" desc:"the address"`
	Propertyid int64  `json:"propertyid" desc:"the property identifier"`
}

func NewOmniGetbalanceCmd(address string, propertyid int64) *OmniGetbalanceCmd {
	return &OmniGetbalanceCmd{
		Address:    address,
		Propertyid: propertyid,
	}
}

// OmniGetallbalancesforid // Returns a list of token balances for a given currency or property identifier.
// example: $ omnicore-cli "omni_getallbalancesforid" 1
type OmniGetallbalancesforidCmd struct {
	Propertyid int64 `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniGetallbalancesforidCmd() *OmniGetallbalancesforidCmd {
	return &OmniGetallbalancesforidCmd{}
}

// OmniGetallbalancesforaddress // Returns a list of all token balances for a given address.
// example: $ omnicore-cli "omni_getallbalancesforaddress" "1EXoDusjGwvnjZUyKkxZ4UHEf77z6A5S4P"
type OmniGetallbalancesforaddressCmd struct {
	Address string `json:"address" desc:"the address"`
}

func NewOmniGetallbalancesforaddressCmd() *OmniGetallbalancesforaddressCmd {
	return &OmniGetallbalancesforaddressCmd{}
}

// OmniGetwalletbalances // Returns a list of the total token balances of the whole wallet.
// example: $ omnicore-cli "omni_getwalletbalances"
type OmniGetwalletbalancesCmd struct {
}

func NewOmniGetwalletbalancesCmd() *OmniGetwalletbalancesCmd {
	return &OmniGetwalletbalancesCmd{}
}

// OmniGetwalletaddressbalances // Returns a list of all token balances for every wallet address.
// example: $ omnicore-cli "omni_getwalletaddressbalances"
type OmniGetwalletaddressbalancesCmd struct {
}

func NewOmniGetwalletaddressbalancesCmd() *OmniGetwalletaddressbalancesCmd {
	return &OmniGetwalletaddressbalancesCmd{}
}

// OmniGettransaction // Get detailed information about an Omni transaction.
// example: $ omnicore-cli "omni_gettransaction" "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
type OmniGettransactionCmd struct {
	Txid *string `json:"txid" desc:"the hash of the transaction to lookup"`
}

func NewOmniGettransactionCmd() *OmniGettransactionCmd {
	return &OmniGettransactionCmd{}
}

// OmniListtransactions // List wallet transactions, optionally filtered by an address and block boundaries.
// example: $ omnicore-cli "omni_listtransactions"
type OmniListtransactionsCmd struct {
	Txid       *string `json:"txid" desc:"address filter (default: "*")"`
	Count      *int64  `json:"count" desc:"show at most n transactions (default: 10)"`
	Skip       *int64  `json:"skip" desc:"skip the first n transactions (default: 0)"`
	Startblock *int64  `json:"startblock" desc:"first block to begin the search (default: 0)"`
	Endblock   *int64  `json:"endblock" desc:"last block to include in the search (default: 999999999)"`
}

func NewOmniListtransactionsCmd(txid *string, count *int64, skip *int64, startblock *int64, endblock *int64) *OmniListtransactionsCmd {
	return &OmniListtransactionsCmd{
		Txid:       txid,
		Count:      count,
		Skip:       skip,
		Startblock: startblock,
		Endblock:   endblock,
	}
}

// OmniListblocktransactions // Lists all Omni transactions in a block.
// example: $ omnicore-cli "omni_listblocktransactions" 279007
type OmniListblocktransactionsCmd struct {
	Height int64 `json:"height" desc:"specific height to query"`
}

func NewOmniListblocktransactionsCmd() *OmniListblocktransactionsCmd {
	return &OmniListblocktransactionsCmd{}
}

// OmniListpendingtransactions // Returns a list of unconfirmed Omni transactions, pending in the memory pool.
// Note: the validity of pending transactions is uncertain, and the state of the memory pool may change at any moment. It is recommended to check transactions after confirmation, and pending transactions should be considered as invalid.
// example: $ omnicore-cli "omni_listpendingtransactions"
type OmniListpendingtransactionsCmd struct {
	Address *string `json:"address" desc:"filter results by address (default: "" for no filter)"`
}

func NewOmniListpendingtransactionsCmd() *OmniListpendingtransactionsCmd {
	return &OmniListpendingtransactionsCmd{}
}

// OmniGetactivedexsells // Returns currently active offers on the distributed exchange.
// example: $ omnicore-cli "omni_getactivedexsells"
type OmniGetactivedexsellsCmd struct {
}

func NewOmniGetactivedexsellsCmd() *OmniGetactivedexsellsCmd {
	return &OmniGetactivedexsellsCmd{}
}

// OmniListproperties // Lists all tokens or smart properties.
// example: $ omnicore-cli "omni_listproperties"
type OmniListpropertiesCmd struct {
}

func NewOmniListpropertiesCmd() *OmniListpropertiesCmd {
	return &OmniListpropertiesCmd{}
}

// OmniGetproperty // Returns details for about the tokens or smart property to lookup.
// example: $ omnicore-cli "omni_getproperty" 3
type OmniGetpropertyCmd struct {
	Propertyid    int64  `json:"propertyid" desc:"the identifier of the tokens or property"`
	CurrentHeight *int64 `json:"height" desc:"current block height"`
}

func NewOmniGetpropertyCmd() *OmniGetpropertyCmd {
	return &OmniGetpropertyCmd{}
}

// OmniGetactivecrowdsales // Lists currently active crowdsales.
// example: $ omnicore-cli "omni_getactivecrowdsales"
type OmniGetactivecrowdsalesCmd struct {
}

func NewOmniGetactivecrowdsalesCmd() *OmniGetactivecrowdsalesCmd {
	return &OmniGetactivecrowdsalesCmd{}
}

// OmniGetcrowdsale // Returns information about a crowdsale.
// example: $ omnicore-cli "omni_getcrowdsale" 3 true
type OmniGetcrowdsaleCmd struct {
	Propertyid int64 `json:"propertyid" desc:"the identifier of the crowdsale"`
	Verbose    *bool `json:"verbose" desc:"list crowdsale participants (default: false)"`
}

func NewOmniGetcrowdsaleCmd(propertyid int64, verbose *bool) *OmniGetcrowdsaleCmd {
	return &OmniGetcrowdsaleCmd{
		Propertyid: propertyid,
		Verbose:    verbose,
	}
}

// OmniGetgrants // Returns information about granted and revoked units of managed tokens.
// example: $ omnicore-cli "omni_getgrants" 31
type OmniGetgrantsCmd struct {
	PropertyId int64
}

func NewOmniGetgrantsCmd() *OmniGetgrantsCmd {
	return &OmniGetgrantsCmd{}
}

// OmniGetsto // Get information and recipients of a send-to-owners transaction.
// example: $ omnicore-cli "omni_getsto" "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d" "*"
type OmniGetstoCmd struct {
	Txid            string  `json:"txid" desc:"the hash of the transaction to lookup"`
	Recipientfilter *string `json:"recipientfilter" desc:"a filter for recipients (wallet by default, "*" for all)"`
}

func NewOmniGetstoCmd(txid string, recipientfilter *string) *OmniGetstoCmd {
	return &OmniGetstoCmd{
		Txid:            txid,
		Recipientfilter: recipientfilter,
	}
}

// OmniGettrade // Get detailed information and trade matches for orders on the distributed token exchange.
// example: $ omnicore-cli "omni_gettrade" "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
type OmniGettradeCmd struct {
	Txid string `json:"txid" desc:"the hash of the order to lookup"`
}

func NewOmniGettradeCmd() *OmniGettradeCmd {
	return &OmniGettradeCmd{}
}

// OmniGetorderbook // List active offers on the distributed token exchange.
// example: $ omnicore-cli "omni_getorderbook" 2
type OmniGetorderbookCmd struct {
	SalePropertyid    int64  `json:"salePropertyid" desc:"filter orders by propertyid for sale"`
	Desiredpropertyid *int64 `json:"desiredpropertyid" desc:"filter orders by propertyid desired"`
}

func NewOmniGetorderbookCmd(salePropertyid int64, desiredpropertyid *int64) *OmniGetorderbookCmd {
	return &OmniGetorderbookCmd{
		SalePropertyid:    salePropertyid,
		Desiredpropertyid: desiredpropertyid,
	}
}

// OmniGettradehistoryforpair // Retrieves the history of trades on the distributed token exchange for the specified market.
// example: $ omnicore-cli "omni_gettradehistoryforpair" 1 12 500
type OmniGettradehistoryforpairCmd struct {
	FirstPropertyid  int64  `json:"firstPropertyid" desc:"the first side of the traded pair"`
	SecondPropertyid int64  `json:"secondPropertyid" desc:"the second side of the traded pair"`
	Count            *int64 `json:"count" desc:"number of trades to retrieve (default: 10)"`
}

func NewOmniGettradehistoryforpairCmd(firstPropertyid int64, secondPropertyid int64, count *int64) *OmniGettradehistoryforpairCmd {
	return &OmniGettradehistoryforpairCmd{
		FirstPropertyid:  firstPropertyid,
		SecondPropertyid: secondPropertyid,
		Count:            count,
	}
}

// OmniGettradehistoryforaddress // Retrieves the history of orders on the distributed exchange for the supplied address.
// example: $ omnicore-cli "omni_gettradehistoryforaddress" "1MCHESTptvd2LnNp7wmr2sGTpRomteAkq8"
type OmniGettradehistoryforaddressCmd struct {
	Address    string `json:"address" desc:"address to retrieve history for"`
	Count      *int64 `json:"count" desc:"number of orders to retrieve (default: 10)"`
	Propertyid *int64 `json:"propertyid" desc:"filter by propertyid transacted (default: no filter)"`
}

func NewOmniGettradehistoryforaddressCmd(address string, count *int64, propertyid *int64) *OmniGettradehistoryforaddressCmd {
	return &OmniGettradehistoryforaddressCmd{
		Address:    address,
		Count:      count,
		Propertyid: propertyid,
	}
}

// OmniGetactivations // Returns pending and completed feature activations.
// example: $ omnicore-cli "omni_getactivations"
type OmniGetactivationsCmd struct {
}

func NewOmniGetactivationsCmd() *OmniGetactivationsCmd {
	return &OmniGetactivationsCmd{}
}

// OmniGetpayload // Get the payload for an Omni transaction.
// example: $ omnicore-cli "omni_getactivations" "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
type OmniGetpayloadCmd struct {
	TxHash string
}

func NewOmniGetpayloadCmd() *OmniGetpayloadCmd {
	return &OmniGetpayloadCmd{}
}

// OmniGetseedblocks // Returns a list of blocks containing Omni transactions for use in seed block filtering.
// WARNING: The Exodus crowdsale is not stored in LevelDB, thus this is currently only safe to use to generate seed blocks after block 255365.
// example: $ omnicore-cli "omni_getseedblocks" 290000 300000
type OmniGetseedblocksCmd struct {
	Startblock int64 `json:"startblock" desc:"the first block to look for Omni transactions (inclusive)"`
	Endblock   int64 `json:"endblock" desc:"the last block to look for Omni transactions (inclusive)"`
}

func NewOmniGetseedblocksCmd(startblock int64, endblock int64) *OmniGetseedblocksCmd {
	return &OmniGetseedblocksCmd{
		Startblock: startblock,
		Endblock:   endblock,
	}
}

// OmniGetcurrentconsensushash // Returns the consensus hash covering the state of the current block.
// example: $ omnicore-cli "omni_getcurrentconsensushash"
type OmniGetcurrentconsensushashCmd struct {
}

func NewOmniGetcurrentconsensushashCmd() *OmniGetcurrentconsensushashCmd {
	return &OmniGetcurrentconsensushashCmd{}
}

// OmniDecodetransaction // Decodes an Omni transaction.
// If the inputs of the transaction are not in the chain, then they must be provided, because the transaction inputs are used to identify the sender of a transaction.
// A block height can be provided, which is used to determine the parsing rules.
// example: $ omnicore-cli "omni_decodetransaction" "010000000163af14ce6d477e1c793507e32a5b7696288fa89705c0d02a3f66beb3c \     5b8afee0100000000ffffffff02ac020000000000004751210261ea979f6a06f9dafe00fb1263ea0aca959875a7073556a088cdf \     adcd494b3752102a3fd0a8a067e06941e066f78d930bfc47746f097fcd3f7ab27db8ddf37168b6b52ae22020000000000001976a \     914946cb2e08075bcbaf157e47bcb67eb2b2339d24288ac00000000" \     "[{\"txid\":\"eeafb8c5b3be663f2ad0c00597a88f2896765b2ae30735791c7e476dce14af63\",\"vout\":1, \     \"scriptPubKey\":\"76a9149084c0bd89289bc025d0264f7f23148fb683d56c88ac\",\"value\":0.0001123}]"
type OmniDecodetransactionCmd struct {
	Rawtx   string  `json:"rawtx" desc:"the raw transaction to decode"`
	Prevtxs *string `json:"prevtxs" desc:"a JSON array of transaction inputs (default: none)"`
	Height  *int64  `json:"height" desc:"the parsing block height (default: 0 for chain height)"`
}

func NewOmniDecodetransactionCmd(rawtx string, prevtxs *string, height *int64) *OmniDecodetransactionCmd {
	return &OmniDecodetransactionCmd{
		Rawtx:   rawtx,
		Prevtxs: prevtxs,
		Height:  height,
	}
}

// OmniCreaterawtxOpreturn // Adds a payload with class C (op-return) encoding to the transaction.
// If no raw transaction is provided, a new transaction is created.
// If the data encoding fails, then the transaction is not modified.
// example: $ omnicore-cli "omni_createrawtx_opreturn" "01000000000000000000" "00000000000000020000000006dac2c0"
type OmniCreaterawtxOpreturnCmd struct {
	Rawtx   string `json:"rawtx" desc:"the raw transaction to extend (can be null)"`
	Payload string `json:"payload" desc:"the hex-encoded payload to add"`
}

func NewOmniCreaterawtxOpreturnCmd(rawtx string, payload string) *OmniCreaterawtxOpreturnCmd {
	return &OmniCreaterawtxOpreturnCmd{
		Rawtx:   rawtx,
		Payload: payload,
	}
}

// OmniCreaterawtxMultisig // Adds a payload with class B (bare-multisig) encoding to the transaction.
// If no raw transaction is provided, a new transaction is created.
// If the data encoding fails, then the transaction is not modified.
// example: $ omnicore-cli "omni_createrawtx_multisig" \     "0100000001a7a9402ecd77f3c9f745793c9ec805bfa2e14b89877581c734c774864247e6f50400000000ffffffff01aa0a00000 \     00000001976a9146d18edfe073d53f84dd491dae1379f8fb0dfe5d488ac00000000" \     "00000000000000020000000000989680"     "1LifmeXYHeUe2qdKWBGVwfbUCMMrwYtoMm" \     "0252ce4bdd3ce38b4ebbc5a6e1343608230da508ff12d23d85b58c964204c4cef3"
type OmniCreaterawtxMultisigCmd struct {
	Rawtx             string `json:"rawtx" desc:"the raw transaction to extend (can be null)"`
	AddPayload        string `json:"addPayload" desc:"the hex-encoded payload to add"`
	Seed              string `json:"seed" desc:"the seed for obfuscation"`
	RedemptionPayload string `json:"redemptionPayload" desc:"a public key or address for dust redemption"`
}

func NewOmniCreaterawtxMultisigCmd(rawtx string, addPayload string, seed string, redemptionPayload string) *OmniCreaterawtxMultisigCmd {
	return &OmniCreaterawtxMultisigCmd{
		Rawtx:             rawtx,
		AddPayload:        addPayload,
		Seed:              seed,
		RedemptionPayload: redemptionPayload,
	}
}

// OmniCreaterawtxInput // Adds a transaction input to the transaction.
// If no raw transaction is provided, a new transaction is created.
// example: $ omnicore-cli "omni_createrawtx_input" \     "01000000000000000000" "b006729017df05eda586df9ad3f8ccfee5be340aadf88155b784d1fc0e8342ee" 0
type OmniCreaterawtxInputCmd struct {
	Rawtx string `json:"rawtx" desc:"the raw transaction to extend (can be null)"`
	Txid  string `json:"txid" desc:"the hash of the input transaction"`
	N     int64  `json:"n" desc:"the index of the transaction output used as input"`
}

func NewOmniCreaterawtxInputCmd(rawtx string, txid string, n int64) *OmniCreaterawtxInputCmd {
	return &OmniCreaterawtxInputCmd{
		Rawtx: rawtx,
		Txid:  txid,
		N:     n,
	}
}

// OmniCreaterawtxReference // Adds a reference output to the transaction.
// If no raw transaction is provided, a new transaction is created.
// The output value is set to at least the dust threshold.
// example: $ omnicore-cli "omni_createrawtx_reference" \     "0100000001a7a9402ecd77f3c9f745793c9ec805bfa2e14b89877581c734c774864247e6f50400000000ffffffff03aa0a00000     00000001976a9146d18edfe073d53f84dd491dae1379f8fb0dfe5d488ac5c0d0000000000004751210252ce4bdd3ce38b4ebbc5a     6e1343608230da508ff12d23d85b58c964204c4cef3210294cc195fc096f87d0f813a337ae7e5f961b1c8a18f1f8604a909b3a51     21f065b52aeaa0a0000000000001976a914946cb2e08075bcbaf157e47bcb67eb2b2339d24288ac00000000" \     "1CE8bBr1dYZRMnpmyYsFEoexa1YoPz2mfB" \     0.005
type OmniCreaterawtxReferenceCmd struct {
	Rawtx       string `json:"rawtx" desc:"the raw transaction to extend (can be null)"`
	Destination string `json:"destination" desc:"the reference address or destination"`
	Amount      *int64 `json:"amount" desc:"the optional reference amount (minimal by default)"`
}

func NewOmniCreaterawtxReferenceCmd(rawtx string, destination string, amount *int64) *OmniCreaterawtxReferenceCmd {
	return &OmniCreaterawtxReferenceCmd{
		Rawtx:       rawtx,
		Destination: destination,
		Amount:      amount,
	}
}

// OmniCreaterawtxChange // Adds a change output to the transaction.
// The provided inputs are not added to the transaction, but only used to determine the change. It is assumed that the inputs were previously added, for example via `"createrawtransaction"`.
// Optionally a position can be provided, where the change output should be inserted, starting with `0`. If the number of outputs is smaller than the position, then the change output is added to the end. Change outputs should be inserted before reference outputs, and as per default, the change output is added to the`first position.
// If the change amount would be considered as dust, then no change output is added.
// example: $ omnicore-cli "omni_createrawtx_change" \     "0100000001b15ee60431ef57ec682790dec5a3c0d83a0c360633ea8308fbf6d5fc10a779670400000000ffffffff025c0d00000 \     000000047512102f3e471222bb57a7d416c82bf81c627bfcd2bdc47f36e763ae69935bba4601ece21021580b888ff56feb27f17f \     08802ebed26258c23697d6a462d43fc13b565fda2dd52aeaa0a0000000000001976a914946cb2e08075bcbaf157e47bcb67eb2b2 \     339d24288ac00000000" \     "[{\"txid\":\"6779a710fcd5f6fb0883ea3306360c3ad8c0a3c5de902768ec57ef3104e65eb1\",\"vout\":4, \     \"scriptPubKey\":\"76a9147b25205fd98d462880a3e5b0541235831ae959e588ac\",\"value\":0.00068257}]" \     "1CE8bBr1dYZRMnpmyYsFEoexa1YoPz2mfB" 0.000035 1
type OmniCreaterawtxChangeCmd struct {
	Rawtx       string `json:"rawtx" desc:"the raw transaction to extend"`
	Prevtxs     string `json:"prevtxs" desc:"a JSON array of transaction inputs"`
	Destination string `json:"destination" desc:"the destination for the change"`
	Fee         int64  `json:"fee" desc:"the desired transaction fees"`
	Position    *int64 `json:"position" desc:"the position of the change output (default: first position)"`
}

func NewOmniCreaterawtxChangeCmd(rawtx string, prevtxs string, destination string, fee int64, position *int64) *OmniCreaterawtxChangeCmd {
	return &OmniCreaterawtxChangeCmd{
		Rawtx:       rawtx,
		Prevtxs:     prevtxs,
		Destination: destination,
		Fee:         fee,
		Position:    position,
	}
}

// OmniCreatepayloadSimplesend // Create the payload for a simple send transaction.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_simplesend" 1 "100.0"
type OmniCreatepayloadSimplesendCmd struct {
	Propertyid int64  `json:"propertyid" desc:"the identifier of the tokens to send"`
	Amount     string `json:"amount" desc:"the amount to send"`
}

func NewOmniCreatepayloadSimplesendCmd(propertyid int64, amount string) *OmniCreatepayloadSimplesendCmd {
	return &OmniCreatepayloadSimplesendCmd{
		Propertyid: propertyid,
		Amount:     amount,
	}
}

// OmniCreatepayloadSendall // Create the payload for a send all transaction.
// example: $ omnicore-cli "omni_createpayload_sendall" 2
type OmniCreatepayloadSendallCmd struct {
}

func NewOmniCreatepayloadSendallCmd() *OmniCreatepayloadSendallCmd {
	return &OmniCreatepayloadSendallCmd{}
}

// OmniCreatepayloadDexsell // Create a payload to place, update or cancel a sell offer on the traditional distributed OMNI/BTC exchange.
// example: $ omnicore-cli "omni_createpayload_dexsell" 1 "1.5" "0.75" 25 "0.0005" 1
type OmniCreatepayloadDexsellCmd struct {
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale (must be 1 for OMNI or 2 for TOMNI)"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to list for sale"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of bitcoins desired"`
	Paymentwindow     int64  `json:"paymentwindow" desc:"a time limit in blocks a buyer has to pay following a successful accepting order"`
	Minacceptfee      string `json:"minacceptfee" desc:"a minimum mining fee a buyer has to pay to accept the offer"`
	Action            int64  `json:"action" desc:"the action to take (1 for new offers, 2 to update\", 3 to cancel)"`
}

func NewOmniCreatepayloadDexsellCmd(propertyidforsale int64, amountforsale string, amountdesired string, paymentwindow int64, minacceptfee string, action int64) *OmniCreatepayloadDexsellCmd {
	return &OmniCreatepayloadDexsellCmd{
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Amountdesired:     amountdesired,
		Paymentwindow:     paymentwindow,
		Minacceptfee:      minacceptfee,
		Action:            action,
	}
}

// OmniCreatepayloadDexaccept // Create the payload for an accept offer for the specified token and amount.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_dexaccept" 1 "15.0"
type OmniCreatepayloadDexacceptCmd struct {
	Propertyid int64  `json:"propertyid" desc:"the identifier of the token to purchase"`
	Amount     string `json:"amount" desc:"the amount to accept"`
}

func NewOmniCreatepayloadDexacceptCmd(propertyid int64, amount string) *OmniCreatepayloadDexacceptCmd {
	return &OmniCreatepayloadDexacceptCmd{
		Propertyid: propertyid,
		Amount:     amount,
	}
}

// OmniCreatepayloadSto // Creates the payload for a send-to-owners transaction.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_sto" 3 "5000"
type OmniCreatepayloadStoCmd struct {
	Propertyid           int64  `json:"propertyid" desc:"the identifier of the token to distribute"`
	Amount               string `json:"amount" desc:"the amount to distribute"`
	Distributionproperty *int64 `json:"distributionproperty" desc:"the identifier of the property holders to distribute to"`
}

func NewOmniCreatepayloadStoCmd(propertyid int64, amount string, distributionproperty *int64) *OmniCreatepayloadStoCmd {
	return &OmniCreatepayloadStoCmd{
		Propertyid:           propertyid,
		Amount:               amount,
		Distributionproperty: distributionproperty,
	}
}

// OmniCreatepayloadIssuancefixed // Creates the payload for a new tokens issuance with fixed supply.
// example: $ omnicore-cli "omni_createpayload_issuancefixed" 2 1 0 "Companies" "Bitcoin Mining" "Quantum Miner" "" "" "1000000"
type OmniCreatepayloadIssuancefixedCmd struct {
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo        int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid  int64  `json:"previousid" desc:"an identifier of a predecessor token (use 0 for new tokens)"`
	Category    string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name        string `json:"name" desc:"the name of the new tokens to create"`
	Url         string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data        string `json:"data" desc:"a description for the new tokens (can be "")"`
	Amount      string `json:"amount" desc:"the number of tokens to create"`
}

func NewOmniCreatepayloadIssuancefixedCmd(ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string, amount string) *OmniCreatepayloadIssuancefixedCmd {
	return &OmniCreatepayloadIssuancefixedCmd{
		Ecosystem:   ecosystem,
		Typo:        typo,
		Previousid:  previousid,
		Category:    category,
		Subcategory: subcategory,
		Name:        name,
		Url:         url,
		Data:        data,
		Amount:      amount,
	}
}

// OmniCreatepayloadIssuancecrowdsale // Creates the payload for a new tokens issuance with crowdsale.
// example: $ omnicore-cli "omni_createpayload_issuancecrowdsale" 2 1 0 "Companies" "Bitcoin Mining" "Quantum Miner" "" "" 2 "100" 1483228800 30 2
type OmniCreatepayloadIssuancecrowdsaleCmd struct {
	Ecosystem         int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo              int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid        int64  `json:"previousid" desc:"an identifier of a predecessor token (use 0 for new tokens)"`
	Category          string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory       string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name              string `json:"name" desc:"the name of the new tokens to create"`
	Url               string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data              string `json:"data" desc:"a description for the new tokens (can be "")"`
	Propertyiddesired int64  `json:"propertyiddesired" desc:"the identifier of a token eligible to participate in the crowdsale"`
	Tokensperunit     string `json:"tokensperunit" desc:"the amount of tokens granted per unit invested in the crowdsale"`
	Deadline          int64  `json:"deadline" desc:"the deadline of the crowdsale as Unix timestamp"`
	Earlybonus        int64  `json:"earlybonus" desc:"an early bird bonus for participants in percent per week"`
	Issuerpercentage  int64  `json:"issuerpercentage" desc:"a percentage of tokens that will be granted to the issuer"`
}

func NewOmniCreatepayloadIssuancecrowdsaleCmd(ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string, propertyiddesired int64, tokensperunit string, deadline int64, earlybonus int64, issuerpercentage int64) *OmniCreatepayloadIssuancecrowdsaleCmd {
	return &OmniCreatepayloadIssuancecrowdsaleCmd{
		Ecosystem:         ecosystem,
		Typo:              typo,
		Previousid:        previousid,
		Category:          category,
		Subcategory:       subcategory,
		Name:              name,
		Url:               url,
		Data:              data,
		Propertyiddesired: propertyiddesired,
		Tokensperunit:     tokensperunit,
		Deadline:          deadline,
		Earlybonus:        earlybonus,
		Issuerpercentage:  issuerpercentage,
	}
}

// OmniCreatepayloadIssuancemanaged // Creates the payload for a new tokens issuance with manageable supply.
// example: $ omnicore-cli "omni_createpayload_issuancemanaged" 2 1 0 "Companies" "Bitcoin Mining" "Quantum Miner" "" ""
type OmniCreatepayloadIssuancemanagedCmd struct {
	Ecosystem   int64  `json:"ecosystem" desc:"the ecosystem to create the tokens in (1 for main ecosystem, 2 for test ecosystem)"`
	Typo        int64  `json:"type" desc:"the type of the tokens to create: (1 for indivisible tokens, 2 for divisible tokens)"`
	Previousid  int64  `json:"previousid" desc:"an identifier of a predecessor token (use 0 for new tokens)"`
	Category    string `json:"category" desc:"a category for the new tokens (can be "")"`
	Subcategory string `json:"subcategory" desc:"a subcategory for the new tokens (can be "")"`
	Name        string `json:"name" desc:"the name of the new tokens to create"`
	Url         string `json:"url" desc:"an URL for further information about the new tokens (can be "")"`
	Data        string `json:"data" desc:"a description for the new tokens (can be "")"`
}

func NewOmniCreatepayloadIssuancemanagedCmd(ecosystem int64, typo int64, previousid int64, category string, subcategory string, name string, url string, data string) *OmniCreatepayloadIssuancemanagedCmd {
	return &OmniCreatepayloadIssuancemanagedCmd{
		Ecosystem:   ecosystem,
		Typo:        typo,
		Previousid:  previousid,
		Category:    category,
		Subcategory: subcategory,
		Name:        name,
		Url:         url,
		Data:        data,
	}
}

// OmniCreatepayloadClosecrowdsale // Creates the payload to manually close a crowdsale.
// example: $ omnicore-cli "omni_createpayload_closecrowdsale" 70
type OmniCreatepayloadClosecrowdsaleCmd struct {
}

func NewOmniCreatepayloadClosecrowdsaleCmd() *OmniCreatepayloadClosecrowdsaleCmd {
	return &OmniCreatepayloadClosecrowdsaleCmd{}
}

// OmniCreatepayloadGrant // Creates the payload to issue or grant new units of managed tokens.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_grant" 51 "7000"
type OmniCreatepayloadGrantCmd struct {
	Propertyid int64   `json:"propertyid" desc:"the identifier of the tokens to grant"`
	Amount     string  `json:"amount" desc:"the amount of tokens to create"`
	Memo       *string `json:"memo" desc:"a text note attached to this transaction (none by default)"`
}

func NewOmniCreatepayloadGrantCmd(propertyid int64, amount string, memo *string) *OmniCreatepayloadGrantCmd {
	return &OmniCreatepayloadGrantCmd{
		Propertyid: propertyid,
		Amount:     amount,
		Memo:       memo,
	}
}

// OmniCreatepayloadRevoke // Creates the payload to revoke units of managed tokens.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!f
// example: $ omnicore-cli "omni_createpayload_revoke" 51 "100"
type OmniCreatepayloadRevokeCmd struct {
	Propertyid int64   `json:"propertyid" desc:"the identifier of the tokens to revoke"`
	Amount     string  `json:"amount" desc:"the amount of tokens to revoke"`
	Memo       *string `json:"memo" desc:"a text note attached to this transaction (none by default)"`
}

func NewOmniCreatepayloadRevokeCmd(propertyid int64, amount string, memo *string) *OmniCreatepayloadRevokeCmd {
	return &OmniCreatepayloadRevokeCmd{
		Propertyid: propertyid,
		Amount:     amount,
		Memo:       memo,
	}
}

// OmniCreatepayloadChangeissuer // Creates the payload to change the issuer on record of the given tokens.
// example: $ omnicore-cli "omni_createpayload_changeissuer" 3
type OmniCreatepayloadChangeissuerCmd struct {
	Propertyid int64 `json:"propertyid" desc:"the identifier of the tokens to revoke"`
}

func NewOmniCreatepayloadChangeissuerCmd() *OmniCreatepayloadChangeissuerCmd {
	return &OmniCreatepayloadChangeissuerCmd{}
}

// OmniCreatepayloadTrade // Creates the payload to place a trade offer on the distributed token exchange.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_trade" 31 "250.0" 1 "10.0"
type OmniCreatepayloadTradeCmd struct {
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to list for sale"`
	Propertyiddesired int64  `json:"propertyiddesired" desc:"the identifier of the tokens desired in exchange"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of tokens desired in exchange"`
}

func NewOmniCreatepayloadTradeCmd(propertyidforsale int64, amountforsale string, propertyiddesired int64, amountdesired string) *OmniCreatepayloadTradeCmd {
	return &OmniCreatepayloadTradeCmd{
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Propertyiddesired: propertyiddesired,
		Amountdesired:     amountdesired,
	}
}

// OmniCreatepayloadCanceltradesbyprice // Creates the payload to cancel offers on the distributed token exchange with the specified price.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_canceltradesbyprice" 31 "100.0" 1 "5.0"
type OmniCreatepayloadCanceltradesbypriceCmd struct {
	Propertyidforsale int64  `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale"`
	Amountforsale     string `json:"amountforsale" desc:"the amount of tokens to list for sale"`
	Propertyiddesired int64  `json:"propertyiddesired" desc:"the identifier of the tokens desired in exchange"`
	Amountdesired     string `json:"amountdesired" desc:"the amount of tokens desired in exchange"`
}

func NewOmniCreatepayloadCanceltradesbypriceCmd(propertyidforsale int64, amountforsale string, propertyiddesired int64, amountdesired string) *OmniCreatepayloadCanceltradesbypriceCmd {
	return &OmniCreatepayloadCanceltradesbypriceCmd{
		Propertyidforsale: propertyidforsale,
		Amountforsale:     amountforsale,
		Propertyiddesired: propertyiddesired,
		Amountdesired:     amountdesired,
	}
}

// OmniCreatepayloadCanceltradesbypair // Creates the payload to cancel all offers on the distributed token exchange with the given currency pair.
// example: $ omnicore-cli "omni_createpayload_canceltradesbypair" 1 31
type OmniCreatepayloadCanceltradesbypairCmd struct {
	Propertyidforsale int64 `json:"propertyidforsale" desc:"the identifier of the tokens to list for sale"`
	Propertyiddesired int64 `json:"propertyiddesired" desc:"the identifier of the tokens desired in exchange"`
}

func NewOmniCreatepayloadCanceltradesbypairCmd(propertyidforsale int64, propertyiddesired int64) *OmniCreatepayloadCanceltradesbypairCmd {
	return &OmniCreatepayloadCanceltradesbypairCmd{
		Propertyidforsale: propertyidforsale,
		Propertyiddesired: propertyiddesired,
	}
}

// OmniCreatepayloadCancelalltrades // Creates the payload to cancel all offers on the distributed token exchange with the given currency pair.
// example: $ omnicore-cli "omni_createpayload_cancelalltrades" 1
type OmniCreatepayloadCancelalltradesCmd struct {
}

func NewOmniCreatepayloadCancelalltradesCmd() *OmniCreatepayloadCancelalltradesCmd {
	return &OmniCreatepayloadCancelalltradesCmd{}
}

// OmniCreatepayloadEnablefreezing // Creates the payload to enable address freezing for a centrally managed property.
// example: $ omnicore-cli "omni_createpayload_enablefreezing" 3
type OmniCreatepayloadEnablefreezingCmd struct {
	Propertyid int64 `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniCreatepayloadEnablefreezingCmd() *OmniCreatepayloadEnablefreezingCmd {
	return &OmniCreatepayloadEnablefreezingCmd{}
}

// OmniCreatepayloadDisablefreezing // Creates the payload to disable address freezing for a centrally managed property.
// IMPORTANT NOTE:  Disabling freezing for a property will UNFREEZE all frozen addresses for that property!
// example: $ omnicore-cli "omni_createpayload_disablefreezing" 3
type OmniCreatepayloadDisablefreezingCmd struct {
	Propertyid int64 `json:"propertyid" desc:"the identifier of the tokens"`
}

func NewOmniCreatepayloadDisablefreezingCmd() *OmniCreatepayloadDisablefreezingCmd {
	return &OmniCreatepayloadDisablefreezingCmd{}
}

// OmniCreatepayloadFreeze // Creates the payload to freeze an address for a centrally managed token.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_freeze" "3HTHRxu3aSDV4deakjC7VmsiUp7c6dfbvs" 31 "100"
type OmniCreatepayloadFreezeCmd struct {
	Toaddress  string `json:"toaddress" desc:"the address to freeze tokens for"`
	Propertyid int64  `json:"propertyid" desc:"the property to freeze tokens for (must be managed type and have freezing option enabled)"`
	Amount     string `json:"amount" desc:"the amount of tokens to freeze (note: this is unused - once frozen an address cannot send any transactions)"`
}

func NewOmniCreatepayloadFreezeCmd(toaddress string, propertyid int64, amount string) *OmniCreatepayloadFreezeCmd {
	return &OmniCreatepayloadFreezeCmd{
		Toaddress:  toaddress,
		Propertyid: propertyid,
		Amount:     amount,
	}
}

// OmniCreatepayloadUnfreeze // Creates the payload to unfreeze an address for a centrally managed token.
// Note: if the server is not synchronized, amounts are considered as divisible, even if the token may have indivisible units!
// example: $ omnicore-cli "omni_createpayload_unfreeze" "3HTHRxu3aSDV4deakjC7VmsiUp7c6dfbvs" 31 "100"
type OmniCreatepayloadUnfreezeCmd struct {
	Toaddress  string `json:"toaddress" desc:"the address to unfreeze tokens for"`
	Propertyid int64  `json:"propertyid" desc:"the property to unfreeze tokens for (must be managed type and have freezing option enabled)"`
	Amount     string `json:"amount" desc:"the amount of tokens to unfreeze (note: this is unused)"`
}

func NewOmniCreatepayloadUnfreezeCmd(toaddress string, propertyid int64, amount string) *OmniCreatepayloadUnfreezeCmd {
	return &OmniCreatepayloadUnfreezeCmd{
		Toaddress:  toaddress,
		Propertyid: propertyid,
		Amount:     amount,
	}
}

// OmniGetfeecache // Obtains the current amount of fees cached (pending distribution).
// If a property ID is supplied the results will be filtered to show this property ID only.  If no property ID is supplied the results will contain all properties that currently have fees cached pending distribution.
// example: $ omnicore-cli "omni_getfeecache" 31
type OmniGetfeecacheCmd struct {
	PropertyId int64
}

func NewOmniGetfeecacheCmd() *OmniGetfeecacheCmd {
	return &OmniGetfeecacheCmd{}
}

// OmniGetfeetrigger // Obtains the amount at which cached fees will be distributed.
// If a property ID is supplied the results will be filtered to show this property ID only.  If no property ID is supplied the results will contain all properties.
// example: $ omnicore-cli "omni_getfeetrigger" 31
type OmniGetfeetriggerCmd struct {
}

func NewOmniGetfeetriggerCmd() *OmniGetfeetriggerCmd {
	return &OmniGetfeetriggerCmd{}
}

// OmniGetfeeshare // Obtains the current percentage share of fees addresses would receive if a distribution were to occur.
// If an address is supplied the results will be filtered to show this address only.  If no address is supplied the results will be filtered to show wallet addresses only.  If a wildcard is provided (```"*"```) the results will contain all addresses that would receive a share.
// If an ecosystem is supplied the results will reflect the fee share for that ecosystem (main or test).  If no ecosystem is supplied the results will reflect the main ecosystem.
// example: $ omnicore-cli "omni_getfeeshare" "1CE8bBr1dYZRMnpmyYsFEoexa1YoPz2mfB" 1
type OmniGetfeeshareCmd struct {
	Address   *string `json:"address" desc:"the address to filter results on"`
	Ecosystem *int64  `json:"ecosystem" desc:"the ecosystem to obtain the current percentage fee share (1 = main, 2 = test)"`
}

func NewOmniGetfeeshareCmd(address *string, ecosystem *int64) *OmniGetfeeshareCmd {
	return &OmniGetfeeshareCmd{
		Address:   address,
		Ecosystem: ecosystem,
	}
}

// OmniGetfeedistribution // Obtains data for a past distribution of fees.
// A distribution ID must be supplied to identify the distribution to obtain data for.
// example: $ omnicore-cli "omni_getfeedistribution" 1
type OmniGetfeedistributionCmd struct {
	Distributionid int64
}

func NewOmniGetfeedistributionCmd() *OmniGetfeedistributionCmd {
	return &OmniGetfeedistributionCmd{}
}

// OmniGetfeedistributions // Obtains data for past distributions of fees for a property.
// A property ID must be supplied to retrieve past distributions for.
// example: $ omnicore-cli "omni_getfeedistributions" 31
type OmniGetfeedistributionsCmd struct {
	PropertyId int64
}

func NewOmniGetfeedistributionsCmd() *OmniGetfeedistributionsCmd {
	return &OmniGetfeedistributionsCmd{}
}

// OmniSetautocommit // Sets the global flag that determines whether transactions are automatically committed and broadcasted.
// example: $ omnicore-cli "omni_setautocommit" false
type OmniSetautocommitCmd struct {
	AutoCommit bool
}

func NewOmniSetautocommitCmd() *OmniSetautocommitCmd {
	return &OmniSetautocommitCmd{}
}

type OmniProcessTxCmd struct {
	Sender       string
	Reference    string
	TxHash       string
	BlockHash    string
	Block        uint32
	Idx          int
	ScriptEncode string
	Fee          int64
	Time         int64
}

func NewOmniProcessTxCmd() *OmniProcessTxCmd {
	return &OmniProcessTxCmd{}
}

// OmniReadAllTxHashCmd
// Sets the global flag that determines whether transactions are automatically committed and broadcasted.
// example: $ omnicore-cli "OmniReadAllTxHashCmd" false
type OmniReadAllTxHashCmd struct {
}

func NewOmniReadAllTxHashCmd() *OmniReadAllTxHashCmd {
	return &OmniReadAllTxHashCmd{}
}

type OmniPendingAddCmd struct {
	TxId       string
	Sender     string
	MscType    int
	Propertyid uint32
	Amount     string
	Subtract   bool
}

func NewOmniPendingAddCmd() *OmniPendingAddCmd {
	return &OmniPendingAddCmd{}
}

type OmniProcessPaymentCmd struct {
	Sender    string
	Reference string
	TxHash    string
	Amount    int64
	Block     uint32
	Idx       int
}

func NewOmniProcessPaymentCmd() *OmniProcessPaymentCmd {
	return &OmniProcessPaymentCmd{}
}

// OmniRoolBackCmd
// Sets the global flag that determines whether transactions are automatically committed and broadcasted.
// example: $ omnicore-cli "OmniRoolBackCmd" false
type OmniRollBackCmd struct {
	Height uint32
	Hashs  *[]string
}

func NewOmniRollBackCmd(height uint32, hashs *[]string) *OmniRollBackCmd {
	return &OmniRollBackCmd{
		Height: height,
		Hashs:  hashs,
	}
}

// OmniClearCmd
// Sets the global flag that determines whether transactions are automatically committed and broadcasted.
// example: $ omnicore-cli "OmniClearCmd" false
type OmniClearCmd struct {}

func NewOmniClearCmd(height uint32, hashs  *[]string) *OmniClearCmd {
	return &OmniClearCmd{}
}

// TXExodusFundraiserCmd
// Sets the global flag that determines whether transactions are automatically committed and broadcasted.
// example: $ omnicore-cli "TXExodusFundraiserCmd" false
type OmniTXExodusFundraiserCmd struct{
	Hash string
	StrSender  string
	NBlock int
	AmountInvested int64
	NTime int32
}

func NewTXExodusFundraiserCmd() *OmniTXExodusFundraiserCmd {
	return &OmniTXExodusFundraiserCmd{}
}

func init() {
	// The commands in this file are only usable with a wallet server.
	flags := OMiniOnly
	MustRegisterCmd("omni_send", (*OmniSendCmd)(nil), flags)
	MustRegisterCmd("omni_senddexsell", (*OmniSenddexsellCmd)(nil), flags)
	MustRegisterCmd("omni_senddexaccept", (*OmniSenddexacceptCmd)(nil), flags)
	MustRegisterCmd("omni_sendissuancecrowdsale", (*OmniSendissuancecrowdsaleCmd)(nil), flags)
	MustRegisterCmd("omni_sendissuancefixed", (*OmniSendissuancefixedCmd)(nil), flags)
	MustRegisterCmd("omni_sendissuancemanaged", (*OmniSendissuancemanagedCmd)(nil), flags)
	MustRegisterCmd("omni_sendsto", (*OmniSendstoCmd)(nil), flags)
	MustRegisterCmd("omni_sendgrant", (*OmniSendgrantCmd)(nil), flags)
	MustRegisterCmd("omni_sendrevoke", (*OmniSendrevokeCmd)(nil), flags)
	MustRegisterCmd("omni_sendclosecrowdsale", (*OmniSendclosecrowdsaleCmd)(nil), flags)
	MustRegisterCmd("omni_sendtrade", (*OmniSendtradeCmd)(nil), flags)
	MustRegisterCmd("omni_sendcanceltradesbyprice", (*OmniSendcanceltradesbypriceCmd)(nil), flags)
	MustRegisterCmd("omni_sendcanceltradesbypair", (*OmniSendcanceltradesbypairCmd)(nil), flags)
	MustRegisterCmd("omni_sendcancelalltrades", (*OmniSendcancelalltradesCmd)(nil), flags)
	MustRegisterCmd("omni_sendchangeissuer", (*OmniSendchangeissuerCmd)(nil), flags)
	MustRegisterCmd("omni_sendall", (*OmniSendallCmd)(nil), flags)
	MustRegisterCmd("omni_sendenablefreezing", (*OmniSendenablefreezingCmd)(nil), flags)
	MustRegisterCmd("omni_senddisablefreezing", (*OmniSenddisablefreezingCmd)(nil), flags)
	MustRegisterCmd("omni_sendfreeze", (*OmniSendfreezeCmd)(nil), flags)
	MustRegisterCmd("omni_sendunfreeze", (*OmniSendunfreezeCmd)(nil), flags)
	MustRegisterCmd("omni_sendrawtx", (*OmniSendrawtxCmd)(nil), flags)
	MustRegisterCmd("omni_funded_send", (*OmniFundedSendCmd)(nil), flags)
	MustRegisterCmd("omni_funded_sendall", (*OmniFundedSendallCmd)(nil), flags)
	MustRegisterCmd("omni_getinfo", (*OmniGetinfoCmd)(nil), flags)
	MustRegisterCmd("omni_getbalance", (*OmniGetbalanceCmd)(nil), flags)
	MustRegisterCmd("omni_getallbalancesforid", (*OmniGetallbalancesforidCmd)(nil), flags)
	MustRegisterCmd("omni_getallbalancesforaddress", (*OmniGetallbalancesforaddressCmd)(nil), flags)
	MustRegisterCmd("omni_getwalletbalances", (*OmniGetwalletbalancesCmd)(nil), flags)
	MustRegisterCmd("omni_getwalletaddressbalances", (*OmniGetwalletaddressbalancesCmd)(nil), flags)
	MustRegisterCmd("omni_gettransaction", (*OmniGettransactionCmd)(nil), flags)
	MustRegisterCmd("omni_listtransactions", (*OmniListtransactionsCmd)(nil), flags)
	MustRegisterCmd("omni_listblocktransactions", (*OmniListblocktransactionsCmd)(nil), flags)
	MustRegisterCmd("omni_listpendingtransactions", (*OmniListpendingtransactionsCmd)(nil), flags)
	MustRegisterCmd("omni_getactivedexsells", (*OmniGetactivedexsellsCmd)(nil), flags)
	MustRegisterCmd("omni_listproperties", (*OmniListpropertiesCmd)(nil), flags)
	MustRegisterCmd("omni_getproperty", (*OmniGetpropertyCmd)(nil), flags)
	MustRegisterCmd("omni_getactivecrowdsales", (*OmniGetactivecrowdsalesCmd)(nil), flags)
	MustRegisterCmd("omni_getcrowdsale", (*OmniGetcrowdsaleCmd)(nil), flags)
	MustRegisterCmd("omni_getgrants", (*OmniGetgrantsCmd)(nil), flags)
	MustRegisterCmd("omni_getsto", (*OmniGetstoCmd)(nil), flags)
	MustRegisterCmd("omni_gettrade", (*OmniGettradeCmd)(nil), flags)
	MustRegisterCmd("omni_getorderbook", (*OmniGetorderbookCmd)(nil), flags)
	MustRegisterCmd("omni_gettradehistoryforpair", (*OmniGettradehistoryforpairCmd)(nil), flags)
	MustRegisterCmd("omni_gettradehistoryforaddress", (*OmniGettradehistoryforaddressCmd)(nil), flags)
	MustRegisterCmd("omni_getactivations", (*OmniGetactivationsCmd)(nil), flags)
	MustRegisterCmd("omni_getpayload", (*OmniGetpayloadCmd)(nil), flags)
	MustRegisterCmd("omni_getseedblocks", (*OmniGetseedblocksCmd)(nil), flags)
	MustRegisterCmd("omni_getcurrentconsensushash", (*OmniGetcurrentconsensushashCmd)(nil), flags)
	MustRegisterCmd("omni_decodetransaction", (*OmniDecodetransactionCmd)(nil), flags)
	MustRegisterCmd("omni_createrawtx_opreturn", (*OmniCreaterawtxOpreturnCmd)(nil), flags)
	MustRegisterCmd("omni_createrawtx_multisig", (*OmniCreaterawtxMultisigCmd)(nil), flags)
	MustRegisterCmd("omni_createrawtx_input", (*OmniCreaterawtxInputCmd)(nil), flags)
	MustRegisterCmd("omni_createrawtx_reference", (*OmniCreaterawtxReferenceCmd)(nil), flags)
	MustRegisterCmd("omni_createrawtx_change", (*OmniCreaterawtxChangeCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_simplesend", (*OmniCreatepayloadSimplesendCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_sendall", (*OmniCreatepayloadSendallCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_dexsell", (*OmniCreatepayloadDexsellCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_dexaccept", (*OmniCreatepayloadDexacceptCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_sto", (*OmniCreatepayloadStoCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_issuancefixed", (*OmniCreatepayloadIssuancefixedCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_issuancecrowdsale", (*OmniCreatepayloadIssuancecrowdsaleCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_issuancemanaged", (*OmniCreatepayloadIssuancemanagedCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_closecrowdsale", (*OmniCreatepayloadClosecrowdsaleCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_grant", (*OmniCreatepayloadGrantCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_revoke", (*OmniCreatepayloadRevokeCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_changeissuer", (*OmniCreatepayloadChangeissuerCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_trade", (*OmniCreatepayloadTradeCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_canceltradesbyprice", (*OmniCreatepayloadCanceltradesbypriceCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_canceltradesbypair", (*OmniCreatepayloadCanceltradesbypairCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_cancelalltrades", (*OmniCreatepayloadCancelalltradesCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_enablefreezing", (*OmniCreatepayloadEnablefreezingCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_disablefreezing", (*OmniCreatepayloadDisablefreezingCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_freeze", (*OmniCreatepayloadFreezeCmd)(nil), flags)
	MustRegisterCmd("omni_createpayload_unfreeze", (*OmniCreatepayloadUnfreezeCmd)(nil), flags)
	MustRegisterCmd("omni_getfeecache", (*OmniGetfeecacheCmd)(nil), flags)
	MustRegisterCmd("omni_getfeetrigger", (*OmniGetfeetriggerCmd)(nil), flags)
	MustRegisterCmd("omni_getfeeshare", (*OmniGetfeeshareCmd)(nil), flags)
	MustRegisterCmd("omni_getfeedistribution", (*OmniGetfeedistributionCmd)(nil), flags)
	MustRegisterCmd("omni_getfeedistributions", (*OmniGetfeedistributionsCmd)(nil), flags)
	MustRegisterCmd("omni_setautocommit", (*OmniSetautocommitCmd)(nil), flags)
	MustRegisterCmd("omni_processtx", (*OmniProcessTxCmd)(nil), flags)
	MustRegisterCmd("omni_readalltxhash", (*OmniReadAllTxHashCmd)(nil), flags)
	MustRegisterCmd("omni_rollback", (*OmniRollBackCmd)(nil), flags)
	MustRegisterCmd("omni_pending_add", (*OmniPendingAddCmd)(nil), flags)
	MustRegisterCmd("omni_processpayment", (*OmniProcessPaymentCmd)(nil), flags)
	MustRegisterCmd("omni_clear", (*OmniClearCmd)(nil), flags)
	MustRegisterCmd("omni_txexodus_fundraiser", (*OmniTXExodusFundraiserCmd)(nil), flags)
}
