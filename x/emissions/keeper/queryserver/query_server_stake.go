package queryserver

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	cosmosMath "cosmossdk.io/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-chain/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalStake defines the handler for the Query/TotalStake RPC method.
func (qs queryServer) GetTotalStake(ctx context.Context, req *types.QueryTotalStakeRequest) (*types.QueryTotalStakeResponse, error) {
	totalStake, err := qs.k.GetTotalStake(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryTotalStakeResponse{Amount: totalStake}, nil
}

// Retrieves all stake in a topic for a given reputer address,
// including reputer's stake in themselves and stake delegated to them.
// Also includes stake that is queued for removal.
func (qs queryServer) GetReputerStakeInTopic(ctx context.Context, req *types.QueryReputerStakeInTopicRequest) (*types.QueryReputerStakeInTopicResponse, error) {
	if err := qs.k.ValidateStringIsBech32(req.Address); err != nil {
		return nil, err
	}
	stake, err := qs.k.GetStakeOnReputerInTopic(ctx, req.TopicId, req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryReputerStakeInTopicResponse{Amount: stake}, nil
}

// Retrieves all stake in a topic for a given set of reputer addresses,
// including their stake in themselves and stake delegated to them.
// Also includes stake that is queued for removal.
func (qs queryServer) GetMultiReputerStakeInTopic(ctx context.Context, req *types.QueryMultiReputerStakeInTopicRequest) (*types.QueryMultiReputerStakeInTopicResponse, error) {
	maxLimit := types.DefaultParams().MaxPageLimit
	moduleParams, err := qs.k.GetParams(ctx)
	if err == nil {
		maxLimit = moduleParams.MaxPageLimit
	}

	if uint64(len(req.Addresses)) > maxLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("cannot query more than %d addresses at once", maxLimit))
	}

	stakes := make([]*types.StakeInfo, len(req.Addresses))
	for i, address := range req.Addresses {
		stake := cosmosMath.ZeroInt()
		if err := qs.k.ValidateStringIsBech32(address); err == nil {
			stake, err = qs.k.GetStakeOnReputerInTopic(ctx, req.TopicId, address)
			if err != nil {
				stake = cosmosMath.ZeroInt()
			}
		}
		stakes[i] = &types.StakeInfo{TopicId: req.TopicId, Reputer: address, Amount: stake}
	}

	return &types.QueryMultiReputerStakeInTopicResponse{Amounts: stakes}, nil
}

// Retrieves total delegate stake on a given reputer address in a given topic
func (qs queryServer) GetDelegateStakeInTopicInReputer(ctx context.Context, req *types.QueryDelegateStakeInTopicInReputerRequest) (*types.QueryDelegateStakeInTopicInReputerResponse, error) {
	if err := qs.k.ValidateStringIsBech32(req.ReputerAddress); err != nil {
		return nil, err
	}
	stake, err := qs.k.GetDelegateStakeUponReputer(ctx, req.TopicId, req.ReputerAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegateStakeInTopicInReputerResponse{Amount: stake}, nil
}

func (qs queryServer) GetStakeFromDelegatorInTopicInReputer(ctx context.Context, req *types.QueryStakeFromDelegatorInTopicInReputerRequest) (*types.QueryStakeFromDelegatorInTopicInReputerResponse, error) {
	if err := qs.k.ValidateStringIsBech32(req.ReputerAddress); err != nil {
		return nil, errors.Wrapf(err, "reputer address")
	}
	if err := qs.k.ValidateStringIsBech32(req.DelegatorAddress); err != nil {
		return nil, errors.Wrapf(err, "delegator address")
	}
	stake, err := qs.k.GetDelegateStakePlacement(ctx, req.TopicId, req.DelegatorAddress, req.ReputerAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakeFromDelegatorInTopicInReputerResponse{Amount: stake.Amount.SdkIntTrim()}, nil
}

func (qs queryServer) GetStakeFromDelegatorInTopic(ctx context.Context, req *types.QueryStakeFromDelegatorInTopicRequest) (*types.QueryStakeFromDelegatorInTopicResponse, error) {
	if err := qs.k.ValidateStringIsBech32(req.DelegatorAddress); err != nil {
		return nil, err
	}
	stake, err := qs.k.GetStakeFromDelegatorInTopic(ctx, req.TopicId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakeFromDelegatorInTopicResponse{Amount: stake}, nil
}

// Retrieves total stake in a given topic
func (qs queryServer) GetTopicStake(ctx context.Context, req *types.QueryTopicStakeRequest) (*types.QueryTopicStakeResponse, error) {
	stake, err := qs.k.GetTopicStake(ctx, req.TopicId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTopicStakeResponse{Amount: stake}, nil
}

func (qs queryServer) GetStakeRemovalsForBlock(
	ctx context.Context,
	req *types.QueryStakeRemovalsForBlockRequest,
) (*types.QueryStakeRemovalsForBlockResponse, error) {
	removals, err := qs.k.GetStakeRemovalsForBlock(ctx, req.BlockHeight)
	if err != nil {
		return nil, err
	}
	maxLimit := types.DefaultParams().MaxPageLimit
	moduleParams, err := qs.k.GetParams(ctx)
	if err == nil {
		maxLimit = moduleParams.MaxPageLimit
	}
	if uint64(len(removals)) > maxLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("cannot query more than %d removals at once", maxLimit))
	}
	removalPointers := make([]*types.StakeRemovalInfo, len(removals))
	for i, removal := range removals {
		removalPointers[i] = &removal
	}
	return &types.QueryStakeRemovalsForBlockResponse{Removals: removalPointers}, nil
}

func (qs queryServer) GetDelegateStakeRemovalsForBlock(
	ctx context.Context,
	req *types.QueryDelegateStakeRemovalsForBlockRequest,
) (*types.QueryDelegateStakeRemovalsForBlockResponse, error) {
	removals, err := qs.k.GetDelegateStakeRemovalsForBlock(ctx, req.BlockHeight)
	if err != nil {
		return nil, err
	}
	maxLimit := types.DefaultParams().MaxPageLimit
	moduleParams, err := qs.k.GetParams(ctx)
	if err == nil {
		maxLimit = moduleParams.MaxPageLimit
	}
	if uint64(len(removals)) > maxLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("cannot query more than %d removals at once", maxLimit))
	}
	removalPointers := make([]*types.DelegateStakeRemovalInfo, len(removals))
	for i, removal := range removals {
		removalPointers[i] = &removal
	}
	return &types.QueryDelegateStakeRemovalsForBlockResponse{Removals: removalPointers}, nil
}

func (qs queryServer) GetStakeRemovalInfo(
	ctx context.Context,
	req *types.QueryStakeRemovalInfoRequest,
) (*types.QueryStakeRemovalInfoResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	removal, found, err := qs.k.GetStakeRemovalForReputerAndTopicId(sdkCtx, req.Reputer, req.TopicId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, status.Error(codes.NotFound, "no stake removal found")
	}
	return &types.QueryStakeRemovalInfoResponse{Removal: &removal}, nil
}

func (qs queryServer) GetDelegateStakeRemovalInfo(
	ctx context.Context,
	req *types.QueryDelegateStakeRemovalInfoRequest,
) (*types.QueryDelegateStakeRemovalInfoResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	removal, found, err := qs.k.
		GetDelegateStakeRemovalForDelegatorReputerAndTopicId(sdkCtx, req.Delegator, req.Reputer, req.TopicId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, status.Error(codes.NotFound, "no delegate stake removal found")
	}
	return &types.QueryDelegateStakeRemovalInfoResponse{Removal: &removal}, nil
}
