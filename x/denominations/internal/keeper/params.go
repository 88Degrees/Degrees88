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

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

// GetParams gets params from the store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.GetNomineeParams(ctx))
}

// SetParams updates params in the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetNomineeParams get nominee params from store
func (k Keeper) GetNomineeParams(ctx sdk.Context) []string {
	var nominees []string
	k.paramSubspace.Get(ctx, types.KeyNominees, &nominees)
	return nominees
}
