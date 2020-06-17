# Connect

### Setup
- build sample app(./example)
- prepare relayer
- start two-chainz
- start relayer
  - create client, connection, channel

```bash
❯ make initialize
```

- first transaction
```bash
# create nft
❯ ./build/simappcli tx nft mint crypto-kitties ck1 <token_owner_address> --from <sender> --chain-id ibc1 --home scripts/data/ibc1/n0/simappcli --keyring-backend test

# transfer to ibc1 chain
❯ ./build/simappcli tx connect nft-transfer connect ibconexfer 1600 <reciever> crypto-kitties ck1 --chain-id ibc0 --home scripts/data/ibc0/n0/simappcli --from <sender> --keyring-backend test

# check ibc result
❯ ./build/simappcli query nft token connect/ibczeroxfer/crypto-kitties connect/ibczeroxfer/ck1 --chain-id ibc1 --home scripts/data/ibc1/n0/simappcli --keyring-backend test
{"type":"cosmos-sdk/BaseNFT","value":{"id":"connect/ibczeroxfer/ck1","nft_owner":"<owner_address>","token_uri":""}}
```

- About denom, token_id
  - (denom, token_id)_id(in ibc0 chain) - ibc -> port_id/ibc0_channel_id/(denom, token_id)_id(in ibc1 chain) - ibc -> (denom, token_id)_id(in ibc0 chain)
