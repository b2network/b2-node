package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/evmos/ethermint/testutil/bridge/keeper"
	"github.com/evmos/ethermint/testutil/bridge/nullify"
	"github.com/evmos/ethermint/x/bridge/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCallerGroupQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCallerGroup(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetCallerGroupRequest
		response *types.QueryGetCallerGroupResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCallerGroupRequest{
				Name: msgs[0].Name,
			},
			response: &types.QueryGetCallerGroupResponse{CallerGroup: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCallerGroupRequest{
				Name: msgs[1].Name,
			},
			response: &types.QueryGetCallerGroupResponse{CallerGroup: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCallerGroupRequest{
				Name: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.CallerGroup(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestCallerGroupQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCallerGroup(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCallerGroupRequest {
		return &types.QueryAllCallerGroupRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CallerGroupAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CallerGroup), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CallerGroup),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CallerGroupAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CallerGroup), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CallerGroup),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.CallerGroupAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.CallerGroup),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.CallerGroupAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
