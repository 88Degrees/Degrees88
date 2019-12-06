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

package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type PostIssueReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	types.IssueParams `json:"issue"`
}
type PostDescriptionReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Description string       `json:"description"`
}
type PostIssueBaseReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/issue", postIssueHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/approve/{issue-id}/{accAddress}/{amount}", postIssueApproveHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/approve/increase/{issue-id}/{accAddress}/{amount}", postIssueIncreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/approve/decrease/{issue-id}/{accAddress}/{amount}", postIssueDecreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/burn/{issue-id}/{amount}", postBurnHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/burn-from/{issue-id}/{accAddress}/{amount}", postBurnFromHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/freeze/{freeze-type}/{issue-id}/{accAddress}", postIssueFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/unfreeze/{freeze-type}/{issue-id}/{accAddress}", postIssueUnFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/send-from/{issue-id}/{from}/{to}/{amount}", postIssueSendFrom(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/mint/{issue-id}/{amount}/{to}", postMintHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/ownership/transfer/{issue-id}/{to}", postTransferOwnershipHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/description/{issue-id}", postDescribeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/issue/feature/disable/{issue-id}/{feature}", postDisableFeatureHandlerFn(cliCtx)).Methods("POST")
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/issue/params", queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/query/{issue-id}", queryIssueHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/list", queryIssuesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/search/{symbol}", queryIssueSearchHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/freeze/{issue-id}/{address}", queryIssueFreezeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/freezes/{issue-id}", queryIssueFreezesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/issue/allowance/{issue-id}/{address}/{spender_address}", queryIssueAllowanceHandlerFn(cliCtx)).Methods("GET")
}

func postIssueHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(req.Description) > 0 && !json.Valid([]byte(req.Description)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrCoinDescriptionNotValid().Error())
			return
		}
		// create the message
		msg := types.NewMsgIssue(fromAddr, &req.IssueParams)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postMintHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]

		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars["amount"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)
		toAddr, err := sdk.AccAddressFromBech32(vars["to"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}

		msg := types.NewMsgIssueMint(issueID, fromAddr, toAddr, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postDisableFeatureHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]

		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		feature := vars["feature"]
		_, ok := types.Features[feature]
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrUnknownFeatures().Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		_, err = types.IssueOwnerCheck(cliCtx, fromAddress, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgIssueDisableFeature(issueID, fromAddress, feature)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}

}
func postDescribeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostDescriptionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		if len(req.Description) <= 0 || !json.Valid([]byte(req.Description)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrCoinDescriptionNotValid().Error())
			return
		}
		msg := types.NewMsgIssueDescription(issueID, fromAddress, []byte(req.Description))
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = types.IssueOwnerCheck(cliCtx, fromAddress, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func postTransferOwnershipHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		to, err := sdk.AccAddressFromBech32(vars["to"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueTransferOwnership(issueID, fromAddress, to)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = types.IssueOwnerCheck(cliCtx, fromAddress, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueFreezeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cliCtx, true)
}
func postIssueUnFreezeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cliCtx, false)
}
func issueFreezeHandlerFn(cliCtx context.CLIContext, freeze bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		vars := mux.Vars(r)

		msg, err := GetIssueFreezeMsg(cliCtx, fromAddress, vars["freeze-type"], vars["issue-id"], vars["accAddress"], freeze)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postBurnHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cliCtx, types.BurnHolder)
}
func postBurnFromHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cliCtx, types.BurnFrom)
}
func postBurnFromAddressHandlerFn(cliCtx context.CLIContext, burnFromType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)

		issueID := vars["issue-id"]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, ok := sdk.NewIntFromString(vars["amount"])
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Amount not a valid int")
			return
		}

		//burn sender
		accAddress := fromAddr

		if types.BurnFrom == burnFromType {
			//burn from holder address
			accAddress, err = sdk.AccAddressFromBech32(vars["accAddress"])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg, err := GetBurnMsg(cliCtx, fromAddr, accAddress, issueID, amount, burnFromType, false)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func postIssueSendFrom(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		from, err := sdk.AccAddressFromBech32(vars["from"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		to, err := sdk.AccAddressFromBech32(vars["to"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars["amount"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		_, err = GetIssueByID(cliCtx, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := CheckAllowance(cliCtx, issueID, from, sender, amount); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = CheckFreeze(cliCtx, issueID, from, to); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueSendFrom(issueID, sender, from, to, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueApproveHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.Approve)
}
func postIssueIncreaseApproval(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.IncreaseApproval)
}
func postIssueDecreaseApproval(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.DecreaseApproval)
}
func issueApproveHandlerFn(cliCtx context.CLIContext, approveType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)

		issueID := vars["issue-id"]
		accAddr, err := sdk.AccAddressFromBech32(vars["accAddress"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount, ok := sdk.NewIntFromString(vars["amount"])
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Amount not a valid int")
			return
		}

		msg, err := GetIssueApproveMsg(cliCtx, issueID, fromAddr, accAddr, approveType, amount, false)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func burnCheck(sender sdk.AccAddress, burnFrom sdk.AccAddress, issueInfo types.Issue, amount sdk.Int, burnType string, cli bool) error {
	//coins := sender.GetCoins()
	switch burnType {
	case types.BurnOwner:
		{
			if !sender.Equals(issueInfo.GetOwner()) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if !sender.Equals(burnFrom) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if issueInfo.IsBurnOwnerDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			break
		}
	case types.BurnHolder:
		{
			if issueInfo.IsBurnHolderDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			if !sender.Equals(burnFrom) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			break
		}
	case types.BurnFrom:
		{
			if !sender.Equals(issueInfo.GetOwner()) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if issueInfo.IsBurnFromDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			if issueInfo.GetOwner().Equals(burnFrom) {
				//burnFrom
				if issueInfo.IsBurnOwnerDisabled() {
					return types.Errorf(types.ErrCanNotBurn())
				}
			}
			break
		}
	default:
		{
			panic("not support")
		}

	}
	// TODO validate enough funds, need to get an accGetter
	// ensure account has enough coins
	/*if !coins.IsAllGTE(sdk.NewCoins(sdk.NewCoin(issueInfo.GetIssueId(), amount))) {
		return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", sender.GetAddress())
	}*/
	return nil
}

func GetBurnMsg(
	cliCtx context.CLIContext,
	sender sdk.AccAddress,
	burnFrom sdk.AccAddress,
	issueID string,
	amount sdk.Int,
	burnFromType string,
	cli bool,
) (sdk.Msg, error) {
	issueInfo, err := GetIssueByID(cliCtx, issueID)
	if err != nil {
		return nil, err
	}
	if types.BurnHolder == burnFromType {
		if issueInfo.GetOwner().Equals(sender) {
			burnFromType = types.BurnOwner
		}
	}
	err = burnCheck(sender, burnFrom, issueInfo, amount, burnFromType, cli)
	if err != nil {
		return nil, err
	}
	var msg sdk.Msg
	switch burnFromType {

	case types.BurnOwner:
		msg = types.NewMsgIssueBurnOwner(issueID, sender, amount)
		break
	case types.BurnHolder:
		msg = types.NewMsgIssueBurnHolder(issueID, sender, amount)
		break
	case types.BurnFrom:
		msg = types.NewMsgIssueBurnFrom(issueID, sender, burnFrom, amount)
		break
	default:
		return nil, types.ErrCanNotBurn()
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, types.Errorf(err)
	}
	return msg, nil
}

func GetIssueFreezeMsg(
	cliCtx context.CLIContext,
	account sdk.AccAddress,
	freezeType string,
	issueID string,
	address string,
	freeze bool,
) (sdk.Msg, error) {
	_, ok := types.FreezeTypes[freezeType]
	if !ok {
		return nil, types.Errorf(types.ErrUnknownFreezeType())
	}
	if err := types.CheckIssueId(issueID); err != nil {
		return nil, types.Errorf(err)
	}
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	issueInfo, err := types.IssueOwnerCheck(cliCtx, account, issueID)
	if err != nil {
		return nil, err
	}
	if freeze {
		if issueInfo.IsFreezeDisabled() {
			return nil, types.ErrCanNotFreeze()
		}
		msg := types.NewMsgIssueFreeze(issueID, account, accAddress, freezeType)
		if err := msg.ValidateService(); err != nil {
			return msg, types.Errorf(err)
		}
		return msg, nil
	}
	msg := types.NewMsgIssueUnFreeze(issueID, account, accAddress, freezeType)
	if err := msg.ValidateBasic(); err != nil {
		return msg, types.Errorf(err)
	}
	return msg, nil
}

func GetIssueApproveMsg(
	cliCtx context.CLIContext,
	issueID string,
	account sdk.AccAddress,
	accAddress sdk.AccAddress,
	approveType string,
	amount sdk.Int,
	cli bool,
) (sdk.Msg, error) {
	if err := types.CheckIssueId(issueID); err != nil {
		return nil, types.Errorf(err)
	}
	var msg sdk.Msg
	switch approveType {
	case types.Approve:
		msg = types.NewMsgIssueApprove(issueID, account, accAddress, amount)
		break
	case types.IncreaseApproval:
		msg = types.NewMsgIssueIncreaseApproval(issueID, account, accAddress, amount)
		break
	case types.DecreaseApproval:
		msg = types.NewMsgIssueDecreaseApproval(issueID, account, accAddress, amount)
		break
	default:
		return nil, sdk.ErrInternal("not support")
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, types.Errorf(err)
	}
	return msg, nil
}

func CheckAllowance(
	cliCtx context.CLIContext,
	issueID string,
	owner sdk.AccAddress,
	spender sdk.AccAddress,
	amount sdk.Int,
) error {
	res, height, err := cliCtx.QueryWithData(types.GetQueryIssueAllowancePath(issueID, owner, spender), nil)
	if err != nil {
		return err
	}
	cliCtx = cliCtx.WithHeight(height)

	var approval types.Approval
	cliCtx.Codec.MustUnmarshalJSON(res, &approval)

	if approval.Amount.LT(amount) {
		return types.Errorf(types.ErrNotEnoughAmountToTransfer())
	}
	return nil
}

func GetIssueByID(cliCtx context.CLIContext, issueID string) (types.Issue, error) {
	var issueInfo types.Issue
	// Query the issue
	res, height, err := cliCtx.QueryWithData(types.GetQueryIssuePath(issueID), nil)
	if err != nil {
		return nil, err
	}
	cliCtx = cliCtx.WithHeight(height)

	cliCtx.Codec.MustUnmarshalJSON(res, &issueInfo)
	return issueInfo, nil
}

func CheckFreeze(cliCtx context.CLIContext, issueID string, from sdk.AccAddress, to sdk.AccAddress) error {
	res, height, err := cliCtx.QueryWithData(types.GetQueryIssueFreezePath(issueID, from), nil)
	if err != nil {
		return err
	}
	cliCtx = cliCtx.WithHeight(height)

	var freeze types.IssueFreeze
	cliCtx.Codec.MustUnmarshalJSON(res, &freeze)

	res, height, err = cliCtx.QueryWithData(types.GetQueryIssueFreezePath(issueID, to), nil)
	if err != nil {
		return err
	}
	cliCtx = cliCtx.WithHeight(height)

	cliCtx.Codec.MustUnmarshalJSON(res, &freeze)
	return nil
}
