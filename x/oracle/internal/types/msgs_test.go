package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func TestMsgSort(t *testing.T) {
	from := sdk.AccAddress([]byte("someName"))
	price, _ := sdk.NewDecFromStr("0.01155578")
	expiry := time.Now()

	msg := NewMsgPostPrice(from, "uftm", price, expiry)

	fee := auth.NewStdFee(200000, nil)
	stdTx := auth.NewStdTx([]sdk.Msg{msg}, fee, []auth.StdSignature{}, "")
	signBytes := auth.StdSignBytes("xar-chain-dora", 4, 1, stdTx.Fee, stdTx.Msgs, stdTx.Memo)

	t.Logf("%s", signBytes)
	signed := auth.StdSignBytes(
		"xar-chain-dora", 4, 1, auth.NewStdFee(200000, nil), []sdk.Msg{msg}, "",
	)
	t.Logf("%s", signed)
}