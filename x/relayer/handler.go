package relayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the ibc type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgLinkChains:
			return handleMsgLinkChains(ctx, k, msg)
		case MsgSendPacket:
			return handleMsgSendPacket(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgLinkChains(ctx sdk.Context, k Keeper, msg MsgLinkChains) (*sdk.Result, error) {
	return nil, nil
}

func handleMsgSendPacket(ctx sdk.Context, k Keeper, msg MsgSendPacket) (*sdk.Result, error) {
	return nil, nil
}
