package cli

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/shiki-tak/connect/x/relayer/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	relayerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	relayerTxCmd.AddCommand(flags.PostCommands(
		GetSendPacketCmd(cdc),
	)...)

	return relayerTxCmd
}

// GetTransferTxCmd returns the command to create a NewMsgTransfer transaction
func GetSendPacketCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-packet [src-port] [src-channel] [dest-height] [receiver]",
		Short: "Transfer non-fungible token through IBC.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)
			sender := cliCtx.GetFromAddress()

			srcPort := args[0]
			srcChannel := args[1]
			destHeight, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			receiver, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgSendPacket(srcPort, srcChannel, uint64(destHeight), receiver, sender)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
