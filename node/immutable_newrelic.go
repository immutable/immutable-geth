// Copyright 2024 The Immutable go-ethereum Authors
// This file is part of the Immutable go-ethereum library.
//
// The Immutable go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Immutable go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Immutable go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
)

var (
	errBodySizeExceedsLimit = fmt.Errorf("request body size exceeds limit")
)

type RequestBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
	Method  string `json:"method"`
}

func newRelicMiddleware(nrApp *newrelic.Application, next http.Handler) http.Handler { //nolint:unused
	// CHANGE(immutable) add NR agent
	if nrApp == nil {
		log.Error("Failed to initialise New Relic middlware: nrApp is nil")
		return next
	}
	_, handler := newrelic.WrapHandleFunc(nrApp, "/", func(w http.ResponseWriter, req *http.Request) {
		txn := newrelic.FromContext(req.Context())
		txn.AddAttribute("x-api-key", req.Header.Get("x-api-key"))

		// Capture SDK version from headers in order to monitor usage and errors
		txn.AddAttribute("x-sdk-version", req.Header.Get("x-sdk-version"))
		txn.AddAttribute("x-forwarded-for", req.Header.Get("x-forwarded-for"))
		txn.AddAttribute("x-zkevm-rpc-sticky", req.Header.Get("x-zkevm-rpc-sticky"))
		txn.AddAttribute("k6-load-test-id", req.Header.Get("k6-load-test-id"))
		// the W3C trace context header entries.
		// see: https://github.com/w3c/trace-context/blob/main/spec/20-http_request_header_format.md
		// RFC: https://www.w3.org/TR/trace-context/
		// NR, k6 and any other services following the W3C standard can generate these header entries.
		txn.AddAttribute("traceparent", req.Header.Get("traceparent"))
		txn.AddAttribute("tracestate", req.Header.Get("tracestate"))

		// Capture the request body
		mb := int64(1024 * 1024)
		body, bodyReader, err := limitedTeeRead(req.Body, mb)
		if err != nil {
			if errors.Is(err, errBodySizeExceedsLimit) {
				log.Warn(err.Error())
				http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
				return
			}

			log.WithError(err).Error("Failed to read request body")
			// respond to the client with a 503 error.
			// We want to respond here because we can't use the body in subsequent handlers anyway.
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		requestBody := RequestBody{}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			log.Debug("Failed to parse request body", "body", body, "err", err.Error())
		} else {
			// Add the RPC method to the transaction attributes
			txn.AddAttribute("rpcMethod", requestBody.Method)
		}

		// Reset the request body to the original state
		req.Body = bodyReader

		// Next handler
		ctx := newrelic.NewContext(req.Context(), txn)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
	return http.HandlerFunc(handler)
}

// limitedTeeRead reads the request body while writing it back to an io.ReadCloser
// that can be used to read the same data again.
// It returns the content of the input reader, the new io.ReadCloser and an error
// if the content size exceeds the limit.
//
// The returned io.ReadCloser uses the same underlying buffer as the returned []byte.
// This is done for performance reasons. Care should therefore be taken to ensure
// that the returned content buffer is not be written to. As this would affect the
// content of the returned io.ReadCloser downstream.
func limitedTeeRead(reader io.ReadCloser, limitBytes int64) ([]byte, io.ReadCloser, error) {
	limitReader := io.LimitReader(reader, limitBytes)
	content, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, nil, err
	}
	if len(content) >= int(limitBytes) {
		return nil, nil, errBodySizeExceedsLimit
	}
	return content, io.NopCloser(bytes.NewBuffer(content)), nil
}
