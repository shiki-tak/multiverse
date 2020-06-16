package connect

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
		// TODO: Define your msg cases
		//
		//Example:
		case MsgTransfer:
			return handleMsgTransfer(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handle<Action> does x
func handleMsgTransfer(ctx sdk.Context, k Keeper, msg MsgTransfer) (*sdk.Result, error) {
	nft, err := k.NFTKeeper.GetNFT(ctx, msg.Denom, msg.TokenID)
	if err != nil {
		return nil, err
	}

	err = k.SendTransfer(ctx, msg.SrcPort, msg.SrcChannel, msg.DestHeight, msg.Receiver, msg.Sender, nft)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}
