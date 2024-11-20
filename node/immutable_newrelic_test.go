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
	"errors"
	"io"
	"testing"
)

func TestLimitedTeeReader_Undersized(t *testing.T) {
	body := []byte("xxxxxxx")
	limitBytes := int64(len(body)) * 2
	oldReader := io.NopCloser(bytes.NewReader(body))
	content, newReader, err := limitedTeeRead(oldReader, limitBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(content, body) {
		t.Fatalf("expected content to be the same, got: %v", content)
	}
	readContent, err := io.ReadAll(newReader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(readContent, body) {
		t.Fatalf("expected read content to be the same, got: %v", readContent)
	}
}

func TestLimitedTeeReader_Oversized(t *testing.T) {
	first := []byte("xxxxxxx")
	second := []byte("yyyyyyy")
	body := append(first, second...)
	limitBytes := int64(len(first))
	readCloser := io.NopCloser(bytes.NewReader(body))
	_, _, err := limitedTeeRead(readCloser, limitBytes)
	if !errors.Is(err, errBodySizeExceedsLimit) {
		t.Fatalf("expected error specific error, got: %v", err)
	}
}
