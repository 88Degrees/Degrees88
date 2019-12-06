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

package nft_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/nft"
	"github.com/xar-network/xar-network/x/nft/internal/types"
)

// nolint: deadcode unused
var (
	denom1    = "test-denom"
	denom2    = "test-denom2"
	denom3    = "test-denom3"
	id        = "1"
	id2       = "2"
	id3       = "3"
	address   = types.CreateTestAddrs(1)[0]
	address2  = types.CreateTestAddrs(2)[1]
	address3  = types.CreateTestAddrs(3)[2]
	tokenURI1 = "https://google.com/token-1.json"
	tokenURI2 = "https://google.com/token-2.json"
)

func createTestApp(isCheckTx bool) (*SimApp, sdk.Context) {
	app := Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	return app, ctx
}

// CheckInvariants checks the invariants
func CheckInvariants(k nft.Keeper, ctx sdk.Context) bool {
	collectionsSupply := make(map[string]int)
	ownersCollectionsSupply := make(map[string]int)

	k.IterateCollections(ctx, func(collection types.Collection) bool {
		collectionsSupply[collection.Denom] = collection.Supply()
		return false
	})

	owners := k.GetOwners(ctx)
	for _, owner := range owners {
		for _, idCollection := range owner.IDCollections {
			ownersCollectionsSupply[idCollection.Denom] += idCollection.Supply()
		}
	}

	for denom, supply := range collectionsSupply {
		if supply != ownersCollectionsSupply[denom] {
			fmt.Printf("denom is %s, supply is %d, ownerSupply is %d", denom, supply, ownersCollectionsSupply[denom])
			return false
		}
	}
	return true
}
