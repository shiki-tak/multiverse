package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/shiki-tak/connect/x/connect/types"

	nftexported "github.com/shiki-tak/connect/x/nft/exported"
	nfttypes "github.com/shiki-tak/connect/x/nft/types"

	ibc20transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
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

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	destPrefix := ibc20transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	srcPrefix := ibc20transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())

	destForID := strings.HasPrefix(data.ID, destPrefix)
	destForDenom := strings.HasPrefix(data.Denom, destPrefix)

	srcForID := strings.HasPrefix(data.ID, srcPrefix)
	srcForDenom := strings.HasPrefix(data.Denom, srcPrefix)

	var nft nfttypes.BaseNFT
	denom := data.Denom

	if srcForID && srcForDenom {
		tokenID := data.ID[len(destPrefix)+1:]
		denom = data.Denom[len(destPrefix)+1:]

		nft = nfttypes.NewBaseNFT(tokenID, data.Receiver, data.TokenURI)
	} else {
		tokenID := types.GenerateKey(destPrefix, data.ID)
		denom = types.GenerateKey(destPrefix, data.Denom)

		nft = nfttypes.NewBaseNFT(tokenID, data.Receiver, data.TokenURI)
	}

	err := k.NFTKeeper.MintNFT(ctx, denom, &nft)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData, ack types.NonFungibleTokenPacketAcknowledgement) error {
	if !ack.Success {
		return k.refundPacketAmount(ctx, packet, data)
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	return k.refundPacketAmount(ctx, packet, data)
}

func (k Keeper) refundPacketAmount(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	return k.OnRecvPacket(ctx, packet, data)
}

func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	chanCap, ok := k.ScopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
	}
	return k.ChannelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := ibctypes.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.ScopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channel.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.ScopedKeeper.GetCapability(ctx, ibctypes.PortPath(portID))
	return ok
}

func (k Keeper) BindPort(ctx sdk.Context, portID string) (*capability.Capability, error) {
	cap := k.PortKeeper.BindPort(ctx, portID)
	if err := k.ClaimCapability(ctx, cap, ibctypes.PortPath(portID)); err != nil {
		return nil, err
	}
	return cap, nil
}

func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.ScopedKeeper.ClaimCapability(ctx, cap, name)
}
