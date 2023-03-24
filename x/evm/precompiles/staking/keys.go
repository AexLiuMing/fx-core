package staking

import (
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const (
	// DelegateGas default delegate use 0.76FX
	DelegateGas = 400000 // if gas price 500Gwei, fee is 0.2FX
	// UndelegateGas default undelegate use 0.82FX
	UndelegateGas = 600000 // if gas price 500Gwei, fee is 0.3FX
	// WithdrawGas default withdraw use 0.56FX
	WithdrawGas          = 300000 // if gas price 500Gwei, fee is 0.15FX
	DelegationGas        = 200000 // if gas price 500Gwei, fee is 0.1FX
	DelegationRewardsGas = 200000

	DelegateMethodName          = "delegate"
	UndelegateMethodName        = "undelegate"
	WithdrawMethodName          = "withdraw"
	DelegationMethodName        = "delegation"
	DelegationRewardsMethodName = "delegationRewards"

	JsonABI = `[{"type":"function","name":"delegate","inputs":[{"name":"validator","type":"string"}],"outputs":[{"name":"shares","type":"uint256"},{"name":"reward","type":"uint256"}],"payable":true,"stateMutability":"payable"},{"type":"function","name":"undelegate","inputs":[{"name":"validator","type":"string"},{"name":"shares","type":"uint256"}],"outputs":[{"name":"amount","type":"uint256"},{"name":"reward","type":"uint256"},{"name":"endTime","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"withdraw","inputs":[{"name":"validator","type":"string"}],"outputs":[{"name":"reward","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"delegation","inputs":[{"name":"validator","type":"string"},{"name":"delegator","type":"address"}],"outputs":[{"name":"delegate","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"delegationRewards","inputs":[{"name":"validator","type":"string"},{"name":"delegator","type":"address"}],"outputs":[{"name":"rewards","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"}]`
)

var precompileAddress = common.HexToAddress(fxtypes.StakingAddress)

func GetPrecompileAddress() common.Address {
	return precompileAddress
}