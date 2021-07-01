package types

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	zero = sdk.ZeroInt()
)

// calculateInitialShares calculates initial shares as sqrt(A*B), the geometric mean of A and B
func calculateInitialShares(reservesA, reservesB sdk.Int) sdk.Int {
	// Big.Int allows multiplication without overflow at 255 bits.
	// In addition, Sqrt converges to a correct solution for inputs
	// where sdk.Int.ApproxSqrt does not converge due to exceeding
	// 100 iterations.
	var result big.Int
	result.Mul(reservesA.BigInt(), reservesB.BigInt()).Sqrt(&result)
	return sdk.NewIntFromBigInt(&result)
}

// BasePool implements a unitless constant-product liquidity pool.
//
// The pool is symmetric. For all A,B,s, any operation F on a pool (A,B,s) and pool (B,A,s)
// will result in equal state values of A', B', s': F(A,B,s) => (A',B',s'), F(B,A,s) => (B',A',s')
//
// In addition, the pool is protected from overflow in intermediate calculations, and will
// only overflow when A, B, or s become larger than the max sdk.Int.
//
// Pool operations with non-positive values are invalid, and all functions on a pool will panic
// when given zero or negative values.
type BasePool struct {
	reservesA   sdk.Int
	reservesB   sdk.Int
	totalShares sdk.Int
}

// NewBasePool returns a pointer to a base pool with reserves and total shares initialized
func NewBasePool(reservesA, reservesB sdk.Int) (*BasePool, error) {
	if reservesA.LTE(zero) || reservesB.LTE(zero) {
		return nil, sdkerrors.Wrap(ErrInvalidPool, "reserves must be greater than zero")
	}

	totalShares := calculateInitialShares(reservesA, reservesB)

	return &BasePool{
		reservesA:   reservesA,
		reservesB:   reservesB,
		totalShares: totalShares,
	}, nil
}

// NewBasePoolWithExistingShares returns a pointer to a base pool with existing shares
func NewBasePoolWithExistingShares(reservesA, reservesB, totalShares sdk.Int) (*BasePool, error) {
	if reservesA.LTE(zero) || reservesB.LTE(zero) {
		return nil, sdkerrors.Wrap(ErrInvalidPool, "reserves must be greater than zero")
	}

	if totalShares.LTE(zero) {
		return nil, sdkerrors.Wrap(ErrInvalidPool, "total shares must be greater than zero")
	}

	return &BasePool{
		reservesA:   reservesA,
		reservesB:   reservesB,
		totalShares: totalShares,
	}, nil
}

// ReservesA returns the A reserves of the pool
func (p *BasePool) ReservesA() sdk.Int {
	return p.reservesA
}

// ReservesB returns the B reserves of the pool
func (p *BasePool) ReservesB() sdk.Int {
	return p.reservesB
}

// IsEmpty returns true if all reserves are zero and
// returns false if reserveA or reserveB is not empty
func (p *BasePool) IsEmpty() bool {
	return p.reservesA.IsZero() && p.reservesB.IsZero()
}

// TotalShares returns the total number of shares in the pool
func (p *BasePool) TotalShares() sdk.Int {
	return p.totalShares
}

// AddLiquidity adds liquidity to the pool returns the actual reservesA, reservesB deposits in addition
// to the number of shares created.  The deposits are always less than or equal to the provided and desired
// values.
func (p *BasePool) AddLiquidity(desiredA sdk.Int, desiredB sdk.Int) (sdk.Int, sdk.Int, sdk.Int) {
	// Panics if provided values are zero
	p.assertDepositsArePositive(desiredA, desiredB)

	// Reinitialize the pool if reserves are empty and return the initialized state.
	if p.IsEmpty() {
		p.reservesA = desiredA
		p.reservesB = desiredB
		p.totalShares = calculateInitialShares(desiredA, desiredB)
		return p.ReservesA(), p.ReservesB(), p.TotalShares()
	}

	// Panics if reserveA or reserveB is zero.
	p.assertReservesArePositive()

	// In order to preserve the reserve ratio of the pool, we must deposit
	// A and B in the same ratio of the existing reserves.  In addition,
	// we should not deposit more funds than requested.
	//
	// To meet these requirements, we first calculate the optimalB to deposit
	// if we keep desiredA fixed.  If this is less than or equal to the desiredB,
	// then we use (desiredA, optimalB) as the deposit.
	//
	// If the optimalB is greater than the desiredB, we calculate the optimalA
	// from the desiredB and use (optimalA, optimalB) as the deposit.
	//
	// These optimal values are calculated as:
	//
	// optimalB = reservesB * desiredA / reservesA
	// optimalA = reservesA * desiredB / reservesB
	//
	// Which shows us:
	//
	// if optimalB < desiredB then optimalA > desiredA
	// if optimalB = desiredB then optimalA = desiredA
	// if optimalB > desiredB then optimalA < desiredA
	//
	// so we first check if optimalB <= desiredB, then deposit
	// (desiredA, optimalB) else deposit (optimalA, desiredA).
	//
	// In order avoid precision loss, we rearrange the inequality
	// of optimalB <= desiredB
	// from:
	//   reservesB * desiredA / reservesA <= desiredB
	// to:
	//   reservesB * desiredA <= desiredB * reservesA
	//
	// which also shares the same intermediate products
	// as the calculations for optimalB and optimalA.
	actualA := desiredA.BigInt()
	actualB := desiredB.BigInt()

	// productA = reservesB * desiredA
	var productA big.Int
	productA.Mul(p.reservesB.BigInt(), actualA)

	// productB = reservesA * desiredB
	var productB big.Int
	productB.Mul(p.reservesA.BigInt(), actualB)

	// optimalB <= desiredB
	if productA.Cmp(&productB) <= 0 {
		actualB.Quo(&productA, p.reservesA.BigInt())
	} else { // optimalA < desiredA
		actualA.Quo(&productB, p.reservesB.BigInt())
	}

	var sharesA big.Int
	sharesA.Mul(actualA, p.totalShares.BigInt()).Quo(&sharesA, p.reservesA.BigInt())

	var sharesB big.Int
	sharesB.Mul(actualB, p.totalShares.BigInt()).Quo(&sharesB, p.reservesB.BigInt())

	// a/A and b/B may not be equal due to discrete math and truncation errors,
	// so use the smallest deposit ratio to calculate the number of shares
	//
	// If we do not use the min or max ratio, then the result becomes
	// dependent on the order of reserves in the pool
	//
	// Min is used to always ensure the share ratio is never larger
	// than the deposit ratio for either A or B, ensuring there are no
	// cases where a withdraw will allow funds to be removed at a higher ratio
	// than it was deposited.
	var shares sdk.Int
	if sharesA.Cmp(&sharesB) <= 0 {
		shares = sdk.NewIntFromBigInt(&sharesA)
	} else {
		shares = sdk.NewIntFromBigInt(&sharesB)
	}

	depositA := sdk.NewIntFromBigInt(actualA)
	depositB := sdk.NewIntFromBigInt(actualB)

	// update internal pool state
	p.reservesA = p.reservesA.Add(depositA)
	p.reservesB = p.reservesB.Add(depositB)
	p.totalShares = p.totalShares.Add(shares)

	return depositA, depositB, shares
}

// RemoveLiquidity removes liquidity from the pool and panics if the
// the shares provided are greater than the total shares of the pool
// or the shares are not positive.
// In addition, also panics if reserves go negative, which should not happen.
// If panic occurs, it is a bug.
func (p *BasePool) RemoveLiquidity(shares sdk.Int) (sdk.Int, sdk.Int) {
	// calculate amount to withdraw from the pool based
	// on the number of shares provided. s/S * reserves
	withdrawA, withdrawB := p.ShareValue(shares)

	// update internal pool state
	p.reservesA = p.reservesA.Sub(withdrawA)
	p.reservesB = p.reservesB.Sub(withdrawB)
	p.totalShares = p.totalShares.Sub(shares)

	// Panics if reserveA or reserveB are negative
	// A zero value (100% withdraw) is OK and should not panic.
	p.assertReservesAreNotNegative()

	return withdrawA, withdrawB
}

// SwapAForB trades a for b with a percentage fee
func (p *BasePool) SwapAForB(a sdk.Int, fee sdk.Dec) sdk.Int {
	// TODO: implementation
	return sdk.ZeroInt()
}

// SwapBForA trades b for a with a percentage fee
func (p *BasePool) SwapBForA(b sdk.Int, fee sdk.Dec) sdk.Int {
	// TODO: implementation
	return sdk.ZeroInt()
}

// ShareValue returns the value of the provided shares and panics
// if the shares are greater than the total shares of the pool or
// if the shares are not positive.
func (p *BasePool) ShareValue(shares sdk.Int) (sdk.Int, sdk.Int) {
	p.assertSharesArePositive(shares)
	p.assertSharesAreLessThanTotal(shares)

	var resultA big.Int
	resultA.Mul(p.reservesA.BigInt(), shares.BigInt())
	resultA.Quo(&resultA, p.totalShares.BigInt())

	var resultB big.Int
	resultB.Mul(p.reservesB.BigInt(), shares.BigInt())
	resultB.Quo(&resultB, p.totalShares.BigInt())

	return sdk.NewIntFromBigInt(&resultA), sdk.NewIntFromBigInt(&resultB)
}

// assertSharesPositive panics if shares is zero or negative
func (p *BasePool) assertSharesArePositive(shares sdk.Int) {
	if !shares.IsPositive() {
		panic("invalid value: shares must be positive")
	}
}

// assertSharesLessThanTotal panics if the number of shares is greater than the total shares
func (p *BasePool) assertSharesAreLessThanTotal(shares sdk.Int) {
	if shares.GT(p.totalShares) {
		panic(fmt.Sprintf("out of bounds: shares %s > total shares %s", shares, p.totalShares))
	}
}

// assertDepositsPositive panics if a deposit is zero or negative
func (p *BasePool) assertDepositsArePositive(depositA, depositB sdk.Int) {
	if !depositA.IsPositive() {
		panic("invalid value: deposit A must be positive")
	}

	if !depositB.IsPositive() {
		panic("invalid state: deposit B must be positive")
	}
}

// assertReservesArePositive panics if any reserves are zero.  This is an invalid
// state that should never happen.  If this panic is seen, it is a bug.
func (p *BasePool) assertReservesArePositive() {
	if !p.reservesA.IsPositive() {
		panic("invalid state: reserves A must be positive")
	}

	if !p.reservesB.IsPositive() {
		panic("invalid state: reserves B must be positive")
	}
}

// assertReservesAreNotNegative panics if any reserves are negative.  This is an invalid
// state that should never happen.  If this panic is seen, it is a bug.
func (p *BasePool) assertReservesAreNotNegative() {
	if p.reservesA.IsNegative() {
		panic("invalid state: reserves A can not be negative")
	}

	if p.reservesB.IsNegative() {
		panic("invalid state: reserves B can not be negative")
	}
}