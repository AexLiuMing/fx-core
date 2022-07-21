package keeper_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/migrate/types"
)

func (suite *KeeperTestSuite) TestMigrateRecord() {
	var (
		req    *types.QueryMigrateRecordRequest
		expRes *types.QueryMigrateRecordResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"fail - no address",
			func() {
				req = &types.QueryMigrateRecordRequest{}
				expRes = &types.QueryMigrateRecordResponse{}
			},
			false,
		},
		{
			"success - address not migrate",
			func() {
				key := secp256k1.GenPrivKey()
				req = &types.QueryMigrateRecordRequest{
					Address: sdk.AccAddress(key.PubKey().Address()).String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found:         false,
					MigrateRecord: types.MigrateRecord{},
				}
			},
			true,
		},
		{
			"success - address from migrate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.app.MigrateKeeper.SetMigrateRecord(suite.ctx, from, to)

				req = &types.QueryMigrateRecordRequest{
					Address: from.String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found: true,
					MigrateRecord: types.MigrateRecord{
						From:   from.String(),
						To:     to.String(),
						Height: suite.ctx.BlockHeight(),
					},
				}
			},
			true,
		},
		{
			"success - address to migrate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.app.MigrateKeeper.SetMigrateRecord(suite.ctx, from, to)

				req = &types.QueryMigrateRecordRequest{
					Address: to.String(),
				}
				expRes = &types.QueryMigrateRecordResponse{
					Found: true,
					MigrateRecord: types.MigrateRecord{
						From:   from.String(),
						To:     to.String(),
						Height: suite.ctx.BlockHeight(),
					},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)
			tc.malleate()

			res, err := suite.queryClient.MigrateRecord(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Found, res.Found)
				suite.Require().Equal(expRes.MigrateRecord, res.MigrateRecord)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMigrateCheckAccount() {
	var (
		req *types.QueryMigrateCheckAccountRequest
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"fail - no address",
			func() {
				req = &types.QueryMigrateCheckAccountRequest{}
			},
			false,
		},
		{
			"fail - no from address",
			func() {
				toKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					To: to.String(),
				}
			},
			false,
		},
		{
			"fail - no to address",
			func() {
				fromKey := secp256k1.GenPrivKey()
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
				}
			},
			false,
		},
		{
			"success - can migrate",
			func() {
				fromKey := secp256k1.GenPrivKey()
				toKey, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)

				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())
				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"failed - has migrated",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				suite.app.MigrateKeeper.SetMigrateRecord(suite.ctx, from, to)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has delegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1 := validators[0]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, from, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has delegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1 := validators[0]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, to.Bytes(), sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has undelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1 := validators[0]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, from, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, from, val1.GetOperator(), sdk.NewDec(1))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has undelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1 := validators[0]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, to.Bytes(), sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, to.Bytes(), val1.GetOperator())
				suite.Require().True(found)

				_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, to.Bytes(), val1.GetOperator(), delegation.Shares)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has redelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1, val2 := validators[0], validators[1]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, from, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				_, err = suite.app.StakingKeeper.BeginRedelegation(suite.ctx, from, val1.GetOperator(), val2.GetOperator(), sdk.NewDec(1))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"fail - to has redelegate",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 10)
				val1, val2 := validators[0], validators[1]
				//delegate
				_, err := suite.app.StakingKeeper.Delegate(suite.ctx, to.Bytes(), sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000)), stakingtypes.Unbonded, val1, true)
				suite.Require().NoError(err)

				delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, to.Bytes(), val1.GetOperator())
				suite.Require().True(found)

				_, err = suite.app.StakingKeeper.BeginRedelegation(suite.ctx, to.Bytes(), val1.GetOperator(), val2.GetOperator(), delegation.Shares)
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
		{
			"success - from has gov deposit inactive",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, from, amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusDepositPeriod)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - to has gov deposit inactive",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, to.Bytes(), amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusDepositPeriod)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - from has gov deposit active",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, from, amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusVotingPeriod)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - to has gov deposit active",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, to.Bytes(), amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusVotingPeriod)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - from has gov vote",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, from, amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusVotingPeriod)

				err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, from, govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - to has gov vote",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, to.Bytes(), amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusVotingPeriod)

				err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, to.Bytes(), govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			true,
		},
		{
			"success - from and to has gov vote",
			func() {
				fromKey := suite.GenerateAcc(1)[0]
				toKey := suite.GenerateEthAcc(1)[0]
				from := sdk.AccAddress(fromKey.PubKey().Address().Bytes())
				to := common.BytesToAddress(toKey.PubKey().Address().Bytes())

				content := govtypes.ContentFromProposalType("title", "description", "Text")
				amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10000))))

				proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
				suite.Require().NoError(err)

				_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, to.Bytes(), amount)
				suite.Require().NoError(err)

				p, found := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
				suite.Require().True(found)
				suite.Require().Equal(p.Status, govtypes.StatusVotingPeriod)

				err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, from, govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
				suite.Require().NoError(err)

				err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, to.Bytes(), govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
				suite.Require().NoError(err)

				req = &types.QueryMigrateCheckAccountRequest{
					From: from.String(),
					To:   to.String(),
				}
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)
			tc.malleate()

			_, err := suite.queryClient.MigrateCheckAccount(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				//suite.T().Log(err)
			}
		})
	}
}
