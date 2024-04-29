package msgserver

import (
	"context"

	"github.com/allora-network/allora-chain/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ms msgServer) AddToWhitelistAdmin(ctx context.Context, msg *types.MsgAddToWhitelistAdmin) (*types.MsgAddToWhitelistAdminResponse, error) {
	// Check that sender is also a whitelist admin
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	targetAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	isAdmin, err := ms.k.IsWhitelistAdmin(ctx, senderAddr)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, types.ErrNotWhitelistAdmin
	}
	// Add the address to the whitelist
	err = ms.k.AddWhitelistAdmin(ctx, targetAddr)
	if err != nil {
		return nil, err
	}
	return &types.MsgAddToWhitelistAdminResponse{}, nil
}

func (ms msgServer) RemoveFromWhitelistAdmin(ctx context.Context, msg *types.MsgRemoveFromWhitelistAdmin) (*types.MsgRemoveFromWhitelistAdminResponse, error) {
	// Check that sender is also a whitelist admin
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	targetAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	isAdmin, err := ms.k.IsWhitelistAdmin(ctx, senderAddr)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, types.ErrNotWhitelistAdmin
	}
	// Remove the address from the whitelist
	err = ms.k.RemoveWhitelistAdmin(ctx, targetAddr)
	if err != nil {
		return nil, err
	}
	return &types.MsgRemoveFromWhitelistAdminResponse{}, nil
}