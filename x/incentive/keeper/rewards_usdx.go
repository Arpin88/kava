package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cdptypes "github.com/kava-labs/kava/x/cdp/types"
	"github.com/kava-labs/kava/x/incentive/types"
)

// AccumulateUSDXMintingRewards updates the rewards accumulated for the input reward period
func (k Keeper) AccumulateUSDXMintingRewards(ctx sdk.Context, rewardPeriod types.RewardPeriod) error {
	previousAccrualTime, found := k.GetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType)
	if !found {
		k.SetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType, ctx.BlockTime())
		return nil
	}
	timeElapsed := CalculateTimeElapsed(rewardPeriod.Start, rewardPeriod.End, ctx.BlockTime(), previousAccrualTime)
	if timeElapsed.IsZero() {
		return nil
	}
	if rewardPeriod.RewardsPerSecond.Amount.IsZero() {
		k.SetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType, ctx.BlockTime())
		return nil
	}
	totalPrincipal := k.cdpKeeper.GetTotalPrincipal(ctx, rewardPeriod.CollateralType, types.PrincipalDenom).ToDec()
	if totalPrincipal.IsZero() {
		k.SetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType, ctx.BlockTime())
		return nil
	}
	newRewards := timeElapsed.Mul(rewardPeriod.RewardsPerSecond.Amount)
	cdpFactor, found := k.cdpKeeper.GetInterestFactor(ctx, rewardPeriod.CollateralType)
	if !found {
		k.SetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType, ctx.BlockTime())
		return nil
	}
	rewardFactor := newRewards.ToDec().Mul(cdpFactor).Quo(totalPrincipal)

	previousRewardFactor, found := k.GetUSDXMintingRewardFactor(ctx, rewardPeriod.CollateralType)
	if !found {
		previousRewardFactor = sdk.ZeroDec()
	}
	newRewardFactor := previousRewardFactor.Add(rewardFactor)
	k.SetUSDXMintingRewardFactor(ctx, rewardPeriod.CollateralType, newRewardFactor)
	k.SetPreviousUSDXMintingAccrualTime(ctx, rewardPeriod.CollateralType, ctx.BlockTime())
	return nil
}

// InitializeUSDXMintingClaim creates or updates a claim such that no new rewards are accrued, but any existing rewards are not lost.
// this function should be called after a cdp is created. If a user previously had a cdp, then closed it, they shouldn't
// accrue rewards during the period the cdp was closed. By setting the reward factor to the current global reward factor,
// any unclaimed rewards are preserved, but no new rewards are added.
func (k Keeper) InitializeUSDXMintingClaim(ctx sdk.Context, cdp cdptypes.CDP) {
	_, found := k.GetUSDXMintingRewardPeriod(ctx, cdp.Type)
	if !found {
		// this collateral type is not incentivized, do nothing
		return
	}
	claim, found := k.GetUSDXMintingClaim(ctx, cdp.Owner)
	if !found { // this is the owner's first usdx minting reward claim
		claim = types.NewUSDXMintingClaim(cdp.Owner, sdk.NewCoin(types.USDXMintingRewardDenom, sdk.ZeroInt()), types.RewardIndexes{})
	}
	globalRewardFactor, found := k.GetUSDXMintingRewardFactor(ctx, cdp.Type)
	if !found {
		globalRewardFactor = sdk.ZeroDec()
	}
	claim.RewardIndexes = claim.RewardIndexes.With(cdp.Type, globalRewardFactor)

	k.SetUSDXMintingClaim(ctx, claim)
}

// SynchronizeUSDXMintingReward updates the claim object by adding any accumulated rewards and updating the reward index value.
// this should be called before a cdp is modified.
func (k Keeper) SynchronizeUSDXMintingReward(ctx sdk.Context, cdp cdptypes.CDP) {

	globalRewardFactor, found := k.GetUSDXMintingRewardFactor(ctx, cdp.Type)
	if !found {
		// The global factor is only not found if
		// - the cdp collateral type has not started accumulating rewards yet (either there is no reward specified in params, or the reward start time hasn't been hit)
		// - OR it was wrongly deleted from state (factors should never be removed while unsynced claims exist)
		// If not found we could either skip this sync, or assume the global factor is zero.
		// Skipping will avoid storing unnecessary factors in the claim for non rewarded denoms.
		// And in the event a global factor is wrongly deleted, it will avoid this function panicking when calculating rewards.
		return
	}
	claim, found := k.GetUSDXMintingClaim(ctx, cdp.Owner)
	if !found {
		claim = types.NewUSDXMintingClaim(
			cdp.Owner,
			sdk.NewCoin(types.USDXMintingRewardDenom, sdk.ZeroInt()),
			types.RewardIndexes{},
		)
	}

	userRewardFactor, found := claim.RewardIndexes.Get(cdp.Type)
	if !found {
		// Normally the factor should always be found, as it is added when the cdp is created in InitializeUSDXMintingClaim.
		// However if a cdp type is not rewarded then becomes rewarded (ie a reward period is added to params), existing cdps will not have the factor in their claims.
		// So assume the factor is the starting value for any global factor: 0.
		userRewardFactor = sdk.ZeroDec()
	}

	newRewardsAmount, err := k.CalculateSingleReward(userRewardFactor, globalRewardFactor, cdp.GetTotalPrincipal().Amount.ToDec())
	if err != nil {
		// Global reward factors should never decrease, as it would lead to a negative update to claim.Rewards.
		// This panics if a global reward factor decreases or disappears between the old and new indexes.
		panic(fmt.Sprintf("corrupted global reward indexes found: %v", err))
	}
	newRewardsCoin := sdk.NewCoin(types.USDXMintingRewardDenom, newRewardsAmount)

	claim.Reward = claim.Reward.Add(newRewardsCoin)
	claim.RewardIndexes = claim.RewardIndexes.With(cdp.Type, globalRewardFactor)

	k.SetUSDXMintingClaim(ctx, claim)
}

// SimulateUSDXMintingSynchronization calculates a user's outstanding USDX minting rewards by simulating reward synchronization
func (k Keeper) SimulateUSDXMintingSynchronization(ctx sdk.Context, claim types.USDXMintingClaim) types.USDXMintingClaim {
	for _, ri := range claim.RewardIndexes {
		_, found := k.GetUSDXMintingRewardPeriod(ctx, ri.CollateralType)
		if !found {
			continue
		}

		globalRewardFactor, found := k.GetUSDXMintingRewardFactor(ctx, ri.CollateralType)
		if !found {
			globalRewardFactor = sdk.ZeroDec()
		}

		// the owner has an existing usdx minting reward claim
		index, hasRewardIndex := claim.HasRewardIndex(ri.CollateralType)
		if !hasRewardIndex { // this is the owner's first usdx minting reward for this collateral type
			claim.RewardIndexes = append(claim.RewardIndexes, types.NewRewardIndex(ri.CollateralType, globalRewardFactor))
		}
		userRewardFactor := claim.RewardIndexes[index].RewardFactor
		rewardsAccumulatedFactor := globalRewardFactor.Sub(userRewardFactor)
		if rewardsAccumulatedFactor.IsZero() {
			continue
		}

		claim.RewardIndexes[index].RewardFactor = globalRewardFactor

		cdp, found := k.cdpKeeper.GetCdpByOwnerAndCollateralType(ctx, claim.GetOwner(), ri.CollateralType)
		if !found {
			continue
		}
		newRewardsAmount := rewardsAccumulatedFactor.Mul(cdp.GetTotalPrincipal().Amount.ToDec()).RoundInt()
		if newRewardsAmount.IsZero() {
			continue
		}
		newRewardsCoin := sdk.NewCoin(types.USDXMintingRewardDenom, newRewardsAmount)
		claim.Reward = claim.Reward.Add(newRewardsCoin)
	}

	return claim
}

// SynchronizeUSDXMintingClaim updates the claim object by adding any rewards that have accumulated.
// Returns the updated claim object
func (k Keeper) SynchronizeUSDXMintingClaim(ctx sdk.Context, claim types.USDXMintingClaim) (types.USDXMintingClaim, error) {
	for _, ri := range claim.RewardIndexes {
		cdp, found := k.cdpKeeper.GetCdpByOwnerAndCollateralType(ctx, claim.Owner, ri.CollateralType)
		if !found {
			// if the cdp for this collateral type has been closed, no updates are needed
			continue
		}
		claim = k.synchronizeRewardAndReturnClaim(ctx, cdp)
	}
	return claim, nil
}

// this function assumes a claim already exists, so don't call it if that's not the case
func (k Keeper) synchronizeRewardAndReturnClaim(ctx sdk.Context, cdp cdptypes.CDP) types.USDXMintingClaim {
	k.SynchronizeUSDXMintingReward(ctx, cdp)
	claim, _ := k.GetUSDXMintingClaim(ctx, cdp.Owner)
	return claim
}

// ZeroUSDXMintingClaim zeroes out the claim object's rewards and returns the updated claim object
func (k Keeper) ZeroUSDXMintingClaim(ctx sdk.Context, claim types.USDXMintingClaim) types.USDXMintingClaim {
	claim.Reward = sdk.NewCoin(claim.Reward.Denom, sdk.ZeroInt())
	k.SetUSDXMintingClaim(ctx, claim)
	return claim
}
