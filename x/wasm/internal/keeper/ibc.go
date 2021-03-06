package keeper

import (
	"strings"

	wasm "github.com/CosmWasm/go-cosmwasm"
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/shiki-tak/multiverse/x/wasm/internal/keeper/cosmwasm"
	"github.com/shiki-tak/multiverse/x/wasm/internal/types"
)

// bindIbcPort will reserve the port.
// returns a string name of the port or error if we cannot bind it.
// this will fail if call twice.
func (k Keeper) bindIbcPort(ctx sdk.Context, portID string) error {
	// TODO: always set up IBC in tests, so we don't need to disable this
	if k.PortKeeper == nil {
		return nil
	}
	cap := k.PortKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// ensureIbcPort is like registerIbcPort, but it checks if we already hold the port
// before calling register, so this is safe to call multiple times.
// Returns success if we already registered or just registered and error if we cannot
// (lack of permissions or someone else has it)
func (k Keeper) ensureIbcPort(ctx sdk.Context, contractAddr sdk.AccAddress) (string, error) {
	// TODO: always set up IBC in tests, so we don't need to disable this
	if k.PortKeeper == nil {
		return PortIDForContract(contractAddr), nil
	}

	portID := PortIDForContract(contractAddr)
	if _, ok := k.ScopedKeeper.GetCapability(ctx, host.PortPath(portID)); ok {
		return portID, nil
	}
	return portID, k.bindIbcPort(ctx, portID)
}

const portIDPrefix = "wasm."

func PortIDForContract(addr sdk.AccAddress) string {
	return portIDPrefix + addr.String()
}

func ContractFromPortID(portID string) (sdk.AccAddress, error) {
	if !strings.HasPrefix(portID, portIDPrefix) {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "without prefix")
	}
	return sdk.AccAddressFromBech32(portID[len(portIDPrefix):])
}

// ClaimCapability allows the transfer module to claim a capability
//that IBC module passes to it
// TODO: make private and inline??
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.ScopedKeeper.ClaimCapability(ctx, cap, name)
}

// IBCContractCallbacks defines the methods for go-cosmwasm to interact with the wasm contract.
// A mock contract would implement the interface to fully simulate a wasm contract's behaviour.
type IBCContractCallbacks interface {
	// Package livecycle

	// OnIBCPacketReceive handles an incoming IBC package
	OnIBCPacketReceive(hash []byte, params cosmwasm.Env, packet cosmwasm.IBCPacket, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCPacketReceiveResponse, uint64, error)
	// OnIBCPacketAcknowledgement handles a IBC package execution on the counterparty chain
	OnIBCPacketAcknowledgement(hash []byte, params cosmwasm.Env, packetAck cosmwasm.IBCAcknowledgement, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCPacketAcknowledgementResponse, uint64, error)
	// OnIBCPacketTimeout reverts state when the IBC package execution does not come in time
	OnIBCPacketTimeout(hash []byte, params cosmwasm.Env, packet cosmwasm.IBCPacket, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCPacketTimeoutResponse, uint64, error)
	// channel livecycle

	// OnIBCChannelOpen does the protocol version negotiation during channel handshake phase
	OnIBCChannelOpen(hash []byte, params cosmwasm.Env, channel cosmwasm.IBCChannel, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCChannelOpenResponse, uint64, error)
	// OnIBCChannelConnect callback when a IBC channel is established
	OnIBCChannelConnect(hash []byte, params cosmwasm.Env, channel cosmwasm.IBCChannel, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCChannelConnectResponse, uint64, error)
	// OnIBCChannelConnect callback when a IBC channel is closed
	OnIBCChannelClose(ctx sdk.Context, hash []byte, params cosmwasm.Env, channel cosmwasm.IBCChannel, meter sdk.GasMeter, gas uint64) (*cosmwasm.IBCChannelCloseResponse, uint64, error)
	Execute(hash []byte, params wasmTypes.Env, msg []byte, store prefix.Store, api wasm.GoAPI, querier QueryHandler, meter sdk.GasMeter, gas uint64) (*cosmwasm.HandleResponse, uint64, error)
}

var MockContracts = make(map[string]IBCContractCallbacks, 0)

func (k Keeper) OnOpenChannel(ctx sdk.Context, contractAddr sdk.AccAddress, channel cosmwasm.IBCChannel) error {
	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return err
	}

	var sender sdk.AccAddress // we don't know the sender
	params := cosmwasm.NewEnv(ctx, sender, nil, contractAddr)

	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	mock, ok := MockContracts[contractAddr.String()]
	if !ok { // hack for testing without wasmer
		panic("not supported")
	}
	res, gasUsed, execErr := mock.OnIBCChannelOpen(codeInfo.CodeHash, params, channel, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}
	if !res.Success { // todo: would it make more sense to let the contract return an error instead?
		return sdkerrors.Wrap(types.ErrInvalid, res.Reason)
	}
	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, contractAddr sdk.AccAddress, packet cosmwasm.IBCPacket) ([]byte, error) {
	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	var sender sdk.AccAddress // we don't know the sender
	params := cosmwasm.NewEnv(ctx, sender, nil, contractAddr)

	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	mock, ok := MockContracts[contractAddr.String()]
	if !ok { // hack for testing without wasmer
		panic("not supported")
	}
	res, gasUsed, execErr := mock.OnIBCPacketReceive(codeInfo.CodeHash, params, packet, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddr)
	ctx.EventManager().EmitEvents(events)

	if err := k.messenger.DispatchV2(ctx, contractAddr, packet.Destination, res.Messages...); err != nil {
		return nil, err
	}
	return res.Acknowledgement, nil
}

func (k Keeper) OnAckPacket(ctx sdk.Context, contractAddr sdk.AccAddress, acknowledgement cosmwasm.IBCAcknowledgement) error {
	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return err
	}

	var sender sdk.AccAddress // we don't know the sender
	params := cosmwasm.NewEnv(ctx, sender, nil, contractAddr)

	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	mock, ok := MockContracts[contractAddr.String()] // hack for testing without wasmer
	if !ok {
		panic("not supported")
	}
	res, gasUsed, execErr := mock.OnIBCPacketAcknowledgement(codeInfo.CodeHash, params, acknowledgement, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddr)
	ctx.EventManager().EmitEvents(events)

	if err := k.messenger.DispatchV2(ctx, contractAddr, acknowledgement.OriginalPacket.Source, res.Messages...); err != nil {
		return err
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, contractAddr sdk.AccAddress, packet cosmwasm.IBCPacket) error {
	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return err
	}

	var sender sdk.AccAddress // we don't know the sender
	params := cosmwasm.NewEnv(ctx, sender, nil, contractAddr)

	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	mock, ok := MockContracts[contractAddr.String()]
	if !ok { // hack for testing without wasmer
		panic("not supported")
	}
	res, gasUsed, execErr := mock.OnIBCPacketTimeout(codeInfo.CodeHash, params, packet, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddr)
	ctx.EventManager().EmitEvents(events)

	if err := k.messenger.DispatchV2(ctx, contractAddr, packet.Source, res.Messages...); err != nil {
		return err
	}
	return nil
}

func (k Keeper) OnConnectChannel(ctx sdk.Context, contractAddr sdk.AccAddress, channel cosmwasm.IBCChannel) error {
	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return err
	}

	var sender sdk.AccAddress // we don't know the sender
	params := cosmwasm.NewEnv(ctx, sender, nil, contractAddr)

	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	mock, ok := MockContracts[contractAddr.String()]
	if !ok { // hack for testing without wasmer
		panic("not supported")
	}
	res, gasUsed, execErr := mock.OnIBCChannelConnect(codeInfo.CodeHash, params, channel, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddr)
	ctx.EventManager().EmitEvents(events)

	if err := k.messenger.DispatchV2(ctx, contractAddr, channel.Endpoint, res.Messages...); err != nil {
		return err
	}
	return nil
}
