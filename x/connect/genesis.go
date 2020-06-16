package connect

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shiki-tak/connect/x/connect/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	_, err := k.BindPort(ctx, types.PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	return NewGenesisState()
}
