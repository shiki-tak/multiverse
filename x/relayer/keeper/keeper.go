package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/shiki-tak/connect/x/relayer/types"
)

const (
	DefaultPacketTimeoutHeight = 1000 // NOTE: in blocks

	DefaultPacketTimeoutTimestamp = 0 // NOTE: in nanoseconds
)

// Keeper of the ibc store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Marshaler
	NFTKeeper     types.NFTKeeper
	ChannelKeeper types.ChannelKeeper
	PortKeeper    types.PortKeeper
	ScopedKeeper  capability.ScopedKeeper
}

// NewKeeper creates a ibc keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey,
	nftKeeper types.NFTKeeper,
	channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper, scopedKeeper capability.ScopedKeeper) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		NFTKeeper:     nftKeeper,
		ChannelKeeper: channelKeeper,
		PortKeeper:    portKeeper,
		ScopedKeeper:  scopedKeeper,
	}
	return keeper
}
