// SPDX-License-Identifier: Apache-2.0
// Copyright © 2019 Intel Corporation

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func putPfdAppTransaction(cliCtx context.Context, pfdTs PfdData,
	afCtx *Context, pfdID string, appID string) (PfdData,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tsRet, resp, err := cli.PfdManagementAppPutAPI.PfdAppTransactionPut(cliCtx,
		afCtx.cfg.AfID, pfdID, appID, pfdTs)

	if err != nil {
		return PfdData{}, resp, err
	}
	return tsRet, resp, nil
}

// PutPfdAppTransaction function
func PutPfdAppTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		pfdTs            PfdData
		resp             *http.Response
		pfdTransactionID string
		appID            string
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("Pfd Management App Put: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&pfdTs); err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pfdTransactionID, err = getPfdTransIDFromURL(r.URL)
	if err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appID, err = getPfdAppIDFromURL(r.URL)
	if err != nil {
		log.Errf("Pfd Management App Put: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, resp, err = putPfdAppTransaction(cliCtx, pfdTs, afCtx,
		pfdTransactionID, appID)
	// TBD how to validate the PUT response
	if err != nil {
		log.Errf("Pfd Management App Put : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	w.WriteHeader(resp.StatusCode)
}
