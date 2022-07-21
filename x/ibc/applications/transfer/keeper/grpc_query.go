package keeper

import (
	"context"
	"fmt"

	"github.com/functionx/fx-core/v2/x/ibc/applications/transfer/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

// DenomTrace implements the Query/DenomTrace gRPC method
func (k Keeper) DenomTrace(c context.Context, req *types.QueryDenomTraceRequest) (*types.QueryDenomTraceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	hash, err := types.ParseHexHash(req.Hash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid denom trace hash %s, %s", req.Hash, err))
	}

	ctx := sdk.UnwrapSDKContext(c)
	denomTrace, found := k.GetDenomTrace(ctx, hash)
	if !found {
		return nil, status.Error(
			codes.NotFound,
			sdkerrors.Wrap(types.ErrTraceNotFound, req.Hash).Error(),
		)
	}

	return &types.QueryDenomTraceResponse{
		DenomTrace: &denomTrace,
	}, nil
}

// DenomTraces implements the Query/DenomTraces gRPC method
func (k Keeper) DenomTraces(c context.Context, req *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	traces := types.Traces{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DenomTraceKey)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		result, err := k.UnmarshalDenomTrace(value)
		if err != nil {
			return err
		}

		traces = append(traces, result)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDenomTracesResponse{
		DenomTraces: traces.Sort(),
		Pagination:  pageRes,
	}, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &params,
	}, nil
}

// DenomHash implements the Query/DenomHash gRPC method
func (k Keeper) DenomHash(c context.Context, req *types.QueryDenomHashRequest) (*types.QueryDenomHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// Convert given request trace path to DenomTrace struct to confirm the path in a valid denom trace format
	denomTrace := types.ParseDenomTrace(req.Trace)
	if err := denomTrace.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	denomHash := denomTrace.Hash()
	found := k.HasDenomTrace(ctx, denomHash)
	if !found {
		return nil, status.Error(
			codes.NotFound,
			sdkerrors.Wrap(types.ErrTraceNotFound, req.Trace).Error(),
		)
	}

	return &types.QueryDenomHashResponse{
		Hash: denomHash.String(),
	}, nil
}
