package cli

import (
	"github.com/xar-network/xar-network/x/authority/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	authorityCmds := &cobra.Command{
		Use:                "authority",
		DisableFlagParsing: false,
		RunE:               client.ValidateCmd,
	}

	authorityCmds.AddCommand(
		client.PostCommands(
			getCmdCreateIssuer(cdc),
			getCmdDestroyIssuer(cdc),
		)...,
	)

	return authorityCmds
}

func getCmdCreateIssuer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "create-issuer [authority_key_or_address] [issuer_address] [denominations]",
		Example: "xarcli authority create-issuer masterkey xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu x2eur,x0jpy",
		Short:   "Create a new issuer",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denoms, err := types.ParseDenominations(args[2])
			if err != nil {
				return err
			}

			msg := types.MsgCreateIssuer{
				Issuer:        issuerAddr,
				Denominations: denoms,
				Authority:     cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdDestroyIssuer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "destroy-issuer [authority_key_or_address] [issuer_address]",
		Example: "xarcli authority destory-issuer masterkey xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Delete an issuer",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgDestroyIssuer{
				Issuer:    issuerAddr,
				Authority: cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
