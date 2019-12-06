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

package order

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/pkg/log"
	"github.com/xar-network/xar-network/types/errs"
	"github.com/xar-network/xar-network/x/order/types"
)

var logger = log.WithModule("order")

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgPost:
			return handleMsgPost(ctx, keeper, msg)
		case types.MsgCancel:
			return handleMsgCancel(ctx, keeper, msg)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unknown message type %v", msg.Type())).Result()
		}
	}
}

func handleMsgPost(ctx sdk.Context, keeper Keeper, msg types.MsgPost) sdk.Result {
	order, err := keeper.Post(
		ctx,
		msg.Owner,
		msg.MarketID,
		msg.Direction,
		msg.Price,
		msg.Quantity,
		msg.TimeInForce,
	)

	if err == nil {
		logger.Info(
			"posted order",
			"id", order.ID.String(),
			"market_id", order.MarketID.String(),
			"price", order.Price.String(),
			"quantity", order.Quantity.String(),
			"direction", order.Direction.String(),
		)
		return sdk.Result{
			Log: fmt.Sprintf("order_id:%s", order.ID),
		}
	}

	return err.Result()
}

func handleMsgCancel(ctx sdk.Context, keeper Keeper, msg types.MsgCancel) sdk.Result {
	order, err := keeper.Get(ctx, msg.OrderID)
	if err != nil {
		return err.Result()
	}
	if !order.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized("cannot cancel unowned order").Result()
	}
	return errs.ErrOrBlankResult(keeper.Del(ctx, order.ID))
}
