package types

import "github.com/cosmos/cosmos-sdk/codec"

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
}

// RegisterCodec registers the necessary types for cdp module
func RegisterCodec(cdc *codec.Codec) {

}

func RegisterMultiSpend(cdc *codec.Codec) {
	cdc.RegisterConcrete(CommunityPoolMultiSpendProposal{}, "kava/CommunityPoolMultiSpendProposal", nil)
}
