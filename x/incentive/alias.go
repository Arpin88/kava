// Code generated by aliasgen tool (github.com/rhuairahrighairidh/aliasgen) DO NOT EDIT.

package incentive

import (
	"github.com/kava-labs/kava/x/incentive/keeper"
	"github.com/kava-labs/kava/x/incentive/types"
)

const (
	BeginningOfMonth                   = keeper.BeginningOfMonth
	MidMonth                           = keeper.MidMonth
	PaymentHour                        = keeper.PaymentHour
	AttributeKeyClaimAmount            = types.AttributeKeyClaimAmount
	AttributeKeyClaimPeriod            = types.AttributeKeyClaimPeriod
	AttributeKeyClaimType              = types.AttributeKeyClaimType
	AttributeKeyClaimedBy              = types.AttributeKeyClaimedBy
	AttributeKeyRewardPeriod           = types.AttributeKeyRewardPeriod
	AttributeValueCategory             = types.AttributeValueCategory
	BondDenom                          = types.BondDenom
	DefaultParamspace                  = types.DefaultParamspace
	EventTypeClaim                     = types.EventTypeClaim
	EventTypeClaimPeriod               = types.EventTypeClaimPeriod
	EventTypeClaimPeriodExpiry         = types.EventTypeClaimPeriodExpiry
	EventTypeRewardPeriod              = types.EventTypeRewardPeriod
	HardLiquidityProviderClaimType     = types.HardLiquidityProviderClaimType
	Large                              = types.Large
	Medium                             = types.Medium
	ModuleName                         = types.ModuleName
	QuerierRoute                       = types.QuerierRoute
	QueryGetClaimPeriods               = types.QueryGetClaimPeriods
	QueryGetHardRewards                = types.QueryGetHardRewards
	QueryGetHardRewardsUnsynced        = types.QueryGetHardRewardsUnsynced
	QueryGetParams                     = types.QueryGetParams
	QueryGetRewardFactors              = types.QueryGetRewardFactors
	QueryGetRewardPeriods              = types.QueryGetRewardPeriods
	QueryGetRewards                    = types.QueryGetRewards
	QueryGetUSDXMintingRewards         = types.QueryGetUSDXMintingRewards
	QueryGetUSDXMintingRewardsUnsynced = types.QueryGetUSDXMintingRewardsUnsynced
	RestClaimCollateralType            = types.RestClaimCollateralType
	RestClaimOwner                     = types.RestClaimOwner
	RestClaimType                      = types.RestClaimType
	RestUnsynced                       = types.RestUnsynced
	RouterKey                          = types.RouterKey
	Small                              = types.Small
	StoreKey                           = types.StoreKey
	USDXMintingClaimType               = types.USDXMintingClaimType
)

var (
	// function aliases
	CalculateTimeElapsed                     = keeper.CalculateTimeElapsed
	NewKeeper                                = keeper.NewKeeper
	NewQuerier                               = keeper.NewQuerier
	DefaultGenesisState                      = types.DefaultGenesisState
	DefaultParams                            = types.DefaultParams
	GetTotalVestingPeriodLength              = types.GetTotalVestingPeriodLength
	NewGenesisAccumulationTime               = types.NewGenesisAccumulationTime
	NewGenesisRewardIndexes                  = types.NewGenesisRewardIndexes
	NewGenesisState                          = types.NewGenesisState
	NewHardLiquidityProviderClaim            = types.NewHardLiquidityProviderClaim
	NewMsgClaimHardReward                    = types.NewMsgClaimHardReward
	NewMsgClaimHardRewardVVesting            = types.NewMsgClaimHardRewardVVesting
	NewMsgClaimUSDXMintingReward             = types.NewMsgClaimUSDXMintingReward
	NewMsgClaimUSDXMintingRewardVVesting     = types.NewMsgClaimUSDXMintingRewardVVesting
	NewMultiRewardIndex                      = types.NewMultiRewardIndex
	NewMultiRewardPeriod                     = types.NewMultiRewardPeriod
	NewMultiplier                            = types.NewMultiplier
	NewParams                                = types.NewParams
	NewPeriod                                = types.NewPeriod
	NewQueryHardRewardsParams                = types.NewQueryHardRewardsParams
	NewQueryHardRewardsUnsyncedParams        = types.NewQueryHardRewardsUnsyncedParams
	NewQueryRewardFactorsParams              = types.NewQueryRewardFactorsParams
	NewQueryRewardsParams                    = types.NewQueryRewardsParams
	NewQueryUSDXMintingRewardsParams         = types.NewQueryUSDXMintingRewardsParams
	NewQueryUSDXMintingRewardsUnsyncedParams = types.NewQueryUSDXMintingRewardsUnsyncedParams
	NewRewardFactor                          = types.NewRewardFactor
	NewRewardIndex                           = types.NewRewardIndex
	NewRewardPeriod                          = types.NewRewardPeriod
	NewUSDXMintingClaim                      = types.NewUSDXMintingClaim
	ParamKeyTable                            = types.ParamKeyTable
	RegisterCodec                            = types.RegisterCodec

	// variable aliases
	DefaultActive                                   = types.DefaultActive
	DefaultClaimEnd                                 = types.DefaultClaimEnd
	DefaultGenesisAccumulationTimes                 = types.DefaultGenesisAccumulationTimes
	DefaultGenesisRewardIndexes                     = types.DefaultGenesisRewardIndexes
	DefaultHardClaims                               = types.DefaultHardClaims
	DefaultMultiRewardPeriods                       = types.DefaultMultiRewardPeriods
	DefaultMultipliers                              = types.DefaultMultipliers
	DefaultRewardPeriods                            = types.DefaultRewardPeriods
	DefaultUSDXClaims                               = types.DefaultUSDXClaims
	ErrAccountNotFound                              = types.ErrAccountNotFound
	ErrClaimExpired                                 = types.ErrClaimExpired
	ErrClaimNotFound                                = types.ErrClaimNotFound
	ErrInsufficientModAccountBalance                = types.ErrInsufficientModAccountBalance
	ErrInvalidAccountType                           = types.ErrInvalidAccountType
	ErrInvalidClaimOwner                            = types.ErrInvalidClaimOwner
	ErrInvalidClaimType                             = types.ErrInvalidClaimType
	ErrInvalidMultiplier                            = types.ErrInvalidMultiplier
	ErrNoClaimsFound                                = types.ErrNoClaimsFound
	ErrRewardPeriodNotFound                         = types.ErrRewardPeriodNotFound
	ErrZeroClaim                                    = types.ErrZeroClaim
	GovDenom                                        = types.GovDenom
	HardBorrowRewardIndexesKeyPrefix                = types.HardBorrowRewardIndexesKeyPrefix
	HardDelegatorRewardFactorKeyPrefix              = types.HardDelegatorRewardFactorKeyPrefix
	HardLiquidityClaimKeyPrefix                     = types.HardLiquidityClaimKeyPrefix
	HardLiquidityRewardDenom                        = types.HardLiquidityRewardDenom
	HardSupplyRewardIndexesKeyPrefix                = types.HardSupplyRewardIndexesKeyPrefix
	IncentiveMacc                                   = types.IncentiveMacc
	KeyClaimEnd                                     = types.KeyClaimEnd
	KeyHardBorrowRewardPeriods                      = types.KeyHardBorrowRewardPeriods
	KeyHardDelegatorRewardPeriods                   = types.KeyHardDelegatorRewardPeriods
	KeyHardSupplyRewardPeriods                      = types.KeyHardSupplyRewardPeriods
	KeyMultipliers                                  = types.KeyMultipliers
	KeyUSDXMintingRewardPeriods                     = types.KeyUSDXMintingRewardPeriods
	ModuleCdc                                       = types.ModuleCdc
	PreviousHardBorrowRewardAccrualTimeKeyPrefix    = types.PreviousHardBorrowRewardAccrualTimeKeyPrefix
	PreviousHardDelegatorRewardAccrualTimeKeyPrefix = types.PreviousHardDelegatorRewardAccrualTimeKeyPrefix
	PreviousHardSupplyRewardAccrualTimeKeyPrefix    = types.PreviousHardSupplyRewardAccrualTimeKeyPrefix
	PreviousUSDXMintingRewardAccrualTimeKeyPrefix   = types.PreviousUSDXMintingRewardAccrualTimeKeyPrefix
	PrincipalDenom                                  = types.PrincipalDenom
	USDXMintingClaimKeyPrefix                       = types.USDXMintingClaimKeyPrefix
	USDXMintingRewardDenom                          = types.USDXMintingRewardDenom
	USDXMintingRewardFactorKeyPrefix                = types.USDXMintingRewardFactorKeyPrefix
)

type (
	Hooks                                 = keeper.Hooks
	Keeper                                = keeper.Keeper
	AccountKeeper                         = types.AccountKeeper
	BaseClaim                             = types.BaseClaim
	BaseMultiClaim                        = types.BaseMultiClaim
	CDPHooks                              = types.CDPHooks
	CdpKeeper                             = types.CdpKeeper
	Claim                                 = types.Claim
	Claims                                = types.Claims
	GenesisAccumulationTime               = types.GenesisAccumulationTime
	GenesisAccumulationTimes              = types.GenesisAccumulationTimes
	GenesisRewardIndexes                  = types.GenesisRewardIndexes
	GenesisRewardIndexesSlice             = types.GenesisRewardIndexesSlice
	GenesisState                          = types.GenesisState
	HARDHooks                             = types.HARDHooks
	HardKeeper                            = types.HardKeeper
	HardLiquidityProviderClaim            = types.HardLiquidityProviderClaim
	HardLiquidityProviderClaims           = types.HardLiquidityProviderClaims
	MsgClaimHardReward                    = types.MsgClaimHardReward
	MsgClaimHardRewardVVesting            = types.MsgClaimHardRewardVVesting
	MsgClaimUSDXMintingReward             = types.MsgClaimUSDXMintingReward
	MsgClaimUSDXMintingRewardVVesting     = types.MsgClaimUSDXMintingRewardVVesting
	MultiRewardIndex                      = types.MultiRewardIndex
	MultiRewardIndexes                    = types.MultiRewardIndexes
	MultiRewardPeriod                     = types.MultiRewardPeriod
	MultiRewardPeriods                    = types.MultiRewardPeriods
	Multiplier                            = types.Multiplier
	MultiplierName                        = types.MultiplierName
	Multipliers                           = types.Multipliers
	Params                                = types.Params
	QueryHardRewardsParams                = types.QueryHardRewardsParams
	QueryHardRewardsUnsyncedParams        = types.QueryHardRewardsUnsyncedParams
	QueryRewardFactorsParams              = types.QueryRewardFactorsParams
	QueryRewardsParams                    = types.QueryRewardsParams
	QueryUSDXMintingRewardsParams         = types.QueryUSDXMintingRewardsParams
	QueryUSDXMintingRewardsUnsyncedParams = types.QueryUSDXMintingRewardsUnsyncedParams
	RewardFactor                          = types.RewardFactor
	RewardFactors                         = types.RewardFactors
	RewardIndex                           = types.RewardIndex
	RewardIndexes                         = types.RewardIndexes
	RewardPeriod                          = types.RewardPeriod
	RewardPeriods                         = types.RewardPeriods
	StakingKeeper                         = types.StakingKeeper
	SupplyKeeper                          = types.SupplyKeeper
	USDXMintingClaim                      = types.USDXMintingClaim
	USDXMintingClaims                     = types.USDXMintingClaims
)
