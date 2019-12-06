/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

package record

import (
	"bytes"

	"github.com/xar-network/xar-network/x/record/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/record/internal/keeper"
)

// GenesisState - all record state that must be provided at genesis
type GenesisState struct {
	StartingRecordId uint64              `json:"starting_record_id"`
	Records          []*types.RecordInfo `json:"records"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(startingRecordId uint64) GenesisState {
	return GenesisState{StartingRecordId: startingRecordId}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.RecordMinId)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	err := keeper.SetInitialRecordStartingRecordId(ctx, data.StartingRecordId)
	if err != nil {
		panic(err)
	}

	for _, record := range data.Records {
		keeper.AddRecord(ctx, record)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	genesisState := GenesisState{}

	genesisState.Records = keeper.List(ctx, types.RecordQueryParams{
		Limit: 99999999,
	})

	return genesisState

}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
