/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	csdtcmd "github.com/xar-network/xar-network/x/synthetic/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

// NewModuleClient creates client for the module
func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group nameservice queries under a subcommand
	csdtQueryCmd := &cobra.Command{
		Use:   "csdt",
		Short: "Querying commands for the csdt module",
	}

	csdtQueryCmd.AddCommand(client.GetCommands(
		csdtcmd.GetCmd_GetCsdt(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetCsdts(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetUnderCollateralizedCsdts(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetParams(mc.storeKey, mc.cdc),
	)...)

	return csdtQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	csdtTxCmd := &cobra.Command{
		Use:   "csdt",
		Short: "csdt transactions subcommands",
	}

	csdtTxCmd.AddCommand(client.PostCommands(
		csdtcmd.GetCmdModifyCsdt(mc.cdc),
	)...)

	return csdtTxCmd
}
