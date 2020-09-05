package cosmwasm

import (
	wasmtypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shiki-tak/multiverse/x/wasm/internal/types"
)

type Env struct {
	Block    wasmtypes.BlockInfo    `json:"block"`
	Message  wasmtypes.MessageInfo  `json:"message"`
	Contract wasmtypes.ContractInfo `json:"contract"`
}

func NewEnv(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAddr sdk.AccAddress) Env {
	// safety checks before casting below
	if ctx.BlockHeight() < 0 {
		panic("Block height must never be negative")
	}
	if ctx.BlockTime().Unix() < 0 {
		panic("Block (unix) time must never be negative ")
	}
	return Env{
		Block: wasmtypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Message: wasmtypes.MessageInfo{
			Sender:    wasmtypes.HumanAddress(creator),
			SentFunds: types.NewWasmCoins(deposit),
		},
		Contract: wasmtypes.ContractInfo{
			Address: wasmtypes.HumanAddress(contractAddr),
		},
	}
}
