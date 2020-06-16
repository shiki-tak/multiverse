package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"

	"github.com/shiki-tak/connect/x/connect/types"

	nftexported "github.com/shiki-tak/connect/x/nft/exported"
)

// Keeper of the ibc store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	NFTKeeper     types.NFTKeeper
	ChannelKeeper types.ChannelKeeper
	PortKeeper    types.PortKeeper
	ScopedKeeper  capability.ScopedKeeper
}

// NewKeeper creates a ibc keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey,
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

/*
	SrcPort:    srcPort,
	SrcChannel: srcChannel,
	DestHeight: destHeight,
	Receiver:   receiver,
	Sender:     sender,
	Denom:      denom,
	TokenID:    tokenID,
*/
func (k Keeper) SendTransfer(ctx sdk.Context, srcPort string, srcChannel string,
	destHeight int, receiver sdk.AccAddress, sender sdk.AccAddress, nft nftexported.NFT) error {
	return nil
}

// Get returns the pubkey from the adddress-pubkey relation
func (k Keeper) Get(ctx sdk.Context, key string) (string, error) {
	return "", nil
}

func (k Keeper) Set(ctx sdk.Context, key string, value string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}

func (k Keeper) Delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}
