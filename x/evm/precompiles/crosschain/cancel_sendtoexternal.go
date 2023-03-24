package crosschain

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/evm/types"
)

var CancelSendToExternalMethod = abi.NewMethod(
	CancelSendToExternalMethodName,
	CancelSendToExternalMethodName,
	abi.Function, "nonpayable", false, false,
	abi.Arguments{
		abi.Argument{Name: "chain", Type: types.TypeString},
		abi.Argument{Name: "txid", Type: types.TypeUint256},
	},
	abi.Arguments{
		abi.Argument{Name: "result", Type: types.TypeBool},
	},
)

func (c *Contract) CancelSendToExternal(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("cancel send to external method not readonly")
	}

	// args
	args, err := CancelSendToExternalMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	chain, ok0 := args[0].(string)
	txID, ok1 := args[1].(*big.Int)
	if !ok0 || !ok1 {
		return nil, errors.New("unexpected arg type")
	}

	if err := crosschaintypes.ValidateModuleName(chain); err != nil {
		return nil, err
	}
	if txID.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid tx id: %s", txID.String())
	}

	sender := contract.Caller()
	route, has := c.router.GetRoute(chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", chain)
	}

	originDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	// NOTE: must be get relation before cancel, cancel will delete it if relation exist
	hasRelation := c.erc20Keeper.HasOutgoingTransferRelation(ctx, txID.Uint64())

	refundCoin, err := route.PrecompileCancelSendToExternal(ctx, txID.Uint64(), sender.Bytes())
	if err != nil {
		return nil, err
	}
	if !hasRelation && refundCoin.Denom == originDenom {
		// add refund to sender in evm state db, because bank keeper add refund to sender
		evm.StateDB.AddBalance(sender, refundCoin.Amount.BigInt())
	}

	return CancelSendToExternalMethod.Outputs.Pack(true)
}