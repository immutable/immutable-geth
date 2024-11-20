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

package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ethereum/go-ethereum/log"
)

type ObjectStore struct {
	os *s3.S3
}

func NewObjectStore(region string) (*ObjectStore, error) {
	os := s3.New(session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(region))
	return &ObjectStore{
		os: os,
	}, nil
}

func (os ObjectStore) PutObject(bucket, key string, content []byte) error {
	log.Info("Pushing S3 object", "name", fmt.Sprintf("%s/%s", bucket, key))
	if _, err := os.os.PutObject(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(content),
		ACL:    aws.String("private"),
	}); err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	return nil
}
