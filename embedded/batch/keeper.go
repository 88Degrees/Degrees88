/*

Copyright 2019 All in Bits, Inc
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

package batch

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/types/errs"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TableKey = "batch"

	batchKeyPrefix = "batch"
)

type Keeper struct {
	as  store.ArchiveStore
	cdc *codec.Codec
}

func NewKeeper(db dbm.DB, cdc *codec.Codec) Keeper {
	return Keeper{
		as:  store.NewTable(db, TableKey),
		cdc: cdc,
	}
}

func (k Keeper) LatestByMarket(marketID store.EntityID) (Batch, sdk.Error) {
	var res Batch
	var found bool
	k.as.ReversePrefixIterator(batchIterKey(marketID), func(_ []byte, v []byte) bool {
		k.cdc.MustUnmarshalBinaryBare(v, &res)
		found = true
		return false
	})

	if !found {
		return res, errs.ErrNotFound("batch not found")
	}

	return res, nil
}

func (k Keeper) OnBatchEvent(event types.Batch) {
	batch := Batch{
		BlockNumber:   event.BlockNumber,
		BlockTime:     event.BlockTime,
		MarketID:      event.MarketID,
		ClearingPrice: event.ClearingPrice,
		Bids:          event.Bids,
		Asks:          event.Asks,
	}
	k.as.Set(batchKey(batch.MarketID, batch.BlockNumber), k.cdc.MustMarshalBinaryBare(batch))
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Batch:
		k.OnBatchEvent(ev)
	}

	return nil
}

func batchKey(marketID store.EntityID, blkNum int64) []byte {
	return store.PrefixKeyBytes(batchIterKey(marketID), store.Int64Subkey(blkNum))
}

func batchIterKey(marketID store.EntityID) []byte {
	return store.PrefixKeyString(batchKeyPrefix, marketID.Bytes())
}
