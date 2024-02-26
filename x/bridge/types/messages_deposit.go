package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDeposit = "create_deposit"
	TypeMsgUpdateDeposit = "update_deposit"
	TypeMsgDeleteDeposit = "delete_deposit"
)

var _ sdk.Msg = &MsgCreateDeposit{}

func NewMsgCreateDeposit(
	creator string,
	txHash string,
	from string,
	to string,
	coinType CoinType,
	value int64,
	data string,
) *MsgCreateDeposit {
	return &MsgCreateDeposit{
		Creator:  creator,
		TxHash:   txHash,
		From:     from,
		To:       to,
		CoinType: coinType,
		Value:    value,
		Data:     data,
	}
}

func (msg *MsgCreateDeposit) Route() string {
	return RouterKey
}

func (msg *MsgCreateDeposit) Type() string {
	return TypeMsgCreateDeposit
}

func (msg *MsgCreateDeposit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDeposit{}

func NewMsgUpdateDeposit(
	creator string,
	txHash string,
	status DepositStatus,
	rollupTxHash string,
	fromAa string,
) *MsgUpdateDeposit {
	return &MsgUpdateDeposit{
		Creator:      creator,
		TxHash:       txHash,
		Status:       status,
		RollupTxHash: rollupTxHash,
		FromAa:       fromAa,
	}
}

func (msg *MsgUpdateDeposit) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDeposit) Type() string {
	return TypeMsgUpdateDeposit
}

func (msg *MsgUpdateDeposit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.RollupTxHash == "" || msg.FromAa == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "rollup tx hash and aa address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteDeposit{}

func NewMsgDeleteDeposit(
	creator string,
	txHash string,
) *MsgDeleteDeposit {
	return &MsgDeleteDeposit{
		Creator: creator,
		TxHash:  txHash,
	}
}

func (msg *MsgDeleteDeposit) Route() string {
	return RouterKey
}

func (msg *MsgDeleteDeposit) Type() string {
	return TypeMsgDeleteDeposit
}

func (msg *MsgDeleteDeposit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}