package types

const (
	// ModuleName is the name of the module
	ModuleName = "connect"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

const (
	TypeTransfer = "transfer"
)
