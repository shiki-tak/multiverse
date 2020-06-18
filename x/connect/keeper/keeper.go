package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/shiki-tak/connect/x/connect/types"
	nftexported "github.com/shiki-tak/connect/x/nft/exported"

	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
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

func (k Keeper) SendTransfer(
	ctx sdk.Context,
	srcPort string,
	srcChannel string,
	destHeight uint64,
	receiver sdk.AccAddress,
	sender sdk.AccAddress,
	denom string,
	nft nftexported.NFT,
) error {
	sourceChannelEnd, found := k.ChannelKeeper.GetChannel(ctx, srcPort, srcChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, srcChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.ChannelKeeper.GetNextSequenceSend(ctx, srcPort, srcChannel)
	if !found {
		return channeltypes.ErrSequenceSendNotFound
	}

	return k.createOutgoingPacket(ctx, sequence, srcPort, srcChannel, destinationPort, destinationChannel, destHeight, denom, nft, sender, receiver)
}

func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
	destHeight uint64,
	denom string,
	nft nftexported.NFT,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
) error {
	channelCap, ok := k.ScopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	if err := k.NFTKeeper.DeleteNFT(ctx, denom, nft.GetID()); err != nil {
		return err
	}

	packetData := types.NewNonFungibleTokenPacketData(
		denom, nft.GetID(), nft.GetTokenURI(), sender, receiver,
	)

	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeoutHeight,
		DefaultPacketTimeoutTimestamp,
	)

	return k.ChannelKeeper.SendPacket(ctx, channelCap, packet)
}
