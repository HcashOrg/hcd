// Copyright (c) 2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcjson

// GetMultisigOutInfoResult models the data returned from the getmultisigoutinfo
// command.
type GetMultisigOutInfoResult struct {
	Address      string   `json:"address"`
	RedeemScript string   `json:"redeemscript"`
	M            uint8    `json:"m"`
	N            uint8    `json:"n"`
	Pubkeys      []string `json:"pubkeys"`
	TxHash       string   `json:"txhash"`
	BlockHeight  uint32   `json:"blockheight"`
	BlockHash    string   `json:"blockhash"`
	Spent        bool     `json:"spent"`
	SpentBy      string   `json:"spentby"`
	SpentByIndex uint32   `json:"spentbyindex"`
	Amount       float64  `json:"amount"`
}

// GetStakeInfoResult models the data returned from the getstakeinfo
// command.
type GetStakeInfoResult struct {
	BlockHeight      int64   `json:"blockheight"`
	PoolSize         uint32  `json:"poolsize"`
	Difficulty       float64 `json:"difficulty"`
	AllMempoolTix    uint32  `json:"allmempooltix"`
	OwnMempoolTix    uint32  `json:"ownmempooltix"`
	Immature         uint32  `json:"immature"`
	Live             uint32  `json:"live"`
	ProportionLive   float64 `json:"proportionlive"`
	Voted            uint32  `json:"voted"`
	TotalSubsidy     float64 `json:"totalsubsidy"`
	Missed           uint32  `json:"missed"`
	ProportionMissed float64 `json:"proportionmissed"`
	Revoked          uint32  `json:"revoked"`
	Expired          uint32  `json:"expired"`

	AiPoolSize         uint32  `json:"aipoolsize"`
	AiDifficulty       float64 `json:"aidifficulty"`
	AiAllMempoolTix    uint32  `json:"aiallmempooltix"`
	AiOwnMempoolTix    uint32  `json:"aiownmempooltix"`
	AiImmature         uint32  `json:"aiimmature"`
	AiLive             uint32  `json:"ailive"`
	AiProportionLive   float64 `json:"aiproportionlive"`
	AiVoted            uint32  `json:"aivoted"`
	AiMissed           uint32  `json:"aimissed"`
	AiProportionMissed float64 `json:"aiproportionmissed"`
	AiRevoked          uint32  `json:"airevoked"`
	AiExpired          uint32  `json:"aiexpired"`
}

// GetTicketsResult models the data returned from the gettickets
// command.
type GetTicketsResult struct {
	Hashes []string `json:"hashes"`
}

type GetAiTicketsResult struct {
	Hashes []string `json:"hashes"`
}

// VoteChoice models the data for a vote choice in the getvotechoices result.
type VoteChoice struct {
	AgendaID          string `json:"agendaid"`
	AgendaDescription string `json:"agendadescription"`
	ChoiceID          string `json:"choiceid"`
	ChoiceDescription string `json:"choicedescription"`
}

// GetVoteChoicesResult models the data returned by the getvotechoices command.
type GetVoteChoicesResult struct {
	Version uint32       `json:"version"`
	Choices []VoteChoice `json:"choices"`
}

// ScriptInfo is the structure representing a redeem script, its hash,
// and its address.
type ScriptInfo struct {
	Hash160      string `json:"hash160"`
	Address      string `json:"address"`
	RedeemScript string `json:"redeemscript"`
}

// ListScriptsResult models the data returned from the listscripts
// command.
type ListScriptsResult struct {
	Scripts []ScriptInfo `json:"scripts"`
}

// RedeemMultiSigOutResult models the data returned from the redeemmultisigout
// command.
type RedeemMultiSigOutResult struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
	Errors   []SignRawTransactionError `json:"errors,omitempty"`
}

// RedeemMultiSigOutsResult models the data returned from the redeemmultisigouts
// command.
type RedeemMultiSigOutsResult struct {
	Results []RedeemMultiSigOutResult `json:"results"`
}

// GetStraightPubKeyResult models the data returned from the getStraightPubKey
// command.
type GetStraightPubKeyResult struct {
	StraightPubKey string `json:"StraightPubKey"`
}

// SendToMultiSigResult models the data returned from the sendtomultisig
// command.
type SendToMultiSigResult struct {
	TxHash       string `json:"txhash"`
	Address      string `json:"address"`
	RedeemScript string `json:"redeemscript"`
}

// SignedTransaction is a signed transaction resulting from a signrawtransactions
// command.
type SignedTransaction struct {
	SigningResult SignRawTransactionResult `json:"signingresult"`
	Sent          bool                     `json:"sent"`
	TxHash        *string                  `json:"txhash,omitempty"`
}

// SignRawTransactionsResult models the data returned from the signrawtransactions
// command.
type SignRawTransactionsResult struct {
	Results []SignedTransaction `json:"results"`
}

// PoolUserTicket is the JSON struct corresponding to a stake pool user ticket
// object.
type PoolUserTicket struct {
	Status        string `json:"status"`
	Ticket        string `json:"ticket"`
	TicketHeight  uint32 `json:"ticketheight"`
	SpentBy       string `json:"spentby"`
	SpentByHeight uint32 `json:"spentbyheight"`
}

// StakePoolUserInfoResult models the data returned from the stakepooluserinfo
// command.
type StakePoolUserInfoResult struct {
	Tickets        []PoolUserTicket `json:"tickets"`
	InvalidTickets []string         `json:"invalid"`
}

// WalletInfoResult models the data returned from the walletinfo
// command.
type WalletInfoResult struct {
	DaemonConnected  bool    `json:"daemonconnected"`
	Unlocked         bool    `json:"unlocked"`
	TxFee            float64 `json:"txfee"`
	TicketFee        float64 `json:"ticketfee"`
	TicketPurchasing bool    `json:"ticketpurchasing"`
	VoteBits         uint16  `json:"votebits"`
	VoteBitsExtended string  `json:"votebitsextended"`
	VoteVersion      uint32  `json:"voteversion"`
	Voting           bool    `json:"voting"`
}
