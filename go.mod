module github.com/shiki-tak/connect

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200522204605-4a07d536a7cc
	github.com/gorilla/mux v1.7.4
	github.com/iqlusioninc/relayer v0.5.4 // indirect
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shiki-tak/relayer v0.5.4
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/shiki-tak/relayer => /dev/relayer // FIXME:
