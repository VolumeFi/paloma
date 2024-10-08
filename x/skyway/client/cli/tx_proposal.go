package cli

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"cosmossdk.io/math"
	"github.com/VolumeFi/whoops"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/palomachain/paloma/v2/x/skyway/types"
	"github.com/spf13/cobra"
)

const (
	flagExcludedTokens  = "excluded-tokens"
	flagExemptAddresses = "exempt-addresses"
)

func CmdSkywayProposalHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skyway",
		Short: "Skyway proposals",
	}
	cmd.AddCommand([]*cobra.Command{
		CmdSetErc20ToDenom(),
		CmdSetBridgeTax(),
		CmdSetBridgeTransferLimit(),
		CmdSetLightNodeSaleContracts(),
	}...)

	return cmd
}

func applyFlags(cmd *cobra.Command) {
	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
}

func getDeposit(cmd *cobra.Command) (sdk.Coins, error) {
	depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
	whoops.Assert(err)
	return sdk.ParseCoinsNormalized(depositStr)
}

func CmdSetErc20ToDenom() *cobra.Command {
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "set-erc20-to-denom [chain-reference-id] [denom] [erc20]",
		Short: "Sets an association between a denom and an erc20 token for a given chain",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoops.Try(func() {
				cliCtx, err := client.GetClientTxContext(cmd)
				whoops.Assert(err)

				chainReferenceID := args[0]
				denom := args[1]
				erc20 := args[2]

				setERC20ToDenomProposal := &types.SetERC20ToDenomProposal{
					Title:            whoops.Must(cmd.Flags().GetString(cli.FlagTitle)),
					Description:      whoops.Must(cmd.Flags().GetString(cli.FlagSummary)),
					ChainReferenceId: chainReferenceID,
					Erc20:            erc20,
					Denom:            denom,
				}

				from := cliCtx.GetFromAddress()

				deposit, err := getDeposit(cmd)
				whoops.Assert(err)

				msg, err := govv1beta1types.NewMsgSubmitProposal(setERC20ToDenomProposal, deposit, from)
				whoops.Assert(err)

				err = tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
				whoops.Assert(err)
			})
		},
	}
	applyFlags(cmd)
	return cmd
}

func CmdSetBridgeTax() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-bridge-tax [token] [tax-rate]",
		Short:   "Sets the bridge tax rate for a token, and optional exempt addresses",
		Long:    "Each outgoing transfer from Paloma will pay a tax. Tax amount is calculated using [tax-rate], which must be non-negative. [tax-rate] is the ratio of tax collected, so for a 20%% tax it must be set to 0.2.",
		Example: "set-bridge-tax ugrain 0.2",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token, rateRaw := args[0], args[1]

			rate, ok := new(big.Rat).SetString(rateRaw)
			if !ok || rate.Sign() < 0 {
				return fmt.Errorf("invalid tax rate: %s", rateRaw)
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagSummary)
			if err != nil {
				return err
			}

			exemptAddresses, err := cmd.Flags().GetStringSlice(flagExemptAddresses)
			if err != nil {
				return err
			}

			prop := &types.SetBridgeTaxProposal{
				Title:           title,
				Description:     description,
				Token:           token,
				Rate:            rateRaw,
				ExemptAddresses: exemptAddresses,
			}

			from := cliCtx.GetFromAddress()

			deposit, err := getDeposit(cmd)
			if err != nil {
				return err
			}

			msg, err := govv1beta1types.NewMsgSubmitProposal(prop, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice(flagExcludedTokens, []string{},
		"Comma separated list of tokens excluded from the bridge tax. Can be passed multiple times.")
	cmd.Flags().StringSlice(flagExemptAddresses, []string{},
		"Comma separated list of addresses exempt from the bridge tax. Can be passed multiple times.")

	applyFlags(cmd)
	return cmd
}

func CmdSetBridgeTransferLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-bridge-transfer-limit [token] [limit] [limit-period]",
		Short: "Sets the bridge transfer limit, and optionally exempt addresses",
		Long: `Set the bridge transfer limit for the specified token.
[limit-period] must be one of: NONE, DAILY, WEEKLY, MONTHLY, YEARLY. Setting it to NONE effectively disables the limit.
[limit-period] will be converted to a block window. At most, [limit] tokens can be transferred within each block window. After that transfers will fail.`,
		Example: "set-bridge-transfer-limit ugrain 1000000 DAILY",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token, limitRaw, limitPeriodRaw := args[0], args[1], args[2]

			limit, ok := math.NewIntFromString(limitRaw)
			if !ok {
				return fmt.Errorf("invalid limit: %v", limitRaw)
			}

			// Accept both lower case and upper case limit period strings
			limitPeriod, ok := types.LimitPeriod_value[strings.ToUpper(limitPeriodRaw)]
			if !ok {
				return fmt.Errorf("invalid limit period: %v", limitPeriodRaw)
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagSummary)
			if err != nil {
				return err
			}

			exemptAddresses, err := cmd.Flags().GetStringSlice(flagExemptAddresses)
			if err != nil {
				return err
			}

			prop := &types.SetBridgeTransferLimitProposal{
				Title:           title,
				Description:     description,
				Token:           token,
				Limit:           limit,
				LimitPeriod:     types.LimitPeriod(limitPeriod),
				ExemptAddresses: exemptAddresses,
			}

			from := cliCtx.GetFromAddress()

			deposit, err := getDeposit(cmd)
			if err != nil {
				return err
			}

			msg, err := govv1beta1types.NewMsgSubmitProposal(prop, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice(flagExemptAddresses, []string{},
		"Comma separated list of addresses exempt from the bridge tax. Can be passed multiple times.")

	applyFlags(cmd)
	return cmd
}

func CmdSetLightNodeSaleContracts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-light-node-sale-contracts [contracts]",
		Short: "Sets the light node sale contract details for all external chains",
		Long: `Contract details should be entered as a JSON map, with the keys being the external chain reference ID.
The existing set of contract details will be entirely replaced with the new set, so if a chain's contract does not exist in the new value set, it will be removed.`,
		Example: `set-light-node-sale-contracts '{"gnosis-main":{"contract_address":"0x950AA3028F1A3A09D4969C3504BEc30D7ac7d6b2"}}'`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Use a map to get data from the CLI as it's easier to type and
			// format
			var contractsMap map[string]types.LightNodeSaleContract
			err = json.Unmarshal([]byte(args[0]), &contractsMap)
			if err != nil {
				return err
			}

			// Make it a slice for the proposal, as that's what we're using
			// everywhere else
			contracts := make([]*types.LightNodeSaleContract, 0, len(contractsMap))
			for chain, contract := range contractsMap {
				contracts = append(contracts, &types.LightNodeSaleContract{
					ChainReferenceId: chain,
					ContractAddress:  contract.ContractAddress,
				})
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagSummary)
			if err != nil {
				return err
			}

			prop := &types.SetLightNodeSaleContractsProposal{
				Title:                  title,
				Description:            description,
				LightNodeSaleContracts: contracts,
			}

			from := cliCtx.GetFromAddress()

			deposit, err := getDeposit(cmd)
			if err != nil {
				return err
			}

			msg, err := govv1beta1types.NewMsgSubmitProposal(prop, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	applyFlags(cmd)
	return cmd
}
