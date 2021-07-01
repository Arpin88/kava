package swap

import (
	"github.com/kava-labs/kava/x/swap/keeper"
	"github.com/kava-labs/kava/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for swap messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		if deadlineMsg, ok := msg.(types.MsgWithDeadline); ok {
			if deadlineExceeded := deadlineMsg.DeadlineExceeded(ctx.BlockTime()); deadlineExceeded {
				return nil, sdkerrors.Wrapf(types.ErrDeadlineExceeded, "block time %d >= deadline %d", ctx.BlockTime().Unix(), deadlineMsg.GetDeadline().Unix())
			}
		}

		switch msg := msg.(type) {
		case types.MsgDeposit:
			return handleMsgDeposit(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgDeposit(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeposit) (*sdk.Result, error) {
	err := k.Deposit(ctx, msg.Depositor, msg.TokenA, msg.TokenB, msg.Slippage)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}