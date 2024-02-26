//go:build norace
// +build norace

package cli_test

import (
	"fmt"
	"github.com/evmos/ethermint/x/bridge/types"
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/evmos/ethermint/testutil/network"
	"github.com/evmos/ethermint/x/bridge/client/cli"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCreateSignerGroup(t *testing.T) {
	net, err := network.New(t, t.TempDir(), network.DefaultConfig())
	require.NoError(t, err)
	val := net.Validators[0]
	ctx := val.ClientCtx

	fields := []string{"xyz", "1", "abc,xyz"}
	for _, tc := range []struct {
		desc   string
		idName string

		args []string
		err  error
		code uint32
	}{
		{
			idName: strconv.Itoa(0),

			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagGasPrices, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10000000000))).String()),
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idName,
			}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateSignerGroup(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}

func TestUpdateSignerGroup(t *testing.T) {
	net, err := network.New(t, t.TempDir(), network.DefaultConfig())
	require.NoError(t, err)
	val := net.Validators[0]
	ctx := val.ClientCtx

	fields := []string{val.Address.String(), "1", "abc,xyz"}
	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10000000000))).String()),
	}
	args := []string{
		"0",
	}
	args = append(args, fields...)
	args = append(args, common...)
	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateSignerGroup(), args)
	require.NoError(t, err)
	for _, tc := range []struct {
		desc   string
		idName string

		args []string
		code uint32
		err  error
	}{
		{
			desc:   "valid",
			idName: strconv.Itoa(0),

			args: common,
		},
		{
			desc:   "key not found",
			idName: strconv.Itoa(100000),

			args: common,
			code: types.ErrIndexNotExist.ABCICode(),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idName,
			}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateSignerGroup(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}

func TestDeleteSignerGroup(t *testing.T) {
	net, err := network.New(t, t.TempDir(), network.DefaultConfig())
	require.NoError(t, err)

	val := net.Validators[0]
	ctx := val.ClientCtx

	fields := []string{val.Address.String(), "1", "abc,xyz"}
	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdkmath.NewInt(10000000000))).String()),
	}
	args := []string{
		"0",
	}
	args = append(args, fields...)
	args = append(args, common...)
	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateSignerGroup(), args)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc   string
		idName string

		args []string
		code uint32
		err  error
	}{
		{
			desc:   "valid",
			idName: strconv.Itoa(0),

			args: common,
		},
		{
			desc:   "key not found",
			idName: strconv.Itoa(100000),

			args: common,
			code: types.ErrIndexNotExist.ABCICode(),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idName,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdDeleteSignerGroup(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}