package core

import (
	"context"

	"github.com/filecoin-project/go-filecoin/types"
)

// Send executes a message pass inside the VM. If error is set it
// will always satisfy either ShouldRevert() or IsFault().
func Send(ctx context.Context, from, to *types.Actor, msg *types.Message, st types.StateTree) (b []byte, i uint8, err error) {
	ctx = log.Start(ctx, "Send")
	defer func() {
		log.SetTag(ctx, "message-method", msg.Method)
		log.SetTag(ctx, "message-nonce", msg.Nonce)
		log.SetTag(ctx, "message-from", msg.From.String())
		log.SetTag(ctx, "message-to", msg.To.String())

		fmstr, _ := from.Cid()
		log.SetTag(ctx, "from", fmstr.String())

		tostr, _ := to.Cid()
		log.SetTag(ctx, "to", tostr.String())

		log.FinishWithErr(ctx, err)
	}()
	vmCtx := NewVMContext(from, to, msg, st)

	if msg.Value != nil {
		if err := transfer(from, to, msg.Value); err != nil {
			return nil, 1, err
		}
	}

	if msg.Method == "" {
		// if only tokens are transferred there is no need for a method
		// this means we can shortcircuit execution
		return nil, 0, nil
	}

	toExecutable, err := LoadCode(to.Code)
	if err != nil {
		return nil, 1, faultErrorWrap(err, "unable to load code for To actor")
	}

	if !toExecutable.Exports().Has(msg.Method) {
		return nil, 1, newRevertErrorf("missing export: %s", msg.Method)
	}

	return MakeTypedExport(toExecutable, msg.Method)(vmCtx)
}

func transfer(from, to *types.Actor, value *types.TokenAmount) error {
	if value.IsNegative() {
		return ErrCannotTransferNegativeValue
	}

	if from.Balance.LessThan(value) {
		return ErrInsufficientBalance
	}

	if to.Balance == nil {
		to.Balance = types.ZeroToken // This would be unsafe if TokenAmount could be mutated.
	}
	from.Balance = from.Balance.Sub(value)
	to.Balance = to.Balance.Add(value)

	return nil
}