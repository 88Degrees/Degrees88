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

package fill

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	TableKey = "fill"
)

type IteratorCB func(fill Fill) bool

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

func (k Keeper) OnFillEvent(event types.Fill) {
	fill := Fill{
		OrderID:     event.OrderID,
		Owner:       event.Owner,
		Pair:        event.Pair,
		Direction:   event.Direction,
		QtyFilled:   event.QtyFilled,
		QtyUnfilled: event.QtyUnfilled,
		BlockNumber: event.BlockNumber,
		Price:       event.Price,
	}
	storedB := k.cdc.MustMarshalBinaryBare(fill)
	k.as.Set(fillKey(event.BlockNumber, event.OrderID), storedB)
}

func (k Keeper) IterOverBlockNumbers(start int64, end int64, cb IteratorCB) {
	k.as.Iterator(fillIterKey(start), fillIterKey(end), func(_ []byte, v []byte) bool {
		var fill Fill
		k.cdc.MustUnmarshalBinaryBare(v, &fill)
		return cb(fill)
	})
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Fill:
		k.OnFillEvent(ev)
	}

	return nil
}

func fillIterKey(blockNum int64) []byte {
	return store.PrefixKeyBytes(store.Int64Subkey(blockNum))
}

func fillKey(blockNum int64, orderId store.EntityID) []byte {
	return store.PrefixKeyBytes(fillIterKey(blockNum), orderId.Bytes())
}
