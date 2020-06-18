package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibc20transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/shiki-tak/connect/x/connect/types"
	nfttypes "github.com/shiki-tak/connect/x/nft/types"
)

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NonFungibleTokenPacketData) error {
	destPrefix := ibc20transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	srcPrefix := ibc20transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())

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
